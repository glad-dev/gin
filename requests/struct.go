package requests

// Query contains the query string and its list of variables.
type Query struct {
	Variables map[string]interface{} `json:"variables"`
	Query     string                 `json:"query"`
}
