package requests

type GraphqlQuery struct {
	Variables map[string]string `json:"variables"`
	Query     string            `json:"query"`
}

type GitHubCommentQuery struct {
	Variables map[string]interface{} `json:"variables"`
	Query     string                 `json:"query"`
}

type ConfigInterface interface {
	GetURL() string
	GetToken() string
}
