package issues

import (
	"net/url"

	"github.com/glad-dev/gin/configuration"
	"github.com/glad-dev/gin/issues/list"
	"github.com/glad-dev/gin/log"
	rt "github.com/glad-dev/gin/remote/type"
	"github.com/glad-dev/gin/repository"
)

// QueryList returns a list of issues associated with a repository. Uses the GraphQL API since it is faster than the
// REST API.
func QueryList(conf *configuration.Config, details []repository.Details, u *url.URL) ([]list.Issue, error) {
	match, projectPath, err := getMatchingConfig(conf, details, u)
	if err != nil {
		log.Error("Failed to get matching config", "error", err)

		return nil, err
	}

	if match.Type == rt.Github {
		return list.QueryGitHub(match, projectPath)
	}

	return list.QueryGitLab(match, projectPath)
}
