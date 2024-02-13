package gitlab

import (
	"fmt"
	"net/url"

	"github.com/glad-dev/gin/log"
	remotetype "github.com/glad-dev/gin/remote/type"
)

// Init checks the token's scope and returns the username and token name associated with the token.
func Init(token string, t remotetype.Type, u *url.URL) (string, string, error) {
	restAPIUrl, err := t.RestAPIURL(u)
	if err != nil {
		log.Error("Failed to get REST API URL", "error", err, "URL", u.String())

		return "", "", fmt.Errorf("getting REST API: %w", err)
	}

	tokenName, err := CheckTokenScope(token, restAPIUrl)
	if err != nil {
		log.Error("Failed to check scope", "error", err)

		return "", "", fmt.Errorf("GitLabDetails.Init: Failed to check scope: %w", err)
	}

	username, err := getUsername(token, u)
	if err != nil {
		log.Error("Failed to get username", "error", err)

		return "", "", fmt.Errorf("GitLabDetails.Init: Failed to get Username: %w", err)
	}

	return username, tokenName, nil
}
