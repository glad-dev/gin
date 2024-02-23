package configuration

import (
	"errors"
	"fmt"
	"strings"

	"github.com/glad-dev/gin/log"
)

type errorStruct struct {
	err       error
	tokenName string
}

// VerifyTokens loads the configuration file and checks if every token is valid.
func VerifyTokens() error {
	config, err := Load()
	if err != nil {
		return err
	}

	invalid := make(map[string][]errorStruct)
	for _, remote := range config.Remotes {
		for _, detail := range remote.Details {
			// We can ignore the gosec linter since go 1.22 creates a new variable for each for-loop iteration.
			_, err = detail.CheckTokenScope(&remote.URL, remote.Type) // nolint: gosec
			if err != nil {
				log.Error("Failed to check token scope.", "error", err, "URL", remote.URL.String())

				invalid[remote.URL.String()] = append(invalid[remote.URL.String()], errorStruct{
					tokenName: detail.TokenName,
					err:       err,
				})
			}
		}
	}

	if len(invalid) == 0 {
		return nil
	}

	out := "The following remotes/tokens could not be verified:\n"
	for urlStr, errorStructs := range invalid {
		for _, errStruct := range errorStructs {
			out += fmt.Sprintf(
				"- Remote '%s' contains invalid token named '%s': %s\n",
				urlStr,
				errStruct.tokenName,
				errStruct.err.Error(),
			)
		}
	}

	log.Error("Not all tokens could be verified.", "error", out)

	return errors.New(strings.TrimSuffix(out, "\n"))
}
