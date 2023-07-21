package remote

import (
	"errors"
	"net/url"

	"github.com/glad-dev/gin/remote/github"
	"github.com/glad-dev/gin/remote/gitlab"
)

type Details struct {
	Token     string
	TokenName string
	Username  string
	Type      Type
}

// Init checks the token's scope and sets the username and token name associated with the token.
func (d *Details) Init(u *url.URL) error {
	apiURL, err := d.Type.ApiURL(u)
	if err != nil {
		return err
	}

	switch d.Type {
	case Github:
		username, tokenName, err := github.Init(d.Token)
		if err != nil {
			return err
		}

		d.Username = username
		d.TokenName = tokenName

		return nil

	case Gitlab:
		username, tokenName, err := gitlab.Init(d.Token, apiURL)
		if err != nil {
			return err
		}

		d.Username = username
		d.TokenName = tokenName

		return nil

	case Bitbucket:
		return errors.New("bitbucket is yet to be implemented")
	}

	return errors.New("invalid type passed to Details.Init")
}

func (d *Details) CheckTokenScope(u *url.URL) (string, error) {
	apiURL, err := d.Type.ApiURL(u)
	if err != nil {
		return "", err
	}

	switch d.Type {
	case Github:
		return "", nil // Not implemented

	case Gitlab:
		return gitlab.CheckTokenScope(d.Token, apiURL)

	case Bitbucket:
		return "", errors.New("bitbucket is not yet implemented")
	}

	return "", errors.New("Details.CheckTokenScope - invalid type")
}
