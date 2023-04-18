package gitlab

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

	err = lab.getUsername(u)
	if err != nil {
		logger.Log.Errorf("Failed to get username: %s", err)

		return nil, fmt.Errorf("GitLabDetails.Init: Failed to get Username: %w", err)
	}

	lab.TokenName = tokenName

	return lab, nil
}

func (lab *Details) getUsername(u *url.URL) error {
	type returnType struct {
		Data struct {
			CurrentUser struct {
				Username string `json:"username"`
			} `json:"currentUser"`
		} `json:"data"`
	}

	query := `
		query {
		  currentUser {
			username
		  }
		}
	`

	response, err := requests.Do(&requests.GraphqlQuery{
		Query:     query,
		Variables: nil,
	}, &remote.Match{
		URL:   *u,
		Token: lab.Token,
	})
	if err != nil {
		return err
	}

	tmp := returnType{}

	dec := json.NewDecoder(bytes.NewBuffer(response))
	dec.DisallowUnknownFields()
	err = dec.Decode(&tmp)
	if err != nil {
		return fmt.Errorf("unmarshle of Username failed: %w", err)
	}

	if len(tmp.Data.CurrentUser.Username) == 0 {
		return errors.New("empty Username: Check API key")
	}

	lab.Username = tmp.Data.CurrentUser.Username

	return nil
}
