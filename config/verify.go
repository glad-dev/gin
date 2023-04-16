package config

import (
	"errors"
	"fmt"
	"strings"

	"gn/logger"
)

type errorStruct struct {
	err       error
	tokenName string
}

func VerifyTokens() error {
	wrapper, err := Load()
	if err != nil {
		return err
	}

	invalid := make(map[string][]errorStruct)
	for _, config := range wrapper.Remotes {
		for _, detail := range config.Details {
			err = detail.checkTokenScope(&config.URL)
			if err != nil {
				logger.Log.Error("Failed to check token scope.", "error", err, "URL", config.URL.String())

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

	logger.Log.Error("Not all tokens could be verified.", "error", out)

	return errors.New(strings.TrimSuffix(out, "\n"))
}
