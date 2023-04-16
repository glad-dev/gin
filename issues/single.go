package issues

import (
	"net/url"

	"gn/config"
	"gn/issues/single"
	"gn/logger"
	"gn/repo"
)

func QuerySingle(conf *config.Wrapper, details []repo.Details, u *url.URL, issueID string) (*single.IssueDetails, error) {
	match, projectPath, err := getMatchingConfig(conf, details, u)
	if err != nil {
		logger.Log.Errorf("Failed to get matching config: %s", err)

		return nil, err
	}

	if match.URL.Host == "github.com" {
		return single.QuerySingleGitHub(match, projectPath, issueID)
	}

	return single.QuerySingleGitLab(match, projectPath, issueID)
}
