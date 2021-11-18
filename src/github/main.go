package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/google/go-github/v39/github"
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
	client := github.NewClient(http.DefaultClient)

	for _, toCheck := range checkers {
		parts := strings.Split(toCheck, ":")
		if len(parts) != 2 {
			continue
		}

		org := parts[0]
		repo := parts[1]

		repositories, _, err := client.Search.Repositories(ctx, fmt.Sprintf("%s/%s", org, repo), nil)
		if err != nil {
			panic(err)
		}

		for _, repository := range repositories.Repositories {
			if repository.GetName() == repo {
				gitDeets[repo] = gitInfo{
					org:    org,
					stars:  repository.GetStargazersCount(),
					issues: repository.GetOpenIssues(),
				}
			}
		}

	}

	fmt.Print(`(box :orientation "h" :space-evenly false`)
	for repo, info := range gitDeets {
		_ = repo
		_ = info
		fmt.Printf(`(button :onclick "eww update revealGithub=false && xdg-open https://github.com/%s/%s" (label :text " %s: %s %d %s %d  | "))  `, info.org, repo, repo, starIcon, info.stars, issueIssue, info.issues)
	}
	fmt.Print(")")
}
