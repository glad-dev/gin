package config

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
	wrapper, err := Load()
	if err != nil {
		return err
	}

	invalid := make(map[string][]errorStruct)
	for _, config := range wrapper.Remotes {
		for _, detail := range config.Details {
			_, err = detail.CheckTokenScope(&config.URL)
			if err != nil {
				log.Error("Failed to check token scope.", "error", err, "URL", config.URL.String())

				invalid[config.URL.String()] = append(invalid[config.URL.String()], errorStruct{
					tokenName: detail.TokenName,
					err:       err,
				})
			}
		}
	}

	if len(invalid) == 0 {
		return nil
	}

	out := "The following configs have issues:\n"
	for urlStr, errorStructs := range invalid {
		for _, errStruct := range errorStructs {
			out += fmt.Sprintf(
				"- Remote '%s' contains token '%s' with error: %s\n",
				urlStr,
				errStruct.tokenName,
				errStruct.err.Error(),
			)
		}
	}

	log.Error("Not all tokens could be verified.", "error", out)

	return errors.New(strings.TrimSuffix(out, "\n"))
}
