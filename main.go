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

func main() {
	// Prints a friendy message to user...
	fmt.Println(yellow + "GitHub user lookup..." + reset)

	rootCmd := &cobra.Command{
		Use:   "github-tool",
		Short: "Look up GitHub inofrmation from the command line",
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

	rootCmd.AddCommand(userCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(red+"Error:", err)
	}
}
