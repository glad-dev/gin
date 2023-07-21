package gitlab

import (
	"fmt"
	"strings"
	"time"

	"github.com/glad-dev/gin/constants"
	"github.com/glad-dev/gin/logger"
	"github.com/xanzy/go-gitlab"
)

// CheckTokenScope checks the scope of the token and returns the token name.
func CheckTokenScope(token string, apiURL string) (string, error) {
	client, err := gitlab.NewClient(token, gitlab.WithBaseURL(apiURL))
	if err != nil {
		logger.Log.Error("Creating gitlab client",
			"error", err,
			"API-URL", apiURL,
		)

		return "", fmt.Errorf("creating gitlab client: %w", err)
	}

	tokenDetails, _, err := client.PersonalAccessTokens.GetSinglePersonalAccessToken()
	if err != nil {
		logger.Log.Error("Requesting personal access token",
			"error", err,
			"API-URL", apiURL,
		)

		return "", fmt.Errorf("requesting personal access token: %w", err)
	}

	// Check if the token has been revoked
	if tokenDetails.Revoked {
		return "", fmt.Errorf("token was revoked")
	}

	// Check if the token has expired
	if tokenDetails.ExpiresAt != nil {
		date := time.Time(*tokenDetails.ExpiresAt)
		if time.Now().After(date) {
			return "", fmt.Errorf("token expired: %s", date)
		}
	}

	// Make a copy of the required scopes slice
	required := make([]string, len(constants.RequiredGitLabScopes))
	copy(required, constants.RequiredGitLabScopes)

	for _, scope := range tokenDetails.Scopes {
		for i, s := range required {
			if s == scope {
				// Remove the matched scope
				required = append(required[:i], required[i+1:]...)

				break
			}
		}
	}

	if len(required) > 0 {
		return "", fmt.Errorf("some scopes are missing: %s", strings.Join(required, ", "))
	}

	return tokenDetails.Name, nil
}
