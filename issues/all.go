package issues

import (
	"net/url"

	"gn/config"
	"gn/logger"
	"gn/repo"
)

func QueryAll(conf *config.Wrapper, details []repo.Details, u *url.URL) ([]Issue, error) {
	match, projectPath, err := getMatchingConfig(conf, details, u)
	if err != nil {
		logger.Log.Errorf("Failed to get matching config: %s", err)

		return nil, err
	}

	if match.URL.Host == "github.com" {
		return queryAllGitHub(match, projectPath)
	}

	return queryAllGitLab(match, projectPath)
}
