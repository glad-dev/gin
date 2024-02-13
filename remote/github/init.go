package github

import (
	"context"
	"fmt"

	"github.com/glad-dev/gin/log"
	"github.com/google/go-github/v59/github"
	"golang.org/x/oauth2"
)

// Init returns the username and token name associated with the token.
func Init(token string) (string, string, error) {
	username, err := getUsername(token)
	if err != nil {
		log.Error("Failed to get username", "error", err)

		return "", "", fmt.Errorf("GitHubDetails.Init: Failed to get Username: %w", err)
	}

	tokenName := "GitHub token for account " + username

	return username, tokenName, nil
}

func getUsername(token string) (string, error) {
	ctx := context.Background()
	tc := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	))

	client := github.NewClient(tc)
	user, _, err := client.Users.Get(context.Background(), "")
	if err != nil {
		return "", fmt.Errorf("getting username: %w", err)
	}

	return user.GetLogin(), nil
}
