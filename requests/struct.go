package requests

type GraphqlQuery struct {
	Variables map[string]string `json:"variables"`
	Query     string            `json:"query"`
}

type scopes struct {
	Data struct {
		Viewer struct {
			Scopes []string `json:"scopes"`
		} `json:"viewer"`
	} `json:"data"`
}
