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

// ApiURL returns the API URL.
func (l *Match) ApiURL() string { //nolint:revive
	if l.URL.Host == "github.com" {
		return "https://api.github.com/graphql"
	}

	return l.URL.JoinPath("/api/graphql").String()
}
