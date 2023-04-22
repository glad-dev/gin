package gitlab

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
