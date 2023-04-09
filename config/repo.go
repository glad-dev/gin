package config

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
	"gn/requests"
)

type Remote struct {
	URL     url.URL
	Details []RemoteDetails
}

type RemoteDetails struct {
	Token     string
	TokenName string
	Username  string
}

var ErrMultipleRepoDetails = errors.New("config contains multiple matching configs")

func (r *Remote) ToMatch() (*Match, error) {
	if len(r.Details) == 0 {
		return nil, errors.New("failed to convert Remote to Match since Remote.Details contains no elements")
	}

	if len(r.Details) == 1 {
		return &Match{
			URL:       r.URL,
			Token:     r.Details[0].Token,
			Username:  r.Details[0].Username,
			TokenName: r.Details[0].TokenName,
		}, nil
	}

	return &Match{
		URL: r.URL,
	}, ErrMultipleRepoDetails
}

func (r *Remote) ToMatchAtIndex(index int) (*Match, error) {
	if index < 0 || index >= len(r.Details) {
		return nil, errors.New("ToMatchAtIndex: invalid index")
	}

	return &Match{
		URL:       r.URL,
		Token:     r.Details[index].Token,
		Username:  r.Details[index].Username,
		TokenName: r.Details[index].TokenName,
	}, nil
}

func (r *Remote) CheckSemantics() error {
	// Check URL
	_, err := checkURLStr(r.URL.String())
	if err != nil {
		return err
	}

	if len(r.Details) == 0 {
		return fmt.Errorf("config contains empty no details")
	}

	for _, details := range r.Details {
		if len(details.Username) == 0 {
			return fmt.Errorf("config contains empty username")
		}

		// Check if token is semantically correct. The tokens validity is not checked.
		if len(details.Token) == 0 { // TODO: Get actual sizes
			return fmt.Errorf("config contains empty token")
		}

		if len(details.TokenName) == 0 {
			return fmt.Errorf("config contains empty token name")
		}
	}

	return nil
}

func (rd *RemoteDetails) Init(u *url.URL) error {
	err := rd.CheckTokenScope(u)
	if err != nil {
		return fmt.Errorf("RemoteDetails.Init: Failed to check scope: %w", err)
	}

	err = rd.GetUsername(u)
	if err != nil {
		return fmt.Errorf("RemoteDetails.Init: Failed to get username: %w", err)
	}

	return nil
}

func (rd *RemoteDetails) GetUsername(u *url.URL) error {
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
		Token: rd.Token,
	})
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

	rd.Username = tmp.Data.CurrentUser.Username

	return nil
}

func (rd *RemoteDetails) CheckTokenScope(u *url.URL) error {
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
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("PRIVATE-TOKEN", rd.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %w", err)
	}

	if resp.StatusCode == 401 {
		return errors.New("the provided URL or token are invalid")
	}

	if resp.StatusCode < 200 || resp.StatusCode > 399 {
		return fmt.Errorf("the status code of %d indicates failure: %s", resp.StatusCode, body)
	}

	dec := json.NewDecoder(bytes.NewBuffer(body))
	dec.DisallowUnknownFields()
	err = dec.Decode(&response)
	if err != nil {
		return fmt.Errorf("failed to decode response: %w - %s", err, body)
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

	rd.TokenName = response.Name

	return nil
}
