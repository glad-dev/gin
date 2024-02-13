package shared

import (
	"errors"
	"net/url"

	"github.com/glad-dev/gin/configuration"
	"github.com/glad-dev/gin/log"
	"github.com/glad-dev/gin/remote"
	"github.com/glad-dev/gin/repository"
	selection "github.com/glad-dev/gin/tui/config/select"
)

// SelectConfig returns the config associated with the URL if a URL is passed. Otherwise, a remote.Details is returned.
func SelectConfig(details []repository.Details, u *url.URL) (*configuration.Config, error) {
	if u != nil {
		return selectConfigForURL(u)
	}

	return selectConfigForLocal(details)
}

func selectConfigForURL(u *url.URL) (*configuration.Config, error) {
	config, err := configuration.Load() // To set the colors
	if err != nil {
		if errors.Is(err, configuration.ErrConfigDoesNotExist) {
			return &configuration.Config{
				Remotes: []configuration.Remote{},
			}, nil
		}

		return nil, err
	}

	for i, conf := range config.Remotes {
		if u.Host == conf.URL.Host {
			return handleMatch(&config.Remotes[i])
		}
	}

	log.Info("Found no matching config", "URL", u.String())

	return &configuration.Config{
		Remotes: []configuration.Remote{},
	}, nil
}

func selectConfigForLocal(details []repository.Details) (*configuration.Config, error) {
	config, err := configuration.Load() // To set the colors
	if err != nil {
		if errors.Is(err, configuration.ErrConfigDoesNotExist) {
			return &configuration.Config{
				Remotes: []configuration.Remote{},
			}, nil
		}

		return nil, err
	}

	for _, detail := range details {
		for i, remoteVar := range config.Remotes {
			if remoteVar.URL.Host == detail.URL.Host {
				return handleMatch(&config.Remotes[i])
			}
		}
	}

	log.Info("Found no matching config", "details", details)

	return &configuration.Config{
		Remotes: []configuration.Remote{},
	}, nil
}

func handleMatch(conf *configuration.Remote) (*configuration.Config, error) {
	if len(conf.Details) == 1 {
		return &configuration.Config{
			Remotes: []configuration.Remote{
				{
					URL:     conf.URL,
					Details: []remote.Details{conf.Details[0]},
				},
			},
		}, nil
	}

	selected, err := selection.Select(conf, "Select the token to use for "+conf.URL.String())
	if err != nil {
		return nil, err
	}

	return &configuration.Config{
		Remotes: []configuration.Remote{
			{
				URL:     conf.URL,
				Details: []remote.Details{*selected},
			},
		},
	}, nil
}
