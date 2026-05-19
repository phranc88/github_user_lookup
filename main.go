package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
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

func main() {
	fmt.Println(yellow + "GitHub user lookup..." + reset)

	if len(os.Args) < 2 {
		fmt.Println(yellow + "Please provide a GitHub username" + reset)
		return
	}

	username := os.Args[1]
	url := "https://api.github.com/users/" + username

	response, err := http.Get(url)

	if err != nil {
		fmt.Println(red+"Request failed..."+reset, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Println("GitHub returned:", yellow+response.Status+reset)
		return
	}
	fmt.Println("Status:", green+response.Status+reset)

	// body, err := io.ReadAll(response.Body)
	// if err != nil {
	// 	fmt.Println(red+"Could not read response:"+reset, err)
	// 	return
	// }

	var user GitHubUser

	err = json.NewDecoder(response.Body).Decode(&user)
	if err != nil {
		fmt.Println(red+"Could not parse JSON..."+reset, err)
		return
	}

	fmt.Println("Login:", green+user.Login+reset)
	fmt.Println("Name:", green+user.Name+reset)
	fmt.Println("Public repos:", orange+strconv.Itoa(user.PublicRepos)+reset)
	fmt.Println("Followers:", magenta+strconv.Itoa(user.Followers)+reset)
	fmt.Println("Profile:", cyan+user.HTMLURL+reset)

}
