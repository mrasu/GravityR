package hasura

type ExplainResponse struct {
	Field string   `json:"field"`
	Plan  []string `json:"plan"`
	SQL   string   `json:"sql"`
}
