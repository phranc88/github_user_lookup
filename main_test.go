package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFetchUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{
			"login": "octocat",
			"name": "The Octocat",
			"public_repos": 8,
			"followers": 100,
			"html_url": "https://github.com/octocat"
		}`)
	}))
	defer server.Close()

	client := GitHubClient{
		httpClient: server.Client(),
		baseURL:    server.URL,
	}

	user, err := client.fetchUser("octocat")
	if err != nil {
		t.Fatal(err)
	}

	if user.Login != "octocat" {
		t.Errorf("expected login %q, got %q", "octocat", user.Login)
	}
}

func TestFetchUserNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", http.StatusNotFound)
	}))
	defer server.Close()

	client := GitHubClient{
		httpClient: server.Client(),
		baseURL:    server.URL,
	}

	user, err := client.fetchUser("missing-user")

	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	if !strings.Contains(err.Error(), `user "missing-user" was not found`) {
		t.Errorf("expected not found error, got %q", err.Error())
	}

	if user != nil {
		t.Errorf("expected user to be nil, got %#v", user)
	}
}

func TestFetchRepo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repos/octocat/Hello-World" {
			t.Errorf("expected path %q, got %q", "/repos/octocat/Hello-World", r.URL.Path)
		}

		fmt.Fprint(w, `{
			"name": "Hello-World",
			"full_name": "octocat/Hello-World",
			"description": "My first repository on GitHub!",
			"stargazers_count": 3000,
			"forks_count": 2500,
			"language": "Go",
			"html_url": "https://github.com/octocat/Hello-World"
		}`)
	}))
	defer server.Close()

	client := GitHubClient{
		httpClient: server.Client(),
		baseURL:    server.URL,
	}

	repo, err := client.fetchRepo("octocat", "Hello-World")
	if err != nil {
		t.Fatal(err)
	}

	if repo.FullName != "octocat/Hello-World" {
		t.Errorf("expected full name %q, got %q", "octocat/Hello-World", repo.FullName)
	}

	if repo.Stars != 3000 {
		t.Errorf("expected stars %d, got %d", 3000, repo.Stars)
	}
}

func TestFetchRepoNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", http.StatusNotFound)
	}))
	defer server.Close()

	client := GitHubClient{
		httpClient: server.Client(),
		baseURL:    server.URL,
	}

	repo, err := client.fetchRepo("octocat", "missing-repo")

	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	if !strings.Contains(err.Error(), `repo "octocat/missing-repo" was not found`) {
		t.Errorf("expected not found error, got %q", err.Error())
	}

	if repo != nil {
		t.Errorf("expected repo to be nil, got %#v", repo)
	}
}
