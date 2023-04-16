package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"

	"gn/config/location"
	"gn/config/remote"
	"gn/constants"
	"gn/logger"
	"gn/style"

	"github.com/BurntSushi/toml"
)

// These are needed since the toml libary is unable to decode not-empty interfaces.
type helperWrapper struct {
	Colors  Colors
	Remotes []helperRemote
	Version uint8
}

type helperRemote struct {
	URL     url.URL
	Details []struct {
		Token     string
		TokenName string
		Username  string
	}
}

var ErrConfigDoesNotExist = errors.New("config does not exist")

const ErrConfigDoesNotExistMsg = "No configuration exists.\nRun `" + constants.ProgramName + " config add` to add remotes"

// Load returns the config located at '~/.config/gn/gn.toml' if it exists. If it does not exist, function returns a
// ErrConfigDoesNotExist error and an initialized Wrapper config.
func Load() (*Wrapper, error) {
	fileLocation, err := location.Get()
	if err != nil {
		return nil, err
	}

	// Load config
	helper := &helperWrapper{}
	metaData, err := toml.DecodeFile(fileLocation, helper)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Log.Infof("Found no configuration file at: %s", fileLocation)

			return &Wrapper{
				Remotes: []Remote{},
				Version: constants.ConfigVersion,
			}, ErrConfigDoesNotExist
		}

		logger.Log.Errorf("toml decode failed: %s", err)

		return nil, fmt.Errorf("could not decode config: %w", err)
	}

	// Check if the config only contains the keys we expect
	if len(metaData.Undecoded()) > 0 {
		logger.Log.Error("Config contains unexpected keys.", "invalidKeys", metaData.Undecoded())

		return nil, fmt.Errorf("config contains unexpected keys: %+v", metaData.Undecoded())
	}

	wrap := &Wrapper{
		Colors:  helper.Colors,
		Version: helper.Version,
		Remotes: make([]Remote, len(helper.Remotes)),
	}

	var toAdd Remote
	for i, r := range helper.Remotes {
		toAdd = Remote{
			URL:     r.URL,
			Details: make([]remote.Details, 0),
		}

		for _, details := range r.Details {
			if r.URL.Host == "github.com" {
				toAdd.Details = append(toAdd.Details, remote.GitHubDetails{
					Token:     details.Token,
					TokenName: details.TokenName,
					Username:  details.Username,
				})

				continue
			}

			toAdd.Details = append(toAdd.Details, remote.GitLabDetails{
				Token:     details.Token,
				TokenName: details.TokenName,
				Username:  details.Username,
			})
		}

		wrap.Remotes[i] = toAdd
	}

	err = wrap.CheckValidity()
	if err != nil {
		return nil, fmt.Errorf("config is invalid: %w", err)
	}

	err = wrap.Colors.setColors()
	if err != nil {
		return nil, err
	}

	style.Init()

	return wrap, nil
}
