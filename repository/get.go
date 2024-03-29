package repository

import (
	"errors"
	"net/url"
	"strings"

	"github.com/glad-dev/gin/log"

	"github.com/go-git/go-git/v5"
)

// Details contains the remote's URL and project path.
type Details struct {
	URL         url.URL
	ProjectPath string
}

// Get opens the git repository at the given path and returns its list of repository.Details.
func Get(path string) ([]Details, error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		log.Error("Failed to open repository.", "error", err, "path", path)

		return nil, err
	}

	remotes, err := r.Remotes()
	if err != nil {
		log.Error("Failed to get remotes.", "error", err, "path", path)

		return nil, err
	}

	if len(remotes) == 0 {
		log.Error("Repository does not contains any remotes.", "path", path)

		return nil, errors.New("repository does not contain any remotes")
	} else if len(remotes) > 1 {
		log.Error("Repository contains multiple remotes.", "path", path)
	}

	repos := make([]Details, 0)
	for _, remote := range remotes {
		for _, origin := range remote.Config().URLs {
			if !strings.HasPrefix(origin, "git@") {
				log.Error("Origin does have 'git@' prefix.", "origin", origin)

				return nil, errors.New("origin does have 'git@' prefix")
			}

			suffixLength := 0
			if strings.HasSuffix(origin, ".git") {
				suffixLength = len(".git")
			}

			if strings.Count(origin, ":") != 1 {
				log.Error("Origin has invalid amount of ':'.", "origin", origin)

				return nil, errors.New("origin contains an invalid amount of ':'")
			}

			u, err := url.Parse(origin[len("git@") : len(origin)-suffixLength])
			if err != nil {
				log.Error("Failed to parse the git URL.", "error", err, "url", origin[len("git@"):len(origin)-len(".git")])

				return nil, err
			}

			repos = append(repos, Details{
				URL: url.URL{
					Host: u.Scheme,
				},
				ProjectPath: u.Opaque,
			})
		}
	}

	return repos, nil
}
