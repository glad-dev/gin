package issues

import (
	"net/url"

	"gn/config"
	"gn/issues/discussion"
	"gn/logger"
	"gn/repo"
)

// QuerySingle returns the discussion associated with the given issueID.
func QuerySingle(conf *config.Wrapper, details []repo.Details, u *url.URL, issueID string) (*discussion.Details, error) {
	match, projectPath, err := getMatchingConfig(conf, details, u)
	if err != nil {
		logger.Log.Errorf("Failed to get matching config: %s", err)

		return nil, err
	}

	if match.URL.Host == "github.com" {
		return discussion.QueryGitHub(match, projectPath, issueID)
	}

	return discussion.QueryGitLab(match, projectPath, issueID)
}
