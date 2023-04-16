package requests

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func checkErrorGitHub(response []byte) error {
	errorResponse := struct {
		Errors []struct {
			Type      string   `json:"type"`
			Message   string   `json:"message"`
			Path      []string `json:"path"`
			Locations []struct {
				Line   int `json:"line"`
				Column int `json:"column"`
			} `json:"locations"`
		} `json:"errors"`
	}{}

	err := json.Unmarshal(response, &errorResponse)
	if err != nil {
		// If unmarshal fails, then the message is not an error, thus we return nil
		return nil //nolint:nilerr
	}

	if len(errorResponse.Errors) == 0 {
		return nil
	}

	out := ""
	for i, element := range errorResponse.Errors {
		out += fmt.Sprintf(
			"Error %d: %s at line %d, column %d of type %s\n",
			i+1,
			element.Message,
			element.Locations[0].Line,
			element.Locations[0].Column,
			element.Type,
		)
	}

	return fmt.Errorf(strings.TrimSuffix(out, "\n"))
}

func checkExistenceGitHub(response []byte) bool {
	emptyResponse := struct {
		Data struct {
			Repository interface{} `json:"repository"`
		} `json:"data"`
	}{}

	err := json.Unmarshal(response, &emptyResponse)
	if err != nil {
		return false
	}

	return !reflect.ValueOf(emptyResponse.Data.Repository).IsValid()
}
