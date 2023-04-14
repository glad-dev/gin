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
)

var (
	ErrProjectDoesNotExist = errors.New("the project does not exist")
	ErrNotFound            = errors.New("received a 404 - not found when contacting API")
)

func Do(query *GraphqlQuery, config ConfigInterface) (*bytes.Buffer, error) {
	requestBody, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	u, err := url.Parse(config.GetURL())
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %w", err)
	}

	u = u.JoinPath("/api/graphql")
	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.GetToken()))

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
		return nil, ErrNotFound
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("request returned invalid status code %d", resp.StatusCode)
	}

	err = checkForError(bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(body), nil
}

func Project(query *GraphqlQuery, config ConfigInterface) (io.Reader, error) {
	body, err := Do(query, config)
	if err != nil {
		return nil, err
	}

	if projectDoesNotExist(bytes.NewBuffer(body.Bytes())) {
		return nil, ErrProjectDoesNotExist
	}

	return body, nil
}

func checkForError(response io.Reader) error {
	errorResponse := struct {
		Errors []struct {
			Extensions struct {
				Code      string `json:"code"`
				TypeName  string `json:"typeName"`
				FieldName string `json:"fieldName"`
			} `json:"extensions"`
			Message   string `json:"message"`
			Locations []struct {
				Line   int `json:"line"`
				Column int `json:"column"`
			} `json:"locations"`
			Path []string `json:"path"`
		} `json:"errors"`
	}{}

	dec := json.NewDecoder(response)
	dec.DisallowUnknownFields()
	err := dec.Decode(&errorResponse)
	if err != nil {
		// If unmarshal fails, then the message is not an error, thus we return nil
		return nil //nolint:nilerr
	}

	out := ""
	for i, err := range errorResponse.Errors {
		out += fmt.Sprintf(
			"Error %d: %s at line %d, column %d\n",
			i+1,
			err.Message,
			err.Locations[0].Line,
			err.Locations[0].Column,
		)
	}

	return fmt.Errorf(out)
}

func projectDoesNotExist(response io.Reader) bool {
	emptyResponse := struct {
		Data struct {
			Project interface{} `json:"project"`
		} `json:"data"`
	}{}

	dec := json.NewDecoder(response)
	dec.DisallowUnknownFields()
	err := dec.Decode(&emptyResponse)
	if err != nil {
		return false
	}

	return !reflect.ValueOf(emptyResponse.Data.Project).IsValid()
}
