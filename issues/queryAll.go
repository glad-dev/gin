package issues

import (
	"encoding/json"
	"fmt"
	"log"
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
				} `json:"nodes"`
			} `json:"issues"`
		} `json:"project"`
	} `json:"data"`
}

func QueryAll(config *config.Gitlab, projectPath string) {
	// ToDo: Get assignees
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
		log.Fatalf("query all issues: %s\n", err)
	}

	rs := queryAllResponse{}
	err = json.Unmarshal(response, &rs)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Raw:\n%s\n", response)

	nodes := rs.Data.Project.Issues.Nodes
	fmt.Printf("List of all (%d) issues:\n", len(nodes))
	for _, node := range nodes {
		fmt.Printf("%s) %s - %s [%s]\n", node.Iid, node.Title, node.Description, node.State)
	}
}
