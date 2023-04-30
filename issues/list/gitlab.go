package list

import (
	"fmt"
	"io"
	"strconv"

	"gn/logger"
	"gn/remote"

	"github.com/xanzy/go-gitlab"
)

// QueryGitLab returns all issues, open and closed, of a given repository.
func QueryGitLab(match *remote.Match, projectPath string) ([]Issue, error) {
	client, err := gitlab.NewClient(match.Token, gitlab.WithBaseURL(match.ApiURL()))
	if err != nil {
		logger.Log.Error("Creating gitlab client",
			"error", err,
			"API-URL", match.ApiURL(),
			"project-path", projectPath,
		)

		return nil, fmt.Errorf("creating gitlab client: %w", err)
	}

	state := "all"
	orderBy := "created_at"
	sort := "desc"

	page := 0
	issueList := make([]Issue, 0)
	for {
		issues, resp, err := client.Issues.ListProjectIssues(projectPath, &gitlab.ListProjectIssuesOptions{
			ListOptions: gitlab.ListOptions{
				Page:    page,
				PerPage: 100,
			},
			State:   &state,
			OrderBy: &orderBy,
			Sort:    &sort,
		})
		if err != nil {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				body = []byte("Failed to read body")
			}

			logger.Log.Error("Requesting issues",
				"error", err,
				"API-URL", match.ApiURL(),
				"project-path", projectPath,
				"response-body", string(body),
			)

			return nil, fmt.Errorf("requesting issues: %w", err)
		}

		for _, issue := range issues {
			assignees := make([]remote.User, len(issue.Assignees))
			for i, assignee := range issue.Assignees {
				assignees[i] = remote.User{
					Name:     assignee.Name,
					Username: assignee.Username,
				}
			}

			issueList = append(issueList, Issue{
				Title:     issue.Title,
				CreatedAt: *issue.CreatedAt,
				UpdatedAt: *issue.UpdatedAt,
				Iid:       strconv.Itoa(issue.IID),
				State:     issue.State,
				Author: remote.User{
					Name:     issue.Author.Name,
					Username: issue.Author.Username,
				},
				Assignees: assignees,
			})
		}

		page = resp.NextPage
		if page == 0 {
			break
		}
	}

	return issueList, nil
}
