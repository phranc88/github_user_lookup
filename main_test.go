package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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

	if user != nil {
		t.Errorf("expected user to be nil, got %#v", user)
	}
}
