package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

// Cool Colors for output.
const (
	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
	orange  = "\033[38;5;208m"
	reset   = "\033[0m"
)

const githubAPIBaseURL = "https://api.github.com"

type GitHubUser struct {
	Login       string `json:"login"`
	Name        string `json:"name"`
	PublicRepos int    `json:"public_repos"`
	Followers   int    `json:"followers"`
	HTMLURL     string `json:"html_url"`
	BIO         string `json:"bio"`
}

type GitHubRepo struct {
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Stars       int    `json:"stargazers_count"`
	Forks       int    `json:"forks_count"`
	Language    string `json:"language"`
	HTMLURL     string `json:"html_url"`
}

type GitHubClient struct {
	httpClient *http.Client
	baseURL    string
}

func newGitHubClient() GitHubClient {
	return GitHubClient{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: githubAPIBaseURL,
	}
}

func (gh GitHubClient) fetchUser(username string) (*GitHubUser, error) {

	apiURL := gh.baseURL + "/users/" + url.PathEscape(username)

	request, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("User-Agent", "go-api-practice")

	response, err := gh.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("user %q was not found", username)
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub returned: %s", response.Status)
	}

	var user GitHubUser

	err = json.NewDecoder(response.Body).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (gh GitHubClient) fetchRepo(owner string, repoName string) (*GitHubRepo, error) {

	apiURL := gh.baseURL + "/repos/" + url.PathEscape(owner) + "/" + url.PathEscape(repoName)

	request, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("User-Agent", "go-api-practice")

	response, err := gh.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("repo %q was not found", owner+"/"+repoName)
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub returned: %s", response.Status)
	}

	var repo GitHubRepo

	err = json.NewDecoder(response.Body).Decode(&repo)
	if err != nil {
		return nil, err
	}

	return &repo, nil
}

func main() {
	// Prints a friendly message to user...
	fmt.Println(yellow + "GitHub lookup..." + reset)

	rootCmd := &cobra.Command{
		Use:   "github-tool",
		Short: "Look up GitHub information from the command line",
	}

	userCmd := &cobra.Command{
		Use:   "user <username>",
		Short: "Look up a GitHub user",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			username := args[0]

			githubClient := newGitHubClient()

			user, err := githubClient.fetchUser(username)
			if err != nil {
				fmt.Println(red+"Error:"+reset, err)
				return
			}
			// print some cool stuff to the terminal :P
			fmt.Println("Login:", green+user.Login+reset)
			fmt.Println("Name:", green+user.Name+reset)
			fmt.Println("Public repos:", orange+strconv.Itoa(user.PublicRepos)+reset)
			fmt.Println("Followers:", magenta+strconv.Itoa(user.Followers)+reset)
			fmt.Println("Profile:", cyan+user.HTMLURL+reset)
			fmt.Println("Bio:", orange+user.BIO+reset)
		},
	}

	repoCmd := &cobra.Command{
		Use:   "repo <owner> <repo>",
		Short: "Look up a GitHub repo",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			owner := args[0]
			repoName := args[1]

			githubClient := newGitHubClient()

			repo, err := githubClient.fetchRepo(owner, repoName)
			if err != nil {
				fmt.Println(red+"Error:"+reset, err)
				return
			}

			fmt.Println("Repo:", green+repo.FullName+reset)
			fmt.Println("Description:", repo.Description)
			fmt.Println("Stars:", yellow+strconv.Itoa(repo.Stars)+reset)
			fmt.Println("Forks:", magenta+strconv.Itoa(repo.Forks)+reset)
			fmt.Println("Language:", blue+repo.Language+reset)
			fmt.Println("URL:", cyan+repo.HTMLURL+reset)
		},
	}

	rootCmd.AddCommand(userCmd, repoCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(red+"Error:", err)
	}
}
