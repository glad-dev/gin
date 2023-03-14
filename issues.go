package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"
)

type QueryAll struct {
	Data struct {
		Project struct {
			Issues struct {
				Nodes []Issue `json:"nodes"`
			} `json:"issues"`
		} `json:"project"`
	} `json:"data"`
}

type Issue struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	WebUrl      string    `json:"webUrl"`
}

func queryAllIssues(config *GitlabConfig, projectPath string) {
	query := `
		query($projectPath: ID!) {
		  project(fullPath: $projectPath) {
			issues {
			  nodes {
				title
				description
				createdAt
				updatedAt
				webUrl
			  }
			}
		  }
		}
	`

	variables := map[string]string{
		"projectPath": projectPath,
	}
	response, err := makeRequest(&graphqlQuery{
		Query:     query,
		Variables: variables,
	}, config)

	if err != nil {
		log.Fatalf("query all issues: %s\n", err)
	}

	rs := QueryAll{}
	err = json.Unmarshal(response, &rs)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Raw:\n%s\n", response)

	nodes := rs.Data.Project.Issues.Nodes
	fmt.Printf("List of all (%d) issues:\n", len(nodes))
	for i, node := range nodes {
		fmt.Printf("%d) %s - %s\n", i+1, node.Title, node.Description)
	}
}

func querySingleIssue(config *GitlabConfig, projectPath string, issueID int) {
	query := `
		query($projectPath: ID!, $issueIID: String!) {
		  project(fullPath: $projectPath) {
		    issue(iid: $issueIID) {
		      title
		      description
		      createdAt
		      updatedAt
		      webUrl
		      discussions {
		        nodes {
		          id
		          notes {
		            nodes {
		              body
		              createdAt
		              updatedAt
		            }
		          }
		        }
		      }
		    }
		  }
		}
	`

	variables := map[string]string{
		"projectPath": projectPath,
		"issueIID":    strconv.Itoa(issueID),
	}

	response, err := makeRequest(&graphqlQuery{
		Query:     query,
		Variables: variables,
	}, config)

	if err != nil {
		log.Fatalf("query single: %s\n", err)
	}

	fmt.Printf("Single issue:\n%s\n", response)
}
