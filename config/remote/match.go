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

func (l *Match) GetURL() string {
	return l.URL.String()
}

func (l *Match) GetToken() string {
	return l.Token
}
