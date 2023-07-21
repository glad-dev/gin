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
	switch t {
	case Github:
		// GitHub can't be self-hosted, thus the library we use doesn't require an API URL.
		return "", nil
	case Gitlab:
		return u.JoinPath("/api/v4").String(), nil
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
