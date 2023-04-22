package remote

import (
	"net/url"
)

// Details interface, implemented by github.Details and gitlab.Details.
type Details interface {
	GetToken() string
	GetTokenName() string
	GetUsername() string

	Init(*url.URL) (Details, error)
	CheckTokenScope(*url.URL) (string, error)
}
