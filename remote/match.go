package remote

import (
	"gn/logger"
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
func (m *Match) ApiURL() string { //nolint:revive
	if m.URL.Host == "github.com" {
		logger.Log.Info("Called ApiURL with a github host")

		return "https://api.github.com/graphql"
	}

	return m.URL.JoinPath("/api/v4").String()
}
