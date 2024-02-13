package configuration

import (
	"errors"
	"fmt"
	"strings"

	"github.com/glad-dev/gin/log"
	r "github.com/glad-dev/gin/remote"
)

// UpdateRemote loads the configuration file and updates the username and token name associated with each token.
func UpdateRemote() error {
	config, err := Load()
	if err != nil {
		return err
	}

	invalid := make(map[string][]errorStruct)
	var d r.Details
	for i, remote := range config.Remotes {
		for k, detail := range remote.Details {
			// Check token's scope and update the username
			d = detail
			// We can ignore the gosec linter since go 1.22 creates a new variable for each for-loop iteration.
			err = d.Init(&remote.URL) // nolint: gosec
			if err != nil {
				invalid[remote.URL.String()] = append(invalid[remote.URL.String()], errorStruct{
					tokenName: d.TokenName,
					err:       err,
				})

				continue
			}

			remote.Details[k] = d
		}

		config.Remotes[i] = remote
	}

	err = write(config)
	if err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	if len(invalid) == 0 {
		// write was successful and there were no issues updating the usernames
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

	log.Error("Not all remotes could be updated.", "error", out)

	return errors.New(strings.TrimSuffix(out, "\n"))
}
