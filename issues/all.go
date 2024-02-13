package issues

import (
	"net/url"

	"github.com/glad-dev/gin/config"
	"github.com/glad-dev/gin/issues/list"
	"github.com/glad-dev/gin/log"
	"github.com/glad-dev/gin/repository"
)

// QueryList returns a list of issues associated with a repository. Uses the GraphQL API since it is faster than the
// REST API.
func QueryList(conf *config.Wrapper, details []repository.Details, u *url.URL) ([]list.Issue, error) {
	match, projectPath, err := getMatchingConfig(conf, details, u)
	if err != nil {
		log.Error("Failed to get matching config", "error", err)

		return nil, err
	}

	if match.URL.Host == "github.com" {
		return list.QueryGitHub(match, projectPath)
	}

	return list.QueryGitLab(match, projectPath)
}
