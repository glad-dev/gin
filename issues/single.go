package issues

import (
	"net/url"

	"gn/config"
	"gn/issues/issue"
	"gn/logger"
	"gn/repo"
)

func QuerySingle(conf *config.Wrapper, details []repo.Details, u *url.URL, issueID string) (*issue.Details, error) {
	match, projectPath, err := getMatchingConfig(conf, details, u)
	if err != nil {
		logger.Log.Errorf("Failed to get matching config: %s", err)

		return nil, err
	}

	if match.URL.Host == "github.com" {
		return issue.QueryGitHub(match, projectPath, issueID)
	}

	return issue.QueryGitLab(match, projectPath, issueID)
}
