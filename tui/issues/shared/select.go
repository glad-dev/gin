package shared

import (
	"errors"

	"gn/config"
	"gn/config/remote"
	"gn/repo"
	selectconfig "gn/tui/config/select"
)

func SelectConfig(details []repo.Details) (*config.Wrapper, error) {
	wrapper, err := config.Load() // To set the colors
	if err != nil {
		return nil, err
	}

	for _, detail := range details {
		for i, conf := range wrapper.Remotes {
			if conf.URL.Host == detail.URL.Host {
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

				selected, err := selectconfig.Select(&wrapper.Remotes[i], "Select the token to use for "+conf.URL.String())
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
		}
	}

	return nil, errors.New("no matching config found")
}
