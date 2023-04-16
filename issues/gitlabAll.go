package issues

import (
	"encoding/json"
	"fmt"
	"time"

	"gn/config"
	"gn/logger"
	"gn/requests"
)

type queryAllGitlab struct {
	Data struct {
		Project struct {
			Issues struct {
				PageInfo struct {
					EndCursor   string `json:"endCursor"`
					HasNextPage bool   `json:"hasNextPage"`
				} `json:"pageInfo"`
				Nodes []struct {
					Title     string    `json:"title"`
					CreatedAt time.Time `json:"createdAt"`
					UpdatedAt time.Time `json:"updatedAt"`
					Iid       string    `json:"iid"`
					State     string    `json:"state"`
					Author    User      `json:"author"`
					Assignees struct {
						Nodes []User `json:"nodes"`
					} `json:"assignees"`
				} `json:"nodes"`
			} `json:"issues"`
		} `json:"project"`
	} `json:"data"`
}

const queryAllQuery = `
		query($projectPath: ID!, $cursor: String) {
		  project(fullPath: $projectPath) {
			issues(first: 100, after: $cursor, sort: CREATED_ASC) {
			  pageInfo {
			    endCursor
			    hasNextPage
			  }
			  nodes {
				title
				createdAt
				updatedAt
				iid
				state
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
			  }
			}
		  }
		}
	`

func gitlabQueryAll(match *config.Match, projectPath string) ([]Issue, error) {
	endCursor := ""
	issueList := make([]Issue, 0)
	variables := map[string]string{
		"projectPath": projectPath,
	}

	for {
		variables["cursor"] = endCursor

		response, err := requests.Project(&requests.GraphqlQuery{
			Query:     queryAllQuery,
			Variables: variables,
		}, match)
		if err != nil {
			return nil, fmt.Errorf("query all issues failed: %w", err)
		}

		queryAll := queryAllGitlab{}

		dec := json.NewDecoder(response)
		dec.DisallowUnknownFields()
		err = dec.Decode(&queryAll)
		if err != nil {
			logger.Log.Errorf("Failed to decode the response: %s", err)

			return nil, fmt.Errorf("unmarshle of issues failed: %w", err)
		}

		// Flatter the Graphql struct to an Issue struct
		var tmp Issue
		for _, issue := range queryAll.Data.Project.Issues.Nodes {
			tmp = Issue{
				Title:     issue.Title,
				CreatedAt: issue.CreatedAt,
				UpdatedAt: issue.UpdatedAt,
				Iid:       issue.Iid,
				State:     issue.State,
				Assignees: issue.Assignees.Nodes,
				Author:    issue.Author,
			}

			tmp.UpdateUsername(match.Username)

			issueList = append(issueList, tmp)
		}

		endCursor = queryAll.Data.Project.Issues.PageInfo.EndCursor
		if !queryAll.Data.Project.Issues.PageInfo.HasNextPage {
			break
		}
	}

	return issueList, nil
}
