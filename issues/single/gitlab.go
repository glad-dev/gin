package single

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"gn/issues/user"
	"gn/logger"
	"gn/remote"
	"gn/requests"
)

type querySingleGitLabResponse struct {
	Data struct {
		Project struct {
			Issue struct {
				Title       string       `json:"title"`
				Description string       `json:"description"`
				CreatedAt   time.Time    `json:"createdAt"`
				UpdatedAt   time.Time    `json:"updatedAt"`
				Author      user.Details `json:"author"`
				Assignees   struct {
					Nodes []user.Details `json:"nodes"`
				} `json:"assignees"`
				Labels struct {
					Nodes []Label `json:"nodes"`
				} `json:"labels"`
				Discussions struct {
					Nodes []struct {
						Notes struct {
							Nodes []struct {
								Body         string       `json:"body"`
								CreatedAt    time.Time    `json:"createdAt"`
								UpdatedAt    time.Time    `json:"updatedAt"`
								LastEditedBy user.Details `json:"lastEditedBy"`
								user.Details `json:"author"`
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

func QuerySingleGitLab(match *remote.Match, projectPath string, issueID string) (*IssueDetails, error) {
	variables := map[string]string{
		"projectPath": projectPath,
		"issueID":     issueID,
	}

	response, err := requests.Project(&requests.GraphqlQuery{
		Query:     querySingleGitLab,
		Variables: variables,
	}, match)
	if err != nil {
		return nil, fmt.Errorf("query single - request failed: %w", err)
	}

	if issueDoesNotExistGitLab(response) {
		logger.Log.Error("Requested issue does not exist.", "issueID", issueID, "response", string(response))

		return nil, ErrIssueDoesNotExist
	}

	querySingle := querySingleGitLabResponse{}

	dec := json.NewDecoder(bytes.NewBuffer(response))
	dec.DisallowUnknownFields()
	err = dec.Decode(&querySingle)
	if err != nil {
		logger.Log.Error("Failed to decode issue.", "error", err, "response", string(response))

		return nil, fmt.Errorf("unmarshal of issues failed: %w", err)
	}

	issueDetails := IssueDetails{
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
	assignees := make([]user.Details, 0)
	for _, assignee := range querySingle.Data.Project.Issue.Assignees.Nodes {
		assignees = append(assignees, user.Details{
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
			Author: user.Details{
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
				Author:       subComment.Details,
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

func issueDoesNotExistGitLab(response []byte) bool {
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
