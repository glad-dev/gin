package requests

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"

	"gn/config"
)

var ErrProjectDoesNotExist = errors.New("the project does not exist")

func MakeRequest(query *GraphqlQuery, config *config.GitLab) ([]byte, error) {
	requestBody, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	u, err := url.Parse(config.URL.String())
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %w", err)
	}

	u = u.JoinPath("/api/graphql")
	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.Token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do HTTP request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read HTTP response body: %w", err)
	}

	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return nil, errors.New("token is invalid")
	case http.StatusNotFound:
		return nil, errors.New("either project does not exist or token is invalid")
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("request returned invalid status code %d with message: %s", resp.StatusCode, body)
	}

	if projectDoesNotExist(body) {
		return nil, ErrProjectDoesNotExist
	}

	return body, nil
}

func projectDoesNotExist(response []byte) bool {
	emptyResponse := struct {
		Data struct {
			Project interface{} `json:"project"`
		} `json:"data"`
	}{}

	err := json.Unmarshal(response, &emptyResponse)
	if err != nil {
		return false
	}

	return !reflect.ValueOf(emptyResponse.Data.Project).IsValid()
}
