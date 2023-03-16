package requests

type GraphqlQuery struct {
	Query     string            `json:"query"`
	Variables map[string]string `json:"variables"`
}
