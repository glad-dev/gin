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
	"gn/remote"
)

var (
	ErrProjectDoesNotExist = errors.New("the project does not exist")
	ErrNotFound            = errors.New("received a 404 - not found when contacting API")
)

func Project(query interface{}, match *remote.Match) ([]byte, error) {
	body, err := Do(query, match)
	if err != nil {
		return nil, err
	}

	if !projectExists(&match.URL, body) {
		logger.Log.Error("Project does not exist", "body", body)

		return nil, ErrProjectDoesNotExist
	}

	return body, nil
}

func Do(query interface{}, match *remote.Match) ([]byte, error) {
	requestBody, err := json.Marshal(query)
	if err != nil {
		logger.Log.Error("Failed to marshal query", "error", err, "query", query)

		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	return makeRequest(requestBody, match)
}

func makeRequest(requestBody []byte, match *remote.Match) ([]byte, error) {
	req, err := http.NewRequest("POST", match.GetApiURL(), bytes.NewBuffer(requestBody))
	if err != nil {
		logger.Log.Error("Failed to create HTTP request", "error", err, "url", match.URL.String(), "body", string(requestBody))

		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", match.Token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Log.Error("Failed to perform HTTP request", "error", err, "request", req)

		return nil, fmt.Errorf("failed to do HTTP request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Log.Error("Failed to read HTTP response body", "error", err)

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

	err = checkError(&match.URL, body)
	if err != nil {
		logger.Log.Error("Request body contains an error", "error", err, "body", string(body))

		return nil, err
	}

	return body, nil
}

func checkError(u *url.URL, response []byte) error {
	if u.Host == "github.com" {
		return checkErrorGitHub(response)
	}

	return checkErrorGitLab(response)
}

func projectExists(u *url.URL, response []byte) bool {
	if u.Host == "github.com" {
		return projectExistsGitHub(response)
	}

	return projectExistsGitLab(response)
}
