package config

import (
	"errors"
	"fmt"
	"strings"
)

func UpdateToken() error {
	wrapper, err := Load()
	if err != nil {
		return err
	}

	invalid := make(map[string][]errorStruct)
	for i, config := range wrapper.Configs {
		for k, detail := range config.Details {
			// Check token's scope and update the username
			err = detail.Init(&config.URL)
			if err != nil {
				invalid[config.URL.String()] = append(invalid[config.URL.String()], errorStruct{
					tokenName: detail.TokenName,
					err:       err,
				})

				continue
			}

			config.Details[k] = detail
		}

		wrapper.Configs[i] = config
	}

	err = writeConfig(wrapper)
	if err != nil {
		return fmt.Errorf("Failed to write config: %w", err)
	}

	if len(invalid) == 0 {
		// Write was successfull and there were no issues updating the usernames
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

	return errors.New(strings.TrimSuffix(out, "\n"))
}
