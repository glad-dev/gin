package list

import (
	"context"
	"fmt"
	"time"

	"github.com/glad-dev/gin/log"
	"github.com/glad-dev/gin/remote"
	"github.com/glad-dev/gin/remote/match"

	"github.com/shurcooL/graphql"
)

type query struct {
	Project struct {
		Issues struct {
			PageInfo struct {
				EndCursor   graphql.String
				HasNextPage graphql.Boolean
			}

			Nodes []struct {
				Title     graphql.String
				CreatedAt graphql.String
				UpdatedAt graphql.String
				Iid       graphql.String
				State     graphql.String
				Author    struct {
					Name     graphql.String
					Username graphql.String
				}
				Assignees struct {
					Nodes []struct {
						Name     graphql.String
						Username graphql.String
					}
				} `graphql:"assignees (first: 100)"`
			}
		} `graphql:"issues(first: 100, after: $cursor, sort:CREATED_DESC)"`
	} `graphql:"project(fullPath: $path)"`
}

// QueryGitLab returns all issues, open and closed, of a given repository.
func QueryGitLab(match *match.Match, projectPath string, channel chan int) ([]Issue, error) {
	client, err := match.GraphqlClient()
	if err != nil {
		// No need to log, since match.GraphqlClient() already logs the error
		return nil, err
	}

	var cursor graphql.String
	q := &query{}
	issueList := make([]Issue, 0)
	for {
		err = client.Query(context.Background(), q, map[string]any{
			"path":   projectPath,
			"cursor": cursor,
		})
		if err != nil {
			log.Error("Requesting issues", "error", err, "projectPath", projectPath)

			return nil, fmt.Errorf("requesting issues: %w", err)
		}

		for _, issue := range q.Project.Issues.Nodes {
			assignees := make([]remote.User, len(issue.Assignees.Nodes))
			for i, assignee := range issue.Assignees.Nodes {
				assignees[i] = remote.User{
					Name:     string(assignee.Name),
					Username: string(assignee.Username),
				}
			}

			creationTime, err := time.Parse(timeLayout, string(issue.CreatedAt))
			if err != nil {
				log.Warn("failed to parse creation time", "time", string(issue.CreatedAt), "error", err)
			}

			updateTime, err := time.Parse(timeLayout, string(issue.UpdatedAt))
			if err != nil {
				log.Warn("failed to parse update time", "time", string(issue.UpdatedAt), "error", err)
			}

			issueList = append(issueList, Issue{
				Title:     string(issue.Title),
				CreatedAt: creationTime,
				UpdatedAt: updateTime,
				Iid:       string(issue.Iid),
				State:     string(issue.State),
				Author: remote.User{
					Name:     string(issue.Author.Name),
					Username: string(issue.Author.Username),
				},
				Assignees: assignees,
			})
		}

		sendCountUpdate(channel, len(issueList))

		cursor = q.Project.Issues.PageInfo.EndCursor
		if !q.Project.Issues.PageInfo.HasNextPage {
			break
		}
	}

	return issueList, nil
}
