package config

import (
	"errors"
	"fmt"
	"strings"

	"gn/logger"
	"gn/remote"
)

func UpdateRemote() error {
	wrapper, err := Load()
	if err != nil {
		return err
	}

	invalid := make(map[string][]errorStruct)
	var d remote.Details
	for i, config := range wrapper.Remotes {
		for k, detail := range config.Details {
			// Check token's scope and update the username
			d, err = detail.Init(&config.URL)
			if err != nil {
				invalid[config.URL.String()] = append(invalid[config.URL.String()], errorStruct{
					tokenName: detail.GetTokenName(),
					err:       err,
				})

				continue
			}

			config.Details[k] = d
		}

		wrapper.Remotes[i] = config
	}

	err = Write(wrapper)
	if err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	if len(invalid) == 0 {
		// Write was successful and there were no issues updating the usernames
		return nil
	}

	out := "The username for the following tokens could not be updated:\n"
	for urlStr, errorStructs := range invalid {
		for _, errStruct := range errorStructs {
			out += fmt.Sprintf(
				"- Remote '%s' with token '%s': %s\n",
				urlStr,
				errStruct.tokenName,
				errStruct.err.Error(),
			)
		}
	}

	logger.Log.Error("Not all remotes could be updated.", "error", out)

	return errors.New(strings.TrimSuffix(out, "\n"))
}
