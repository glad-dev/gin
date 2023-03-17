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

func QueryAll(config *config.Gitlab, projectPath string) ([]Issue, error) {
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
	response, err := requests.MakeRequest(&requests.GraphqlQuery{
		Query:     query,
		Variables: variables,
	}, config)

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
		// Iterate over the issue's assignees
		assignees := make([]string, 0)
		for _, assignee := range issue.Assignees.Nodes {
			assignees = append(assignees, assignee.Name)
		}

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

		fmt.Printf("%d)\tAuthor: %#v\n\tAssignees: %#v\n", len(issues), issue.Assignees.Nodes, issue.Author)
	}

	return issues, nil
}
