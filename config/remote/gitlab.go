package remote

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gn/constants"
	"gn/logger"
	"gn/requests"
)

type GitlabDetails struct {
	Token     string
	TokenName string
	Username  string
}

func (lab GitlabDetails) GetToken() string {
	return lab.Token
}

func (lab GitlabDetails) GetTokenName() string {
	return lab.TokenName
}

func (lab GitlabDetails) GetUsername() string {
	return lab.Username
}

func (lab GitlabDetails) Init(u *url.URL) (Details, error) {
	tokenName, err := lab.CheckTokenScope(u)
	if err != nil {
		logger.Log.Errorf("Failed to check scope: %s", err)

		return nil, fmt.Errorf("GitlabDetails.Init: Failed to check scope: %w", err)
	}

	err = lab.getUsername(u)
	if err != nil {
		logger.Log.Errorf("Failed to get username: %s", err)

		return nil, fmt.Errorf("GitlabDetails.Init: Failed to get Username: %w", err)
	}

	lab.TokenName = tokenName

	return lab, nil
}

func (lab *GitlabDetails) getUsername(u *url.URL) error {
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
	}, &Match{
		URL:   *u,
		Token: lab.Token,
	})
	if err != nil {
		return err
	}

	tmp := returnType{}

	dec := json.NewDecoder(response)
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

func (lab GitlabDetails) CheckTokenScope(u *url.URL) (string, error) {
	response := struct {
		CreatedAt  time.Time   `json:"created_at"`
		LastUsedAt time.Time   `json:"last_used_at"`
		ExpiresAt  interface{} `json:"expires_at"`
		Name       string      `json:"name"`
		Scopes     []string    `json:"scopes"`
		ID         int         `json:"id"`
		UserID     int         `json:"user_id"`
		Revoked    bool        `json:"revoked"`
		Active     bool        `json:"active"`
	}{}

	req, err := http.NewRequest("GET", u.JoinPath("/api/v4/personal_access_tokens/self").String(), bytes.NewBufferString(""))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("PRIVATE-TOKEN", lab.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read body: %w", err)
	}

	if resp.StatusCode == 401 {
		return "", errors.New("the provided URL or token are invalid")
	}

	if resp.StatusCode < 200 || resp.StatusCode > 399 {
		return "", fmt.Errorf("the status code of %d indicates failure: %s", resp.StatusCode, body)
	}

	dec := json.NewDecoder(bytes.NewBuffer(body))
	dec.DisallowUnknownFields()
	err = dec.Decode(&response)
	if err != nil {
		return "", fmt.Errorf("failed to decode response: %w - %s", err, body)
	}

	if response.Revoked {
		return "", fmt.Errorf("token was revoked")
	}

	if response.ExpiresAt != nil {
		date, ok := response.ExpiresAt.(time.Time)
		if ok && time.Now().After(date) {
			return "", fmt.Errorf("token expired: %s", response.ExpiresAt)
		}
	}

	// Make a copy of the required scopes slice
	required := make([]string, len(constants.RequiredGitlabScopes))
	copy(required, constants.RequiredGitlabScopes)

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
		return "", fmt.Errorf("some scopes are missing: %s", strings.Join(required, ", "))
	}

	return response.Name, nil
}
