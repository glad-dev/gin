package gitlab

import (
	"fmt"

	"github.com/glad-dev/gin/logger"
	"github.com/xanzy/go-gitlab"
)

// Init checks the token's scope and returns the username and token name associated with the token.
func Init(token string, apiURL string) (string, string, error) {
	tokenName, err := CheckTokenScope(token, apiURL)
	if err != nil {
		logger.Log.Error("Failed to check scope", "error", err)

		return "", "", fmt.Errorf("GitLabDetails.Init: Failed to check scope: %w", err)
	}

	username, err := getUsername(token, apiURL)
	if err != nil {
		logger.Log.Error("Failed to get username", "error", err)

		return "", "", fmt.Errorf("GitLabDetails.Init: Failed to get Username: %w", err)
	}

	return username, tokenName, nil
}

func getUsername(token string, apiURL string) (string, error) {
	client, err := gitlab.NewClient(token, gitlab.WithBaseURL(apiURL))
	if err != nil {
		logger.Log.Error("Creating gitlab client",
			"error", err,
			"API-URL", apiURL,
		)

		return "", fmt.Errorf("creating gitlab client: %w", err)
	}

	user, _, err := client.Users.CurrentUser()
	if err != nil {
		logger.Log.Error("Requesting current user",
			"error", err,
			"API-URL", apiURL,
		)

		return "", fmt.Errorf("requesting current user: %w", err)
	}

	if len(user.Username) == 0 {
		return "", fmt.Errorf("empty username")
	}

	return user.Username, nil
}
