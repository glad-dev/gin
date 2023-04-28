package list

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"gn/logger"
	"gn/remote"

	"github.com/google/go-github/v52/github"
	"golang.org/x/oauth2"
)

// QueryGitHub returns all issues, open and closed, of a given repository.
func QueryGitHub(match *remote.Match, projectPath string) ([]Issue, error) {
	tmp := strings.Split(projectPath, "/")
	if len(tmp) != 2 {
		logger.Log.Errorf("Project path is invalid: %s", projectPath)

		return nil, errors.New("invalid project path")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: match.Token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	issueList := make([]Issue, 0)
	page := 0
	for {
		// owner is at tmp[0], repo name is at tmp[1]
		issues, response, err := client.Issues.ListByRepo(context.Background(), tmp[0], tmp[1], &github.IssueListByRepoOptions{
			State:     "all",
			Sort:      "created",
			Direction: "desc",
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: 100,
			},
		})
		if err != nil {
			return nil, fmt.Errorf("requesting issues: %w", err)
		}

		for _, issue := range issues {
			if issue.IsPullRequest() {
				continue
			}

			author := issue.GetUser()
			assigned := issue.GetAssignee()

			issueList = append(issueList, Issue{
				Title:     issue.GetTitle(),
				CreatedAt: issue.GetCreatedAt().Time,
				UpdatedAt: issue.GetUpdatedAt().Time,
				Iid:       strconv.Itoa(issue.GetNumber()),
				State:     issue.GetState(),
				Author: remote.User{
					Name:     author.GetName(),
					Username: author.GetLogin(),
				},
				Assignees: []remote.User{
					{
						Name:     assigned.GetName(),
						Username: assigned.GetLogin(),
					},
				},
			})
		}

		page = response.NextPage
		if page == 0 {
			break
		}
	}

	return issueList, nil
}
