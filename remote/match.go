package remote

import (
	"net/url"
)

type Match struct {
	URL       url.URL
	Token     string
	Username  string
	TokenName string
}

func (l *Match) GetApiURL() string {
	if l.URL.Host == "github.com" {
		return "https://api.github.com/graphql"
	}

	return l.URL.JoinPath("/api/graphql").String()
}
