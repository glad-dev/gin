package issues

import (
	"fmt"
	"log"

	"gn/config"
	"gn/requests"
)

func QuerySingle(config *config.Gitlab, projectPath string, issueID string) {
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
		"issueIID":    issueID,
	}

	response, err := requests.MakeRequest(&requests.GraphqlQuery{
		Query:     query,
		Variables: variables,
	}, config)

	if err != nil {
		log.Fatalf("query single: %s\n", err)
	}

	fmt.Printf("Single issue:\n%s\n", response)
}
