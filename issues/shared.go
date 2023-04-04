package issues

import (
	"net/url"
	"strings"

	"gn/config"
	"gn/repo"
)

func getMatchingConfig(conf *config.Wrapper, details []repo.Details, u *url.URL) (*config.GitLab, string, error) {
	if u != nil {
		lab, projectPath := getURLConfig(conf, u)

		return lab, projectPath, nil
	}

	return conf.GetMatchingConfig(details)
}

func getURLConfig(conf *config.Wrapper, u *url.URL) (*config.GitLab, string) {
	// Get project path
	projectPath := u.EscapedPath()
	projectPath = strings.TrimPrefix(projectPath, "/")
	projectPath = strings.TrimSuffix(projectPath, "/")

	// Check if we have a token
	for _, gitLab := range conf.Configs {
		if gitLab.URL.Hostname() == u.Hostname() {
			return &gitLab, projectPath
		}
	}

	// We found no match => Mock up a config
	return &config.GitLab{
		URL: url.URL{
			Scheme: u.Scheme,
			Host:   u.Hostname(),
		},
		Token:    "",
		Username: "",
	}, projectPath
}
