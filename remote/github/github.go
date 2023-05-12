package github

// Details implements the remote.Details interface for a GitHub remote.
type Details struct {
	Token     string
	TokenName string
	Username  string
}

// GetToken implements the remote.Details interface.
func (hub Details) GetToken() string {
	return hub.Token
}

// GetTokenName implements the remote.Details interface.
func (hub Details) GetTokenName() string {
	return hub.TokenName
}

// GetUsername implements the remote.Details interface.
func (hub Details) GetUsername() string {
	return hub.Username
}
