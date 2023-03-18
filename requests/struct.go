package requests

type GraphqlQuery struct {
	Variables map[string]string `json:"variables"`
	Query     string            `json:"query"`
}
