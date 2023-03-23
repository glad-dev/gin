package requests

import (
	"encoding/json"
	"fmt"
	"strings"

	"gn/config"
	"gn/constants"
)

func CheckTokenValidity(lab *config.GitLab) error {
	query := `
		query {
		  viewer {
			scopes
		  }
		}
	`

	response, err := Do(&GraphqlQuery{
		Query:     query,
		Variables: map[string]string{},
	}, lab)

	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	scp := scopes{}
	err = json.Unmarshal(response, &scp)
	if err != nil {
		return fmt.Errorf("unmarshal failed: %w", err)
	}

	// Make a copy of the required scopes slice
	required := make([]string, len(constants.RequiredScopes))
	copy(required, constants.RequiredScopes)

	for _, scope := range scp.Data.Viewer.Scopes {
		for i, s := range required {
			if s == scope {
				// Remove the matched scope
				required = append(required[:i], required[i+1:]...)

				break
			}
		}
	}

	if len(required) > 0 {
		return fmt.Errorf("some scopes are missing: %s", strings.Join(required, ", "))
	}

	return nil
}
