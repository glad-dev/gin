package remote

import (
	"errors"
	"net/url"

	"github.com/glad-dev/gin/remote/github"
	"github.com/glad-dev/gin/remote/gitlab"
	rt "github.com/glad-dev/gin/remote/type"
)

type Details struct {
	Token     string
	TokenName string
	Username  string
}

// Init checks the token's scope and sets the username and token name associated with the token.
func (d *Details) Init(u *url.URL, remoteType rt.Type) error {
	switch remoteType {
	case rt.Github:
		username, tokenName, err := github.Init(d.Token)
		if err != nil {
			return err
		}

		d.Username = username
		d.TokenName = tokenName

		return nil

	case rt.Gitlab:
		username, tokenName, err := gitlab.Init(d.Token, remoteType, u)
		if err != nil {
			return err
		}

		d.Username = username
		d.TokenName = tokenName

		return nil
	}

	return errors.New("invalid type passed to Details.Init")
}

func (d *Details) CheckTokenScope(u *url.URL, remoteType rt.Type) (string, error) {
	apiURL, err := remoteType.GraphqlAPIURL(u)
	if err != nil {
		return "", err
	}

	switch remoteType {
	case rt.Github:
		return "", nil // Not implemented

	case rt.Gitlab:
		return gitlab.CheckTokenScope(d.Token, apiURL)
	}

	return "", errors.New("Details.CheckTokenScope - invalid type")
}
