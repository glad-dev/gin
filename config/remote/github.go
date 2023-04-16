package remote

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"gn/logger"
	"gn/requests"
)

type GithubDetails struct {
	Token     string
	TokenName string
	Username  string
}

func (hub GithubDetails) GetToken() string {
	return hub.Token
}

func (hub GithubDetails) GetTokenName() string {
	return hub.TokenName
}

func (hub GithubDetails) GetUsername() string {
	return hub.Username
}

func (hub GithubDetails) Init(u *url.URL) (Details, error) {
	err := hub.getUsername(u)
	if err != nil {
		logger.Log.Error("Failed to get Username.", "error", err, "GithubDetails", hub)

		return nil, fmt.Errorf("GithubDetails.Init: Failed to get Username: %w", err)
	}

	tokenName, err := hub.CheckTokenScope(u)
	if err != nil {
		logger.Log.Error("Failed to check scope.", "error", err, "GithubDetails", hub)

		return nil, fmt.Errorf("GithubDetails.Init: Failed to check scope: %w", err)
	}

	hub.TokenName = tokenName

	return hub, nil
}

func (hub *GithubDetails) getUsername(u *url.URL) error {
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

	response, err := requests.Do(&requests.GraphqlQuery{
		Query:     query,
		Variables: nil,
	}, &Match{
		URL:   *u,
		Token: hub.Token,
	})
	if err != nil {
		return err
	}

	dec := json.NewDecoder(response)
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

func (hub GithubDetails) CheckTokenScope(_ *url.URL) (string, error) { // TODO
	return "Github token for account " + hub.Username, nil
}
