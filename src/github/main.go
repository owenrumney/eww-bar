package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
)

type gitInfo struct {
	org          string
	stars        int
	issues       int
	pullRequests int
}

var githubIcon = "\uF09B"
var starIcon = "\uF005"
var issueIssue = "\uF192"

func main() {

	gitDeets := make(map[string]gitInfo)

	ctx := context.Background()

	if len(os.Args) < 1 {
		panic("Usage: github-info org:repo,org:repo,org:repo")
	}

	checkers := strings.Split(os.Args[1], ",")
	client := github.NewClient(getTokenContext(ctx))

	for _, toCheck := range checkers {
		parts := strings.Split(toCheck, ":")
		if len(parts) != 2 {
			continue
		}

		org := parts[0]
		repo := parts[1]

		repository, _, err := client.Repositories.Get(ctx, org, repo)
		if err != nil {
			panic(err)
		}

		gitDeets[repo] = gitInfo{
			org:    org,
			stars:  repository.GetStargazersCount(),
			issues: repository.GetOpenIssues(),
		}
	}

	fmt.Print(`(box :orientation "h" :space-evenly false`)
	for repo, info := range gitDeets {
		_ = repo
		_ = info
		fmt.Printf(`(button :onclick "xdg-open https://github.com/%s/%s" (label :text " %s: %s %d %s %d  | "))  `, info.org, repo, repo, starIcon, info.stars, issueIssue, info.issues)
	}
	fmt.Print(")")
}

func getTokenContext(ctx context.Context) *http.Client {
	githubToken, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		githubToken = tryTokenFile()
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	return tc
}

func tryTokenFile() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	tokenFile := filepath.Join(homeDir, ".token")
	token, err := os.ReadFile(tokenFile)
	if err != nil {
		panic("Did not find GITHUB_TOKEN")
	}

	return strings.TrimSpace(string(token))
}
