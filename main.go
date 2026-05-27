package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

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

type GitHubUser struct {
	Login       string `json:"login"`
	Name        string `json:"name"`
	PublicRepos int    `json:"public_repos"`
	Followers   int    `json:"followers"`
	HTMLURL     string `json:"html_url"`
}

func fetchGitHubUser(username string) (*GitHubUser, error) {

	apiURL := "https://api.github.com/users/" + url.PathEscape(username)

	request, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("User-Agent", "go-api-practice")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	response, err := client.Do(request)
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

	// This part checks to see if the user passed in any arguments
	if len(os.Args) < 2 {
		fmt.Println(yellow + "Please provide a GitHub username" + reset)
		return
	}

	// The username will be the argument provided by the user
	username := os.Args[1]
	user, err := fetchGitHubUser(username)
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

}
