package issues

import (
	"net/url"

	"github.com/glad-dev/gin/config"
	"github.com/glad-dev/gin/issues/discussion"
	"github.com/glad-dev/gin/log"
	"github.com/glad-dev/gin/repo"
)

// QuerySingle returns the discussion associated with the given issueID.
func QuerySingle(conf *config.Wrapper, details []repo.Details, u *url.URL, issueID string) (*discussion.Details, error) {
	match, projectPath, err := getMatchingConfig(conf, details, u)
	if err != nil {
		log.Error("Failed to get matching config", "error", err)

		return nil, err
	}

	if match.URL.Host == "github.com" {
		return discussion.QueryGitHub(match, projectPath, issueID)
	}

	return discussion.QueryGitLab(match, projectPath, issueID)
}
