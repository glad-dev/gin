package issueList

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gn/issues/user"
	"gn/logger"
	"gn/remote"
	"gn/requests"
)

type queryAllGitHubResponse struct {
	Data struct {
		Repository struct {
			Issues struct {
				PageInfo pageInfo `json:"pageInfo"`
				Nodes    []struct {
					CreatedAt time.Time `json:"createdAt"`
					UpdatedAt time.Time `json:"updatedAt"`
					Title     string    `json:"title"`
					State     string    `json:"state"`
					Author    struct {
						Login string `json:"login"`
					} `json:"author"`
					Assignees struct {
						Nodes []struct {
							Login string `json:"login"`
						} `json:"nodes"`
					} `json:"assignees"`
					Number int `json:"number"`
				} `json:"nodes"`
			} `json:"issues"`
		} `json:"repository"`
	} `json:"data"`
}

type pageInfo struct {
	EndCursor   string `json:"endCursor"`
	HasNextPage bool   `json:"hasNextPage"`
}

var queryAllFirstRequest = `
	query($owner: String!, $name: String!) { 
	  repository(owner:$owner, name:$name) {
		issues(first: 100, orderBy:{field:CREATED_AT, direction:ASC}) {
		  pageInfo {
			endCursor
			hasNextPage
		  }
		  
		  nodes {
			title
			number
			state
			createdAt
			updatedAt
			author {
			  login
			  
			}
			
			assignees(first:100) {
			  nodes {
				login
			  }
			}
		  }
		}
	  }
	}
`

var queryAllFollowing = `
	query($owner: String!, $name: String!, $cursor: String!) { 
	  repository(owner:$owner, name:$name) {
		issues(first: 100, after: $cursor, orderBy:{field:CREATED_AT, direction:ASC}) {
		  pageInfo {
			endCursor
			hasNextPage
		  }
		  
		  nodes {
			title
			number
			state
			createdAt
			updatedAt
			author {
			  login
			  
			}
			
			assignees(first:100) {
			  nodes {
				login
			  }
			}
		  }
		}
	  }
	}
`

func QueryGitHub(match *remote.Match, projectPath string) ([]Issue, error) {
	tmp := strings.Split(projectPath, "/")
	if len(tmp) != 2 {
		logger.Log.Errorf("Project path is invalid: %s", projectPath)

		return nil, errors.New("invalid project path")
	}

	variables := map[string]string{
		"owner": tmp[0],
		"name":  tmp[1],
	}

	response, err := requests.Project(&requests.Query{
		Query:     queryAllFirstRequest,
		Variables: variables,
	}, match)
	if err != nil {
		return nil, fmt.Errorf("inital query failed: %w", err)
	}

	issueList, info, err := parseResponse(response, match.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if !info.HasNextPage {
		return issueList, nil
	}

	endCursor := info.EndCursor
	var newIssues []Issue
	for {
		variables["cursor"] = endCursor

		response, err = requests.Project(&requests.Query{
			Query:     queryAllFollowing,
			Variables: variables,
		}, match)
		if err != nil {
			return nil, fmt.Errorf("subsequent query failed: %w", err)
		}

		newIssues, info, err = parseResponse(response, match.Username)
		if err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}

		issueList = append(issueList, newIssues...)

		endCursor = info.EndCursor
		if !info.HasNextPage {
			break
		}
	}

	return issueList, nil
}

func parseResponse(response []byte, ownUsername string) ([]Issue, *pageInfo, error) {
	queryAll := queryAllGitHubResponse{}

	dec := json.NewDecoder(bytes.NewBuffer(response))
	dec.DisallowUnknownFields()
	err := dec.Decode(&queryAll)
	if err != nil {
		logger.Log.Errorf("Failed to decode the response: %s", err)

		return nil, nil, fmt.Errorf("unmarshle of issues failed: %w", err)
	}

	// Flatter the Graphql struct to an Issue struct
	var tmp Issue
	issueList := make([]Issue, 0)
	for _, issue := range queryAll.Data.Repository.Issues.Nodes {
		assignees := make([]user.Details, len(issue.Assignees.Nodes))
		for i := range issue.Assignees.Nodes {
			assignees[i] = user.Details{Username: issue.Assignees.Nodes[i].Login}
		}

		tmp = Issue{
			Title:     issue.Title,
			CreatedAt: issue.CreatedAt,
			UpdatedAt: issue.UpdatedAt,
			Iid:       strconv.Itoa(issue.Number),
			State:     issue.State,
			Assignees: assignees,
			Author:    user.Details{Username: issue.Author.Login},
		}

		tmp.UpdateUsername(ownUsername)

		issueList = append(issueList, tmp)
	}

	return issueList, &queryAll.Data.Repository.Issues.PageInfo, nil
}
