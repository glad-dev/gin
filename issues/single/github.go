package single

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"gn/config/remote"
	"gn/issues/user"
	"gn/logger"
	"gn/requests"
)

var ErrIssueDoesNotExist = errors.New("issue with the given iid does not exist")

type querySingleGitHubResponse struct {
	Data struct {
		Repository struct {
			Issue struct {
				Title     string    `json:"title"`
				Body      string    `json:"body"`
				CreatedAt time.Time `json:"createdAt"`
				UpdatedAt time.Time `json:"updatedAt"`
				Author    struct {
					Login string `json:"login"`
				} `json:"author"`
				Assignees struct {
					Nodes []struct {
						Login string `json:"login"`
					} `json:"nodes"`
				} `json:"assignees"`
				Labels struct {
					Nodes []struct {
						Color string `json:"color"`
						Name  string `json:"name"`
					} `json:"nodes"`
				} `json:"labels"`
				Comments struct {
					PageInfo pageInfo `json:"pageInfo"`
					Nodes    []struct {
						CreatedAt    time.Time  `json:"createdAt"`
						UpdatedAt    time.Time  `json:"updatedAt"`
						LastEditedAt *time.Time `json:"lastEditedAt"`
						Author       struct {
							Login string `json:"login"`
						} `json:"author"`
						Body string `json:"body"`
					} `json:"nodes"`
				} `json:"comments"`
			} `json:"issue"`
		} `json:"repository"`
	} `json:"data"`
}

type pageInfo struct {
	EndCursor   string `json:"endCursor"`
	HasNextPage bool   `json:"hasNextPage"`
}

const querySingleGitHubFirst = `
	query($owner: String!, $name: String!, $issueID: Int!) {
	  repository(owner: $owner, name: $name) {
		issue(number: $issueID) {
		  title
		  body
		  createdAt
		  updatedAt
		  author {
			login
		  }
		  assignees(first: 100) {
			nodes {
			  login
			}
		  }
		  labels(first: 100) {
			nodes {
			  color
			  name
			}
		  }
		  comments(first: 100) {
			pageInfo {
			  hasNextPage
			  endCursor
			}
			nodes {
			  
			  author {
				login
			  }
			  body
			  createdAt
			  updatedAt
			  lastEditedAt
			}
		  }
		}
	  }
	}
`

const querySingleGitHubFollowing = `
	query($owner: String!, $name: String!, $cursor: String, $issueID: Int!) {
	  repository(owner: $owner, name: $name) {
		issue(number: $issueID) {
		  comments(first: 100, after: $cursor) {
			pageInfo {
			  hasNextPage
			  endCursor
			}
			nodes {
			  
			  author {
				login
			  }
			  body
			  createdAt
			  updatedAt
			  lastEditedAt
			}
		  }
		}
	  }
	}
`

