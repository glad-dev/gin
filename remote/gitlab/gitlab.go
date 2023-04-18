package gitlab

type Details struct {
	Token     string
	TokenName string
	Username  string
}

func (lab Details) GetToken() string {
	return lab.Token
}

func (lab Details) GetTokenName() string {
	return lab.TokenName
}

func (lab Details) GetUsername() string {
	return lab.Username
}
