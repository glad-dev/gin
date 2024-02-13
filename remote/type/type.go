package rt

import (
	"errors"
	"net/url"
)

type Type uint8

const (
	Github Type = iota
	Gitlab
)

func (t Type) GraphqlAPIURL(u *url.URL) (string, error) { //nolint:revive
	if u == nil {
		return "", errors.New("nil url passed to ApiURL")
	}

	switch t {
	case Github:
		return "https://api.github.com/graphql", nil
	case Gitlab:
		return u.JoinPath("/api/graphql").String(), nil
	}

	return "", errors.New("invalid type")
}

func (t Type) RestAPIURL(u *url.URL) (string, error) {
	if u == nil {
		return "", errors.New("nil url passed to ApiURL")
	}

	switch t {
	case Github:
		return "https://api.github.com/", nil
	case Gitlab:
		return u.JoinPath("/api/v4").String(), nil
	}

	return "", errors.New("invalid type")
}

func (t Type) String() string {
	switch t {
	case Github:
		return "GitHub"
	case Gitlab:
		return "GitLab"
	}

	return "Unknown type"
}
