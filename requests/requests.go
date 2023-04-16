package requests

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"gn/logger"
)

var (
	ErrProjectDoesNotExist = errors.New("the project does not exist")
	ErrNotFound            = errors.New("received a 404 - not found when contacting API")
)

func Do(query *GraphqlQuery, config ConfigInterface) (*bytes.Buffer, error) {
	requestBody, err := json.Marshal(query)
	if err != nil {
		logger.Log.Error("Failed to marshal query", "error", err, "query", query)

		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	u, err := url.Parse(config.GetURL())
	if err != nil {
		logger.Log.Error("Failed to parse url", "error", err, "url", config.GetURL())

		return nil, fmt.Errorf("failed to parse url: %w", err)
	}

	req, err := http.NewRequest("POST", getGraphQLURL(u), bytes.NewBuffer(requestBody))
	if err != nil {
		logger.Log.Error("Failed to create HTTP request", "error", err, "url", u.String(), "body", string(requestBody))

		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.GetToken()))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Log.Error("Failed to perform HTTP request", "error", err, "request", req)

		return nil, fmt.Errorf("failed to do HTTP request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Log.Error("Failed to read HTTP respone body", "error", err)

		return nil, fmt.Errorf("failed to read HTTP response body: %w", err)
	}

	switch resp.StatusCode {
	case http.StatusUnauthorized:
		logger.Log.Error("Request has a 401 status code", "body", string(body))

		return nil, errors.New("token is invalid")
	case http.StatusNotFound:
		logger.Log.Error("Request has a 404 status code", "body", string(body))

		return nil, ErrNotFound
	}

	if resp.StatusCode != 200 {
		logger.Log.Error("Request has an unexpected status code", "statusCode", resp.Status, "body", string(body))

		return nil, fmt.Errorf("request returned invalid status code %d", resp.StatusCode)
	}

	err = checkError(u, body)
	if err != nil {
		logger.Log.Error("Request body contains an error", "error", err, "body", string(body))

		return nil, err
	}

	return bytes.NewBuffer(body), nil
}

func Project(query *GraphqlQuery, config ConfigInterface) (io.Reader, error) {
	body, err := Do(query, config)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(config.GetURL())
	if err != nil {
		logger.Log.Error("Failed to parse url", "error", err, "url", config.GetURL())

		return nil, fmt.Errorf("failed to parse url: %w", err)
	}

	if checkExistence(u, body.Bytes()) {
		logger.Log.Error("Project does not exist", "body", body)

		return nil, ErrProjectDoesNotExist
	}

	return body, nil
}

func getGraphQLURL(u *url.URL) string {
	if u.Host == "github.com" {
		return "https://api.github.com/graphql"
	}

	return u.JoinPath("/api/graphql").String()
}

func checkError(u *url.URL, response []byte) error {
	if u.Host == "github.com" {
		return checkErrorGithub(response)
	}

	return checkErrorGitlab(bytes.NewBuffer(response))
}

func checkExistence(u *url.URL, response []byte) bool {
	if u.Host == "github.com" {
		return checkExistenceGithub(response)
	}

	return checkExistenceGitlab(bytes.NewBuffer(response))
}
