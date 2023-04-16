package requests

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
)

func checkErrorGitlab(response io.Reader) error {
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

func checkExistenceGitlab(response io.Reader) bool {
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
