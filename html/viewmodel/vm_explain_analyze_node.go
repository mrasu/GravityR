package viewmodel

import "gopkg.in/guregu/null.v4"

type VmExplainAnalyzeNode struct {
	Text                     string     `json:"text"`
	Title                    string     `json:"title"`
	TableName                string     `json:"tableName"`
	EstimatedInitCost        null.Float `json:"estimatedInitCost"`
	EstimatedCost            null.Float `json:"estimatedCost"`
	EstimatedReturnedRows    null.Int   `json:"estimatedReturnedRows"`
	ActualTimeFirstRow       null.Float `json:"actualTimeFirstRow"`
	ActualTimeAvg            null.Float `json:"actualTimeAvg"`
	ActualReturnedRows       null.Int   `json:"actualReturnedRows"`
	ActualLoopCount          null.Int   `json:"actualLoopCount"`
	EstimatedActualTotalTime float64    `json:"estimatedActualTotalTime"`

	Children []*VmExplainAnalyzeNode `json:"children"`
}
