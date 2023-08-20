package shared

import (
	"errors"
	"net/url"

	"github.com/glad-dev/gin/config"
	"github.com/glad-dev/gin/logger"
	"github.com/glad-dev/gin/remote"
	"github.com/glad-dev/gin/repo"
	selection "github.com/glad-dev/gin/tui/config/select"
)

// SelectConfig returns the config associated with the URL if a URL is passed. Otherwise, a remote.Details is returned.
func SelectConfig(details []repo.Details, u *url.URL) (*config.Wrapper, error) {
	if u != nil {
		return selectConfigForURL(u)
	}

	return selectConfigForLocal(details)
}

func selectConfigForURL(u *url.URL) (*config.Wrapper, error) {
	wrapper, err := config.Load() // To set the colors
	if err != nil {
		if errors.Is(err, config.ErrConfigDoesNotExist) {
			return &config.Wrapper{
				Remotes: []config.Remote{},
			}, nil
		}

		return nil, err
	}

	for i, conf := range wrapper.Remotes {
		if u.Host == conf.URL.Host {
			return handleMatch(&wrapper.Remotes[i])
		}
	}

	logger.Log.Info("Found no matching config", "URL", u.String())

	return &config.Wrapper{
		Remotes: []config.Remote{},
	}, nil
}

func selectConfigForLocal(details []repo.Details) (*config.Wrapper, error) {
	wrapper, err := config.Load() // To set the colors
	if err != nil {
		if errors.Is(err, config.ErrConfigDoesNotExist) {
			return &config.Wrapper{
				Remotes: []config.Remote{},
			}, nil
		}

		return nil, err
	}

	for _, detail := range details {
		for i, conf := range wrapper.Remotes {
			if conf.URL.Host == detail.URL.Host {
				return handleMatch(&wrapper.Remotes[i])
			}
		}
	}

	logger.Log.Info("Found no matching config", "details", details)

	return &config.Wrapper{
		Remotes: []config.Remote{},
	}, nil
}

func handleMatch(conf *config.Remote) (*config.Wrapper, error) {
	if len(conf.Details) == 1 {
		return &config.Wrapper{
			Remotes: []config.Remote{
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

	return &config.Wrapper{
		Remotes: []config.Remote{
			{
				URL:     conf.URL,
				Details: []remote.Details{*selected},
			},
		},
	}, nil
}
