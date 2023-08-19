package issues

import (
	"net/url"

	"github.com/glad-dev/gin/config"
	"github.com/glad-dev/gin/issues/list"
	"github.com/glad-dev/gin/logger"
	"github.com/glad-dev/gin/repo"
)

// QueryList returns a list of issues associated with a repository.
func QueryList(conf *config.Wrapper, details []repo.Details, u *url.URL) ([]list.Issue, error) {
	match, projectPath, err := getMatchingConfig(conf, details, u)
	if err != nil {
		logger.Log.Error("Failed to get matching config", "error", err)

		return nil, err
	}

	if match.URL.Host == "github.com" {
		return list.QueryGitHub(match, projectPath)
	}

	return list.QueryGitLab(match, projectPath)
}
