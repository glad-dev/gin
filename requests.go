package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type graphqlQuery struct {
	Query     string            `json:"query"`
	Variables map[string]string `json:"variables"`
}

func makeRequest(query *graphqlQuery, config *GitlabConfig) ([]byte, error) {
	requestBody, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", config.Url, bytes.NewBuffer(requestBody))
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

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("request returned invalid status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read HTTP response body: %w", err)
	}

	return body, nil
}
