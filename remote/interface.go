package remote

import (
	"net/url"
)

type Details interface {
	GetToken() string
	GetTokenName() string
	GetUsername() string

	Init(*url.URL) (Details, error)
	CheckTokenScope(*url.URL) (string, error)
}
