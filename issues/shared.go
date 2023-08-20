package issues

import (
	"net/url"
	"strings"

	"github.com/glad-dev/gin/config"
	"github.com/glad-dev/gin/remote/match"
	"github.com/glad-dev/gin/repo"
)

func getMatchingConfig(conf *config.Wrapper, details []repo.Details, u *url.URL) (*match.Match, string, error) {
	if u != nil {
		lab, projectPath := getURLConfig(conf, u)

		return lab, projectPath, nil
	}

	return conf.GetMatchingConfig(details)
}

func getURLConfig(conf *config.Wrapper, u *url.URL) (*match.Match, string) {
	// Get project path
	projectPath := u.EscapedPath()
	projectPath = strings.TrimPrefix(projectPath, "/")
	projectPath = strings.TrimSuffix(projectPath, "/")

	// Check if we have a token
	for _, r := range conf.Remotes {
		if r.URL.Hostname() == u.Hostname() {
			match, err := r.ToMatch()
			if err != nil {
				break
			}

			return match, projectPath
		}
	}

	// We found no match => Mock up a config
	return &match.Match{
		URL: url.URL{
			Scheme: u.Scheme,
			Host:   u.Hostname(),
		},
		Token:    "",
		Username: "",
	}, projectPath
}
