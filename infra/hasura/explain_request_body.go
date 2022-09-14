package hasura

type ExplainRequestBody struct {
	Query *Query            `json:"query"`
	User  map[string]string `json:"user"`
}
