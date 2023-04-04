package issues

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"gn/config"
	"gn/repo"
	"gn/requests"
)

type queryAllResponse struct {
	Data struct {
		Project struct {
			Issues struct {
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
		query($projectPath: ID!) {
		  project(fullPath: $projectPath) {
			issues(sort: CREATED_ASC) {
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

func QueryAll(conf *config.Wrapper, details []repo.Details, u *url.URL) ([]Issue, error) {
	// TODO: Query all issues, not just the first 100
	lab, projectPath, err := getMatchingConfig(conf, details, u)
	if err != nil {
		return nil, err
	}

	variables := map[string]string{
		"projectPath": projectPath,
	}

	response, err := requests.Project(&requests.GraphqlQuery{
		Query:     queryAllQuery,
		Variables: variables,
	}, lab)

	if err != nil {
		return nil, fmt.Errorf("query all issues failed: %w", err)
	}

	queryAll := queryAllResponse{}
	err = json.Unmarshal(response, &queryAll)
	if err != nil {
		return nil, fmt.Errorf("unmarshle of issues failed: %w", err)
	}

	// Flatter the Graphql struct to an Issue struct
	issues := make([]Issue, 0)
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

		tmp.UpdateUsername(lab.Username)

		issues = append(issues, tmp)
	}

	return issues, nil
}
