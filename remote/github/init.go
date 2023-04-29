package github

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/google/go-github/v52/github"
	"golang.org/x/oauth2"

	"gn/logger"
	"gn/remote"
)

// Init queries the username associated with the token and updates the token name as well.
func (hub Details) Init(u *url.URL) (remote.Details, error) {
	if u.Host != "github.com" {
		logger.Log.Errorf("Got GitHubDetails with invalid host: %s", u.Host)

		return nil, errors.New("initializing GitHubDetails with a non-GitHub URL")
	}

	err := hub.getUsername()
	if err != nil {
		logger.Log.Error("Failed to get username.", "error", err, "GitHubDetails", hub)

		return nil, fmt.Errorf("GitHubDetails.Init: Failed to get Username: %w", err)
	}

	hub.TokenName = "GitHub token for account " + hub.Username

	return hub, nil
}

func (hub *Details) getUsername() error {
	ctx := context.Background()
	tc := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: hub.Token},
	))

	client := github.NewClient(tc)
	user, _, err := client.Users.Get(context.Background(), "")
	if err != nil {
		return fmt.Errorf("getting username: %w", err)
	}

	hub.Username = user.GetLogin()

	return nil
}
