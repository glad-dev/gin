package remote

import (
	"errors"
	"net/url"
)

type Type uint8

const (
	Github Type = iota
	Gitlab
	Bitbucket
)

func (t Type) ApiURL(u *url.URL) (string, error) { //nolint:revive
	if u == nil {
		return "", errors.New("nil url passed to ApiURL")
	}

	switch t {
	case Github:
		return "https://api.github.com/graphql", nil
	case Gitlab:
		return u.JoinPath("/api/graphql").String(), nil
	case Bitbucket:
		return "", errors.New("bitbucket is not yet implemented")
	}

	return "", errors.New("invalid type")
}

func (t Type) String() string {
	switch t {
	case Github:
		return "GitHub"
	case Gitlab:
		return "GitLab"
	case Bitbucket:
		return "Bitbucket"
	}

	return "Unknown type"
}
