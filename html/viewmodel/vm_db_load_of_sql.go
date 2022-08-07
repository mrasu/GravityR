package viewmodel

type VmDbLoadOfSql struct {
	Sql         string  `json:"sql"`
	LoadMax     float64 `json:"loadMax"`
	LoadSum     float64 `json:"loadSum"`
	TokenizedId string  `json:"tokenizedId"`
}
