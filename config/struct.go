package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"gn/constants"
	"gn/repo"
	"gn/requests"
)

type Wrapper struct {
	Configs      []GitLab
	MajorVersion int
}

type GitLab struct {
	URL      url.URL
	Token    string
	Username string
}

var ErrNoMatchingConfig = errors.New("no matching config was found")

func (config *Wrapper) CheckValidity() error {
	if len(config.Configs) == 0 {
		return errors.New("config file does not contain []GitLab")
	}

	// Check version
	if config.MajorVersion > constants.CurrentMajorVersion {
		return fmt.Errorf("config was written by a newer version of the tool")
	}

	for _, singleConfig := range config.Configs {
		// Check URL
		_, err := checkURLStr(singleConfig.URL.String())
		if err != nil {
			return err
		}

		// Check if token is semantically correct. The tokens validity is not checked
		if len(singleConfig.Token) < 20 {
			return fmt.Errorf("config contains token that is too short. Expected: at least 20, got %d", len(singleConfig.Token))
		}

		if len(singleConfig.Username) == 0 {
			return fmt.Errorf("config contains empty username")
		}
	}

	return nil
}

func (config *Wrapper) GetMatchingConfig(details []repo.Details) (*GitLab, string, error) {
	for _, detail := range details {
		for _, lab := range config.Configs {
			if lab.URL.Host == detail.URL().Host {
				return &lab, detail.ProjectPath(), nil
			}
		}
	}

	return nil, "", ErrNoMatchingConfig
}

func (l *GitLab) GetURL() string {
	return l.URL.String()
}

func (l *GitLab) GetToken() string {
	return l.Token
}

func (l *GitLab) GetUsername() error {
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
	err = json.Unmarshal(response, &tmp)
	if err != nil {
		return fmt.Errorf("unmarshle of username failed: %w", err)
	}

	if len(tmp.Data.CurrentUser.Username) == 0 {
		return errors.New("empty username: Check API key")
	}

	l.Username = tmp.Data.CurrentUser.Username

	return nil
}

func (l *GitLab) CheckTokenValidity() error {
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
	err = json.Unmarshal(response, &scp)
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
