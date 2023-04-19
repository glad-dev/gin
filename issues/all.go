package issues

import (
	"net/url"

	"gn/config"
	"gn/issues/list"
	"gn/logger"
	"gn/repo"
)

func QueryList(conf *config.Wrapper, details []repo.Details, u *url.URL) ([]list.Issue, error) {
	match, projectPath, err := getMatchingConfig(conf, details, u)
	if err != nil {
		logger.Log.Errorf("Failed to get matching config: %s", err)

		return nil, err
	}

	if match.URL.Host == "github.com" {
		return list.QueryGitHub(match, projectPath)
	}

	return list.QueryGitLab(match, projectPath)
}
