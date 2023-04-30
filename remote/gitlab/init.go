package gitlab

import (
	"errors"
	"fmt"
	"net/url"

	"gn/logger"
	"gn/remote"

	"github.com/xanzy/go-gitlab"
)

// Init checks the token's scope and set's the username and token name associated with the token.
func (lab Details) Init(u *url.URL) (remote.Details, error) {
	if u.Host == "github.com" {
		logger.Log.Errorf("Got GitLabDetails with invalid host: %s", u.Host)

		return nil, errors.New("initializing GitLabDetails with a GitHub URL")
	}

	tokenName, err := lab.CheckTokenScope(u)
	if err != nil {
		logger.Log.Errorf("Failed to check scope: %s", err)

		return nil, fmt.Errorf("GitLabDetails.Init: Failed to check scope: %w", err)
	}

	err = lab.setUsername(u)
	if err != nil {
		logger.Log.Errorf("Failed to get username: %s", err)

		return nil, fmt.Errorf("GitLabDetails.Init: Failed to get Username: %w", err)
	}

	lab.TokenName = tokenName

	return lab, nil
}

func (lab *Details) setUsername(u *url.URL) error {
	api := ApiURL(u)

	client, err := gitlab.NewClient(lab.Token, gitlab.WithBaseURL(api))
	if err != nil {
		logger.Log.Error("Creating gitlab client",
			"error", err,
			"API-URL", api,
		)

		return fmt.Errorf("creating gitlab client: %w", err)
	}

	client.PersonalAccessTokens.GetSinglePersonalAccessToken()

	user, _, err := client.Users.CurrentUser()
	if err != nil {
		logger.Log.Error("Requesting current user",
			"error", err,
			"API-URL", api,
		)

		return fmt.Errorf("requesting current user: %w", err)
	}

	if len(user.Username) == 0 {
		return fmt.Errorf("empty username")
	}

	lab.Username = user.Username

	return nil
}