func QuerySingleGitHub(match *remote.Match, projectPath string, issueID string) (*IssueDetails, error) {
	tmp := strings.Split(projectPath, "/")
	if len(tmp) != 2 {
		logger.Log.Errorf("Project path is invalid: %s", projectPath)

		return nil, errors.New("invalid project path")
	}

	issueNumber, err := strconv.Atoi(issueID)
	if err != nil {
		logger.Log.Error("Failed to convert issueID to int.", "error", err, "issueID", issueID)

		return nil, fmt.Errorf("failed to convert issueID to int: %w", err)
	}

	variables := map[string]interface{}{
		"owner":   tmp[0],
		"name":    tmp[1],
		"issueID": issueNumber,
	}

	responseReader, err := requests.GitHubComment(&requests.GitHubCommentQuery{
		Query:     querySingleGitHubFirst,
		Variables: variables,
	}, match)
	if err != nil {
		return nil, fmt.Errorf("query single - request failed: %w", err)
	}

	response, err := io.ReadAll(responseReader)
	if err != nil {
		logger.Log.Errorf("Failed to read response: %s", err)

		return nil, fmt.Errorf("query single - failed to read request: %w", err)
	}

	if issueDoesNotExistGitHub(bytes.NewBuffer(response)) {
		logger.Log.Error("Requested issue does not exist.", "issueID", issueID, "response", string(response))

		return nil, ErrIssueDoesNotExist
	}

	querySingle := querySingleGitHubResponse{}

	dec := json.NewDecoder(bytes.NewBuffer(response))
	dec.DisallowUnknownFields()
	err = dec.Decode(&querySingle)
	if err != nil {
		logger.Log.Error("Failed to decode issue.", "error", err, "response", string(response))

		return nil, fmt.Errorf("unmarshal of issues failed: %w", err)
	}

	issueDetails := IssueDetails{
		Title:       querySingle.Data.Repository.Issue.Title,
		Description: querySingle.Data.Repository.Issue.Body,
		CreatedAt:   querySingle.Data.Repository.Issue.CreatedAt,
		UpdatedAt:   querySingle.Data.Repository.Issue.UpdatedAt,
		Author:      user.Details{Username: querySingle.Data.Repository.Issue.Author.Login},
		BaseURL:     match.URL,

		Assignees:  nil,
		Labels:     nil,
		Discussion: nil,
	}

	// Flatten response
	// Assignees
	assignees := make([]user.Details, 0)
	for _, assignee := range querySingle.Data.Repository.Issue.Assignees.Nodes {
		assignees = append(assignees, user.Details{
			Username: assignee.Login,
		})
	}
	issueDetails.Assignees = assignees

	// Labels
	labels := make([]Label, 0)
	for _, label := range querySingle.Data.Repository.Issue.Labels.Nodes {
		labels = append(labels, Label{
			Title: label.Name,
			Color: label.Color,
		})
	}
	issueDetails.Labels = labels

	// Parse initial comments
	comments, info, err := parseComments(response)
	if err != nil {
		return nil, err
	}

	issueDetails.Discussion = append(issueDetails.Discussion, comments...)

	if info.HasNextPage {
		endCursor := querySingle.Data.Repository.Issue.Comments.PageInfo.EndCursor

		for {
			variables["cursor"] = endCursor

			responseReader, err = requests.GitHubComment(&requests.GitHubCommentQuery{
				Query:     querySingleGitHubFollowing,
				Variables: variables,
			}, match)
			if err != nil {
				return nil, fmt.Errorf("query single - request failed: %w", err)
			}

			response, err = io.ReadAll(responseReader)
			if err != nil {
				logger.Log.Errorf("Failed to read response: %s", err)

				return nil, fmt.Errorf("query single - failed to read request: %w", err)
			}

			comments, info, err = parseComments(response)
			if err != nil {
				return nil, err
			}

			issueDetails.Discussion = append(issueDetails.Discussion, comments...)

			endCursor = info.EndCursor
			if !info.HasNextPage {
				break
			}
		}
	}

	issueDetails.UpdateUsername(match.Username)

	return &issueDetails, nil
}

func parseComments(response []byte) ([]Comment, *pageInfo, error) {
	querySingle := querySingleGitHubResponse{}

	dec := json.NewDecoder(bytes.NewBuffer(response))
	dec.DisallowUnknownFields()
	err := dec.Decode(&querySingle)
	if err != nil {
		logger.Log.Error("Failed to decode issue.", "error", err, "response", string(response))

		return nil, nil, fmt.Errorf("unmarshal of issues failed: %w", err)
	}

	comments := make([]Comment, len(querySingle.Data.Repository.Issue.Comments.Nodes))
	for i, node := range querySingle.Data.Repository.Issue.Comments.Nodes {
		/*
			if inner[0].System {
				continue
			}*/

		comments[i] = Comment{
			Author: user.Details{
				Username: node.Author.Login,
			},
			Body:         node.Body,
			CreatedAt:    node.CreatedAt,
			UpdatedAt:    node.UpdatedAt,
			Resolved:     false,
			LastEditedBy: user.Details{Username: "unknown"},
			Comments:     make([]Comment, 0),
		}
	}

	return comments, &querySingle.Data.Repository.Issue.Comments.PageInfo, nil
}

func issueDoesNotExistGitHub(response io.Reader) bool {
	emptyResponse := struct {
		Data struct {
			Repository interface{} `json:"repository"`
		} `json:"data"`
	}{}

	dec := json.NewDecoder(response)
	dec.DisallowUnknownFields()
	err := dec.Decode(&emptyResponse)
	if err != nil {
		return false
	}

	return emptyResponse.Data.Repository == nil
}
