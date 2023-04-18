package issues

import (
	"net/url"

	"gn/config"
	"gn/issues/issueList"
	"gn/logger"
	"gn/repo"
)

func QueryAll(conf *config.Wrapper, details []repo.Details, u *url.URL) ([]issueList.Issue, error) {
	match, projectPath, err := getMatchingConfig(conf, details, u)
	if err != nil {
		logger.Log.Errorf("Failed to get matching config: %s", err)

		return nil, err
	}

	if match.URL.Host == "github.com" {
		return issueList.QueryGitHub(match, projectPath)
	}

	return issueList.QueryGitLab(match, projectPath)
}
