package issues

import (
	"net/url"
	"strings"

	"github.com/glad-dev/gin/configuration"
	"github.com/glad-dev/gin/remote/match"
	rt "github.com/glad-dev/gin/remote/type"
	"github.com/glad-dev/gin/repository"
)

func getMatchingConfig(conf *configuration.Config, details []repository.Details, u *url.URL) (*match.Match, string, error) {
	if u != nil {
		lab, projectPath := getURLConfig(conf, u)

		return lab, projectPath, nil
	}

	return conf.GetMatchingConfig(details)
}

func getURLConfig(conf *configuration.Config, u *url.URL) (*match.Match, string) {
	// Get project path
	projectPath := u.EscapedPath()
	projectPath = strings.TrimPrefix(projectPath, "/")
	projectPath = strings.TrimSuffix(projectPath, "/")

	// Check if we have a token
	for _, r := range conf.Remotes {
		if r.URL.Hostname() == u.Hostname() {
			m, err := r.ToMatch()
			if err != nil {
				break
			}

			return m, projectPath
		}
	}

	// We found no match => Mock up a config
	t := rt.Gitlab
	if u.Host == "github.com" {
		t = rt.Github
	}

	return &match.Match{
		URL: url.URL{
			Scheme: u.Scheme,
			Host:   u.Hostname(),
		},
		Token:    "",
		Type:     t,
		Username: "",
	}, projectPath
}
