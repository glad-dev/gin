package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"gn/constants"
	"gn/requests"
)

type GitLab struct { // TODO: Add support for Github, maybe bool field?
	URL      url.URL
	Token    string
	Username string
}

func (l *GitLab) Init() error {
	err := l.CheckTokenScope()
	if err != nil {
		return err
	}

	return l.GetUsername()
}

func (l *GitLab) GetURL() string {
	return l.URL.String()
}

func (l *GitLab) GetToken() string {
	return l.Token
}

func (l *GitLab) CheckValidity() error {
	// Check URL
	_, err := checkURLStr(l.URL.String())
	if err != nil {
		return err
	}

	// Check if token is semantically correct. The tokens validity is not checked
	if len(l.Token) < 20 { // TODO: Get actual sizes
		return fmt.Errorf("config contains token that is too short. Expected: at least 20, got %d", len(l.Token))
	}

	if len(l.Username) == 0 {
		return fmt.Errorf("config contains empty username")
	}

	return nil
}

var debug = true

func (l *GitLab) GetUsername() error {
	if debug {
		l.Username = "Fake username"

		return nil
	}

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
	}, l)
	if err != nil {
		return err
	}

	tmp := returnType{}

	dec := json.NewDecoder(response)
	dec.DisallowUnknownFields()
	err = dec.Decode(&tmp)
	if err != nil {
		return fmt.Errorf("unmarshle of username failed: %w", err)
	}

	if len(tmp.Data.CurrentUser.Username) == 0 {
		return errors.New("empty username: Check API key")
	}

	l.Username = tmp.Data.CurrentUser.Username

	return nil
}

func (l *GitLab) CheckTokenScope() error {
	type scopes struct {
		Data struct {
			Viewer struct {
				Scopes []string `json:"scopes"`
			} `json:"viewer"`
		} `json:"data"`
	}

	query := `
		query {
		  viewer {
			scopes
		  }
		}
	`

	response, err := requests.Do(&requests.GraphqlQuery{
		Query:     query,
		Variables: map[string]string{},
	}, l)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	scp := scopes{}

	dec := json.NewDecoder(response)
	dec.DisallowUnknownFields()
	err = dec.Decode(&scp)
	if err != nil {
		return fmt.Errorf("unmarshal failed: %w", err)
	}

	// Make a copy of the required scopes slice
	required := make([]string, len(constants.RequiredScopes))
	copy(required, constants.RequiredScopes)

	for _, scope := range scp.Data.Viewer.Scopes {
		for i, s := range required {
			if s == scope {
				// Remove the matched scope
				required = append(required[:i], required[i+1:]...)

				break
			}
		}
	}

	if len(required) > 0 {
		return fmt.Errorf("some scopes are missing: %s", strings.Join(required, ", "))
	}

	return nil
}
