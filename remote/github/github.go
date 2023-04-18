package github

type Details struct {
	Token     string
	TokenName string
	Username  string
}

func (hub Details) GetToken() string {
	return hub.Token
}

func (hub Details) GetTokenName() string {
	return hub.TokenName
}

func (hub Details) GetUsername() string {
	return hub.Username
}
