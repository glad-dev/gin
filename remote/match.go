package remote

import (
	"net/url"
)

// Match contains the information needed to request data from a remote.
type Match struct {
	URL       url.URL
	Token     string
	Username  string
	TokenName string
}
