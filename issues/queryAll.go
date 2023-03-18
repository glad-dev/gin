package issues

import (
	"encoding/json"
	"fmt"
	"time"

	"gn/config"
	"gn/requests"
)

type queryAllResponse struct {
	Data struct {
		Project struct {
			Issues struct {
				Nodes []struct {
					Title       string    `json:"title"`
					Description string    `json:"description"`
					CreatedAt   time.Time `json:"createdAt"`
					UpdatedAt   time.Time `json:"updatedAt"`
					Iid         string    `json:"iid"`
					State       string    `json:"state"`
					Author      User      `json:"author"`
					Assignees   struct {
						Nodes []User `json:"nodes"`
					} `json:"assignees"`
				} `json:"nodes"`
			} `json:"issues"`
		} `json:"project"`
	} `json:"data"`
}

func QueryAll(config *config.General, projectPath string, url string) ([]Issue, error) {
	query := `
		query($projectPath: ID!) {
		  project(fullPath: $projectPath) {
			issues(sort: CREATED_ASC) {
			  nodes {
				title
				description
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

	variables := map[string]string{
		"projectPath": projectPath,
	}

	lab, err := config.GetMatchingConfig(url)
	if err != nil {
		return nil, err
	}

	response, err := requests.MakeRequest(&requests.GraphqlQuery{
		Query:     query,
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
	for _, issue := range queryAll.Data.Project.Issues.Nodes {
		issues = append(issues, Issue{
			Title:       issue.Title,
			Description: issue.Description,
			CreatedAt:   issue.CreatedAt,
			UpdatedAt:   issue.UpdatedAt,
			Iid:         issue.Iid,
			State:       issue.State,
			Assignees:   issue.Assignees.Nodes,
			Author:      issue.Author,
		})
	}

	return issues, nil
}
