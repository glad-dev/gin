package issues

import (
	"net/url"
	"strings"

	"gn/config"
	"gn/repo"
)

func getMatchingConfig(conf *config.Wrapper, details []repo.Details, u *url.URL) (*config.Match, string, error) {
	if u != nil {
		lab, projectPath := getURLConfig(conf, u)

		return lab, projectPath, nil
	}

	return conf.GetMatchingConfig(details)
}

func getURLConfig(conf *config.Wrapper, u *url.URL) (*config.Match, string) {
	// Get project path
	projectPath := u.EscapedPath()
	projectPath = strings.TrimPrefix(projectPath, "/")
	projectPath = strings.TrimSuffix(projectPath, "/")

	// Check if we have a token
	for _, remote := range conf.Remotes {
		if remote.URL.Hostname() == u.Hostname() {
			match, err := remote.ToMatch()
			if err != nil {
				break
			}

			return match, projectPath
		}
	}

	// We found no match => Mock up a config
	return &config.Match{
		URL: url.URL{
			Scheme: u.Scheme,
			Host:   u.Hostname(),
		},
		Token:    "",
		Username: "",
	}, projectPath
}
