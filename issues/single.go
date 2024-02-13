package issues

import (
	remotetype "github.com/glad-dev/gin/remote/type"
	"net/url"

	"github.com/glad-dev/gin/configuration"
	"github.com/glad-dev/gin/issues/discussion"
	"github.com/glad-dev/gin/log"
	"github.com/glad-dev/gin/repository"
)

// QuerySingle returns the discussion associated with the given issueID.
func QuerySingle(conf *configuration.Config, details []repository.Details, u *url.URL, issueID string) (*discussion.Details, error) {
	match, projectPath, err := getMatchingConfig(conf, details, u)
	if err != nil {
		log.Error("Failed to get matching config", "error", err)

		return nil, err
	}

	if match.Type == remotetype.Github {
		return discussion.QueryGitHub(match, projectPath, issueID)
	}

	return discussion.QueryGitLab(match, projectPath, issueID)
}
