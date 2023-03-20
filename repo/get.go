package repo

import (
	"errors"
	"net/url"
	"strings"

	"github.com/go-git/go-git/v5"
)

type Details struct {
	url         url.URL
	projectPath string
}

func (details *Details) ProjectPath() string {
	return details.projectPath
}

func (details *Details) URL() url.URL {
	return details.url
}

func Get(path string) ([]Details, error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	remotes, err := r.Remotes()
	if err != nil {
		return nil, err
	}

	repos := make([]Details, 0)
	for _, remote := range remotes {
		for _, origin := range remote.Config().URLs {
			if !strings.HasPrefix(origin, "git@") {
				return nil, errors.New("origin does not start with 'git@'")
			}

			if !strings.HasSuffix(origin, ".git") {
				return nil, errors.New("origin does not end with '.git'")
			}

			if strings.Count(origin, ":") != 1 {
				return nil, errors.New("origin contains an invalid amount of ':'")
			}

			u, err := url.Parse(origin[len("git@") : len(origin)-len(".git")])
			if err != nil {
				return nil, err
			}

			repos = append(repos, Details{
				url: url.URL{
					Host: u.Scheme,
				},
				projectPath: u.Opaque,
			})
		}
	}

	return repos, nil
}
