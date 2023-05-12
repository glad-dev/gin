package gitlab

import "net/url"

// Details implements the remote.Details interface for a GitLab remote.
type Details struct {
	Token     string
	TokenName string
	Username  string
}

// GetToken implements the remote.Details interface.
func (lab Details) GetToken() string {
	return lab.Token
}

// GetTokenName implements the remote.Details interface.
func (lab Details) GetTokenName() string {
	return lab.TokenName
}

// GetUsername implements the remote.Details interface.
func (lab Details) GetUsername() string {
	return lab.Username
}

// ApiURL returns the REST API of a given GitLab host.
func ApiURL(u *url.URL) string { //nolint:revive
	return u.JoinPath("/api/v4").String()
}
