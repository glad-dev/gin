package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type graphqlQuery struct {
	Query     string            `json:"query"`
	Variables map[string]string `json:"variables"`
}

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
	response := makeRequest(graphqlQuery{
		Query:     query,
		Variables: variables,
	}, config)

	rs := QueryAll{}
	err := json.Unmarshal([]byte(response), &rs)
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

	response := makeRequest(graphqlQuery{
		Query:     query,
		Variables: variables,
	}, config)

	fmt.Printf("Single issue:\n%s\n", response)
}

func makeRequest(requestInterface interface{}, config *GitlabConfig) string {
	requestBody, err := json.Marshal(requestInterface)
	if err != nil {
		log.Fatal("Error marshaling JSON:", err)
	}

	req, err := http.NewRequest("POST", config.Url, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatal("Error creating HTTP request:", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.Token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending HTTP request:", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading HTTP response body:", err)
	}

	return string(body)
}
