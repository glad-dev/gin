package discussion

import (
	"context"
	"fmt"
	"time"

	"github.com/glad-dev/gin/logger"
	"github.com/glad-dev/gin/remote"
	"github.com/glad-dev/gin/remote/match"
	"github.com/shurcooL/graphql"
)

type query struct {
	Project struct {
		Issue struct {
			Title       graphql.String
			Description graphql.String
			CreatedAt   graphql.String
			UpdatedAt   graphql.String
			Author      struct {
				Name     graphql.String
				Username graphql.String
			}
			Assignees struct {
				Nodes []struct {
					Name     graphql.String
					Username graphql.String
				}
			}
			Labels struct {
				Nodes []struct {
					Title graphql.String
					Color graphql.String
				}
			}
			Discussions struct {
				Nodes []struct {
					Notes struct {
						Nodes []struct {
							Author struct {
								Name     graphql.String
								Username graphql.String
							}
							LastEditedBy struct {
								Name     graphql.String
								Username graphql.String
							}
							Body      graphql.String
							CreatedAt graphql.String
							UpdatedAt graphql.String
							System    graphql.Boolean
							Resolved  graphql.Boolean
						}
					}
				}
			}
		} `graphql:"issue(iid: $issueID)"`
	} `graphql:"project(fullPath: $projectPath)"`
}

const timeLayout = "2006-01-02T15:04:05Z"

// QueryGitLab returns the discussion associated with the passed issueID. If the requested issue does not exist, an
// ErrIssueDoesNotExist is returned.
func QueryGitLab(match *match.Match, projectPath string, issueID string) (*Details, error) {
	client, err := match.GraphqlClient()
	if err != nil {
		// No need to log, since match.GraphqlClient() already logs the error
		return nil, err
	}

	q := &query{}

	err = client.Query(context.Background(), q, map[string]any{
		"projectPath": projectPath,
		"issueID":     graphql.String(issueID),
	})

	if err != nil {
		logger.Log.Error("Requesting discussion", "error", err, "projectPath", projectPath, "issueID", issueID)

		return nil, fmt.Errorf("requesting discussion: %w", err)
	}

	creationTime, err := time.Parse(timeLayout, string(q.Project.Issue.CreatedAt))
	if err != nil {
		logger.Log.Warn("failed to parse creation time", "time", string(q.Project.Issue.CreatedAt), "error", err)
	}

	updateTime, err := time.Parse(timeLayout, string(q.Project.Issue.UpdatedAt))
	if err != nil {
		logger.Log.Warn("failed to parse update time", "time", string(q.Project.Issue.UpdatedAt), "error", err)
	}

	issueDetails := Details{
		Title:       string(q.Project.Issue.Title),
		Description: string(q.Project.Issue.Description),
		CreatedAt:   creationTime,
		UpdatedAt:   updateTime,
		Author: remote.User{
			Name:     string(q.Project.Issue.Author.Name),
			Username: string(q.Project.Issue.Author.Username),
		},
		BaseURL: match.URL,

		Assignees:  nil,
		Labels:     nil,
		Discussion: nil,
	}

	// Flatten response
	// Assignees
	assignees := make([]remote.User, 0)
	for _, assignee := range q.Project.Issue.Assignees.Nodes {
		assignees = append(assignees, remote.User{
			Name:     string(assignee.Name),
			Username: string(assignee.Username),
		})
	}
	issueDetails.Assignees = assignees

	// Labels
	labels := make([]Label, 0)
	for _, label := range q.Project.Issue.Labels.Nodes {
		labels = append(labels, Label{
			Title: string(label.Title),
			Color: string(label.Color),
		})
	}
	issueDetails.Labels = labels

	// Discussion
	for _, node := range q.Project.Issue.Discussions.Nodes {
		inner := node.Notes.Nodes
		if len(inner) == 0 {
			logger.Log.Info("Discussion without nodes", "response", q)

			continue
		}

		if inner[0].System {
			continue
		}

		creationTime, err = time.Parse(timeLayout, string(inner[0].CreatedAt))
		if err != nil {
			logger.Log.Warn("failed to parse creation time", "time", string(inner[0].CreatedAt), "error", err)
		}

		updateTime, err = time.Parse(timeLayout, string(inner[0].UpdatedAt))
		if err != nil {
			logger.Log.Warn("failed to parse update time", "time", string(inner[0].UpdatedAt), "error", err)
		}

		comment := Comment{
			Author: remote.User{
				Name:     string(inner[0].Author.Name),
				Username: string(inner[0].Author.Username),
			},
			Body:      string(inner[0].Body),
			CreatedAt: creationTime,
			UpdatedAt: updateTime,
			Resolved:  bool(inner[0].Resolved),
			LastEditedBy: remote.User{
				Name:     string(inner[0].LastEditedBy.Name),
				Username: string(inner[0].LastEditedBy.Username),
			},
			Comments: make([]Comment, 0),
		}

		// Get sub comments
		for _, subComment := range inner[1:] {
			creationTime, err = time.Parse(timeLayout, string(subComment.CreatedAt))
			if err != nil {
				logger.Log.Warn("failed to parse creation time", "time", string(subComment.CreatedAt), "error", err)
			}

			updateTime, err = time.Parse(timeLayout, string(subComment.UpdatedAt))
			if err != nil {
				logger.Log.Warn("failed to parse update time", "time", string(subComment.UpdatedAt), "error", err)
			}

			comment.Comments = append(comment.Comments, Comment{
				Author: remote.User{
					Name:     string(subComment.Author.Name),
					Username: string(subComment.Author.Username),
				},
				Body:      string(subComment.Body),
				CreatedAt: creationTime,
				UpdatedAt: updateTime,
				Resolved:  bool(subComment.Resolved),
				LastEditedBy: remote.User{
					Name:     string(subComment.Author.Name),
					Username: string(subComment.Author.Username),
				},
				Comments: nil,
			})
		}

		issueDetails.Discussion = append(issueDetails.Discussion, comment)
	}

	issueDetails.UpdateUsername(match.Username)

	return &issueDetails, nil
}
