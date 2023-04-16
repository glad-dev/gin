package remote

import (
	"errors"
	"net/url"
)

var ErrMultipleRepoDetails = errors.New("config contains multiple matching configs")

type Details interface {
	GetToken() string
	GetTokenName() string
	GetUsername() string

	Init(*url.URL) (Details, error)
	CheckTokenScope(*url.URL) (string, error)
}
