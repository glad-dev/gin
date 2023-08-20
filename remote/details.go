package remote

import (
	"errors"
	"net/url"

	"github.com/glad-dev/gin/remote/github"
	"github.com/glad-dev/gin/remote/gitlab"
	remotetype "github.com/glad-dev/gin/remote/type"
)

type Details struct {
	Token     string
	TokenName string
	Username  string
	Type      remotetype.Type
}

// Init checks the token's scope and sets the username and token name associated with the token.
func (d *Details) Init(u *url.URL) error {
	switch d.Type {
	case remotetype.Github:
		username, tokenName, err := github.Init(d.Token)
		if err != nil {
			return err
		}

		d.Username = username
		d.TokenName = tokenName

		return nil

	case remotetype.Gitlab:
		username, tokenName, err := gitlab.Init(d.Token, d.Type, u)
		if err != nil {
			return err
		}

		d.Username = username
		d.TokenName = tokenName

		return nil

	case remotetype.Bitbucket:
		return errors.New("bitbucket is yet to be implemented")
	}

	return errors.New("invalid type passed to Details.Init")
}

func (d *Details) CheckTokenScope(u *url.URL) (string, error) {
	apiURL, err := d.Type.GraphqlAPIURL(u)
	if err != nil {
		return "", err
	}

	switch d.Type {
	case remotetype.Github:
		return "", nil // Not implemented

	case remotetype.Gitlab:
		return gitlab.CheckTokenScope(d.Token, apiURL)

	case remotetype.Bitbucket:
		return "", errors.New("bitbucket is not yet implemented")
	}

	return "", errors.New("Details.CheckTokenScope - invalid type")
}
