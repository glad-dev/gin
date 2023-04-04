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
	if strings.HasPrefix(projectPath, "/") {
		projectPath = projectPath[1:]
	}

	if strings.HasSuffix(projectPath, "/") {
		projectPath = projectPath[:len(projectPath)-1]
	}

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
