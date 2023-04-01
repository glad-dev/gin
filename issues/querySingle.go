package issues

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"time"

	"gn/config"
	"gn/repo"
	"gn/requests"
)

var ErrIssueDoesNotExist = errors.New("issue with the given iid does not exist")

type querySingleResponse struct {
	Data struct {
		Project struct {
			Issue struct {
				Title       string    `json:"title"`
				Description string    `json:"description"`
				CreatedAt   time.Time `json:"createdAt"`
				UpdatedAt   time.Time `json:"updatedAt"`
				Author      User      `json:"author"`
				Assignees   struct {
					Nodes []User `json:"nodes"`
				} `json:"assignees"`
				Labels struct {
					Nodes []Label `json:"nodes"`
				} `json:"labels"`
				Discussions struct {
					Nodes []struct {
						Notes struct {
							Nodes []struct {
								Body         string    `json:"body"`
								CreatedAt    time.Time `json:"createdAt"`
								UpdatedAt    time.Time `json:"updatedAt"`
								LastEditedBy User      `json:"lastEditedBy"`
								User         `json:"author"`
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

func QuerySingle(config *config.Wrapper, details []repo.Details, issueID string) (*IssueDetails, error) {
	query := `
		query($projectPath: ID!, $issueIID: String!) {
		  project(fullPath: $projectPath) {
			issue(iid: $issueIID) {
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

	lab, projectPath, err := config.GetMatchingConfig(details)
	if err != nil {
		return nil, err
	}

	variables := map[string]string{
		"projectPath": projectPath,
		"issueIID":    issueID,
	}

	response, err := requests.Project(&requests.GraphqlQuery{
		Query:     query,
		Variables: variables,
	}, lab)

	if err != nil {
		return nil, fmt.Errorf("query single - request failed: %w", err)
	}

	if issueDoesNotExist(response) {
		return nil, ErrIssueDoesNotExist
	}

	querySingle := querySingleResponse{}
	err = json.Unmarshal(response, &querySingle)
	if err != nil {
		return nil, fmt.Errorf("unmarshal of issues failed: %w", err)
	}

	issueDetails := IssueDetails{
		Title:       querySingle.Data.Project.Issue.Title,
		Description: querySingle.Data.Project.Issue.Description,
		CreatedAt:   querySingle.Data.Project.Issue.CreatedAt,
		UpdatedAt:   querySingle.Data.Project.Issue.UpdatedAt,
		Author:      querySingle.Data.Project.Issue.Author,
		BaseURL:     lab.URL,

		Assignees:  nil,
		Labels:     nil,
		Discussion: nil,
	}

	// Flatten response
	// Assignees
	assignees := make([]User, 0)
	for _, assignee := range querySingle.Data.Project.Issue.Assignees.Nodes {
		assignees = append(assignees, User{
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
			log.Printf("query single - discussion without nodes?\n%s", response)

			continue
		}

		if inner[0].System {
			continue
		}

		comment := Comment{
			Author: User{
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

	issueDetails.UpdateUsername(lab.Username)

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

	err := json.Unmarshal(response, &emptyResponse)
	if err != nil {
		return false
	}

	return !reflect.ValueOf(emptyResponse.Data.Project.Issue).IsValid()
}
