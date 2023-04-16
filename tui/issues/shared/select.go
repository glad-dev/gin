package shared

import (
	"errors"
	"net/url"

	"gn/config"
	"gn/config/remote"
	"gn/repo"
	selectconfig "gn/tui/config/select"
)

func SelectConfig(details []repo.Details, u *url.URL) (*config.Wrapper, error) {
	if u != nil {
		return selectConfigForURL(u)
	}

	return selectConfigForLocal(details)
}

func selectConfigForURL(u *url.URL) (*config.Wrapper, error) {
	wrapper, err := config.Load() // To set the colors
	if err != nil {
		return nil, err
	}

	for i, conf := range wrapper.Remotes {
		if u.Host == conf.URL.Host {
			return handleMatch(&wrapper.Remotes[i])
		}
	}

	return nil, errors.New("no matching config found")
}

func selectConfigForLocal(details []repo.Details) (*config.Wrapper, error) {
	wrapper, err := config.Load() // To set the colors
	if err != nil {
		return nil, err
	}

	for _, detail := range details {
		for i, conf := range wrapper.Remotes {
			if conf.URL.Host == detail.URL.Host {
				return handleMatch(&wrapper.Remotes[i])
			}
		}
	}

	return nil, errors.New("no matching config found")
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

	selected, err := selectconfig.Select(conf, "Select the token to use for "+conf.URL.String())
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
