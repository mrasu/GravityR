package hasura

type RunSQLResponse struct {
	ResultType string     `json:"result_type"`
	Result     [][]string `json:"result"`
}
