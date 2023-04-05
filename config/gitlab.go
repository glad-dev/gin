package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

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
	response := struct {
		Id         int         `json:"id"`
		Name       string      `json:"name"`
		Revoked    bool        `json:"revoked"`
		CreatedAt  time.Time   `json:"created_at"`
		Scopes     []string    `json:"scopes"`
		UserId     int         `json:"user_id"`
		LastUsedAt time.Time   `json:"last_used_at"`
		Active     bool        `json:"active"`
		ExpiresAt  interface{} `json:"expires_at"`
	}{}

	req, err := http.NewRequest("GET", "https://gitlab.com/api/v4/personal_access_tokens/self", bytes.NewBufferString(""))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("PRIVATE-TOKEN", l.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %s", err)
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	dec.DisallowUnknownFields()
	err = dec.Decode(&response)
	if err != nil {
		return fmt.Errorf("failed to decode response: %s", err)
	}

	if response.Revoked {
		return fmt.Errorf("token was revoked")
	}

	if response.ExpiresAt != nil {
		date, ok := response.ExpiresAt.(time.Time)
		if ok && time.Now().After(date) {
			return fmt.Errorf("token expired: %s", response.ExpiresAt)
		}
	}

	// Make a copy of the required scopes slice
	required := make([]string, len(constants.RequiredScopes))
	copy(required, constants.RequiredScopes)

	for _, scope := range response.Scopes {
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
