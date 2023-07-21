package discussion

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/glad-dev/gin/logger"
	"github.com/glad-dev/gin/remote"
)

type querySingleGitLabResponse struct {
	Data struct {
		Project struct {
			Issue struct {
				Title       string      `json:"title"`
				Description string      `json:"description"`
				CreatedAt   time.Time   `json:"createdAt"`
				UpdatedAt   time.Time   `json:"updatedAt"`
				Author      remote.User `json:"author"`
				Assignees   struct {
					Nodes []remote.User `json:"nodes"`
				} `json:"assignees"`
				Labels struct {
					Nodes []Label `json:"nodes"`
				} `json:"labels"`
				Discussions struct {
					Nodes []struct {
						Notes struct {
							Nodes []struct {
								Body         string      `json:"body"`
								CreatedAt    time.Time   `json:"createdAt"`
								UpdatedAt    time.Time   `json:"updatedAt"`
								LastEditedBy remote.User `json:"lastEditedBy"`
								remote.User  `json:"author"`
								System       bool `json:"system"`
								Resolved     bool `json:"resolved"`
							} `json:"nodes"`
						} `json:"notes"`
					} `json:"nodes"`
				} `json:"discussions"`
			} `json:"issue"`
		} `json:"project"`
	} `json:"data"`
}

const querySingleGitLab = `
	query($projectPath: ID!, $issueID: String!) {
		project(fullPath: $projectPath) {
			issue(iid: $issueID) {
				title
				description
				createdAt
				updatedAt
				author {
					name
					username
				}
				assignees {
					nodes {
						name
						username
					}
				}
				labels {
					nodes {
						title
						color	
					}
				}
				discussions {
					nodes {
						notes {
							nodes {
								system
								author {
									name
									username
								}
								body
								createdAt
								updatedAt
								resolved
								lastEditedBy {
								name
								username
								}
							}
						}
					}
				}
			}
		}
	}
`

// QueryGitLab returns the discussion associated with the passed issueID. If the requested issue does not exist, an
// ErrIssueDoesNotExist is returned.
func QueryGitLab(match *remote.Match, projectPath string, issueID string) (*Details, error) {
	variables := map[string]string{
		"projectPath": projectPath,
		"issueID":     issueID,
	}

	response, err := graphQLRequest(&query{
		Query:     querySingleGitLab,
		Variables: variables,
	}, match)
	if err != nil {
		return nil, fmt.Errorf("query single - request failed: %w", err)
	}

	if issueDoesNotExist(response) {
		logger.Log.Error("Requested discussion does not exist.", "issueID", issueID, "response", string(response))

		return nil, ErrIssueDoesNotExist
	}

	querySingle := querySingleGitLabResponse{}

	dec := json.NewDecoder(bytes.NewBuffer(response))
	dec.DisallowUnknownFields()
	err = dec.Decode(&querySingle)
	if err != nil {
		logger.Log.Error("Failed to decode discussion.", "error", err, "response", string(response))

		return nil, fmt.Errorf("unmarshal of issues failed: %w", err)
	}

	issueDetails := Details{
		Title:       querySingle.Data.Project.Issue.Title,
		Description: querySingle.Data.Project.Issue.Description,
		CreatedAt:   querySingle.Data.Project.Issue.CreatedAt,
		UpdatedAt:   querySingle.Data.Project.Issue.UpdatedAt,
		Author:      querySingle.Data.Project.Issue.Author,
		BaseURL:     match.URL,

		Assignees:  nil,
		Labels:     nil,
		Discussion: nil,
	}

	// Flatten response
	// Assignees
	assignees := make([]remote.User, 0)
	for _, assignee := range querySingle.Data.Project.Issue.Assignees.Nodes {
		assignees = append(assignees, remote.User{
			Name:     assignee.Name,
			Username: assignee.Username,
		})
	}
	issueDetails.Assignees = assignees

	// Labels
	labels := make([]Label, 0)
	for _, label := range querySingle.Data.Project.Issue.Labels.Nodes {
		labels = append(labels, Label{
			Title: label.Title,
			Color: label.Color,
		})
	}
	issueDetails.Labels = labels

	// Discussion
	for _, node := range querySingle.Data.Project.Issue.Discussions.Nodes {
		inner := node.Notes.Nodes
		if len(inner) == 0 {
			logger.Log.Info("Discussion without nodes", "response", string(response))

			continue
		}

		if inner[0].System {
			continue
		}

		comment := Comment{
			Author: remote.User{
				Name:     inner[0].Name,
				Username: inner[0].Username,
			},
			Body:         inner[0].Body,
			CreatedAt:    inner[0].CreatedAt,
			UpdatedAt:    inner[0].UpdatedAt,
			Resolved:     inner[0].Resolved,
			LastEditedBy: inner[0].LastEditedBy,
			Comments:     make([]Comment, 0),
		}

		// Get sub comments
		for _, subComment := range inner[1:] {
			comment.Comments = append(comment.Comments, Comment{
				Author:       subComment.User,
				Body:         subComment.Body,
				CreatedAt:    subComment.CreatedAt,
				UpdatedAt:    subComment.UpdatedAt,
				Resolved:     subComment.Resolved,
				LastEditedBy: subComment.LastEditedBy,
				Comments:     nil,
			})
		}

		issueDetails.Discussion = append(issueDetails.Discussion, comment)
	}

	issueDetails.UpdateUsername(match.Username)

	return &issueDetails, nil
}

func issueDoesNotExist(response []byte) bool {
	emptyResponse := struct {
		Data struct {
			Project struct {
				Issue interface{} `json:"issue"`
			} `json:"project"`
		} `json:"data"`
	}{}

	dec := json.NewDecoder(bytes.NewBuffer(response))
	dec.DisallowUnknownFields()
	err := dec.Decode(&emptyResponse)
	if err != nil {
		return false
	}

	return !reflect.ValueOf(emptyResponse.Data.Project.Issue).IsValid()
}