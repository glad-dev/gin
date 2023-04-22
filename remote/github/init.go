package github

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"gn/logger"
	"gn/remote"
	"gn/requests"
)

// Init queries the username associated with the token and updates the token name as well.
func (hub Details) Init(u *url.URL) (remote.Details, error) {
	if u.Host != "github.com" {
		logger.Log.Errorf("Got GitHubDetails with invalid host: %s", u.Host)

		return nil, errors.New("initializing GitHubDetails with a non-GitHub URL")
	}

	err := hub.getUsername(u)
	if err != nil {
		logger.Log.Error("Failed to get Username.", "error", err, "GitHubDetails", hub)

		return nil, fmt.Errorf("GitHubDetails.Init: Failed to get Username: %w", err)
	}

	tokenName, err := hub.CheckTokenScope(u)
	if err != nil {
		logger.Log.Error("Failed to check scope.", "error", err, "GitHubDetails", hub)

		return nil, fmt.Errorf("GitHubDetails.Init: Failed to check scope: %w", err)
	}

	hub.TokenName = tokenName

	return hub, nil
}

func (hub *Details) getUsername(u *url.URL) error {
	responseType := struct {
		Data struct {
			Viewer struct {
				Login string `json:"login"`
			} `json:"viewer"`
		} `json:"data"`
	}{}

	query := `
		query {
			viewer {
				login
			}
		}
	`

	response, err := requests.Do(&requests.Query{
		Query:     query,
		Variables: nil,
	}, &remote.Match{
		URL:   *u,
		Token: hub.Token,
	})
	if err != nil {
		return err
	}

	dec := json.NewDecoder(bytes.NewBuffer(response))
	dec.DisallowUnknownFields()
	err = dec.Decode(&responseType)
	if err != nil {
		return fmt.Errorf("unmarshle of Username failed: %w", err)
	}

	if len(responseType.Data.Viewer.Login) == 0 {
		return errors.New("empty Username: Check API key")
	}

	hub.Username = responseType.Data.Viewer.Login

	return nil
}
