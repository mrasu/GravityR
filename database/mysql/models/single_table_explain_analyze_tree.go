package models

import (
	"fmt"
)

type SingleTableExplainAnalyzeTree struct {
	TableName string

	EstimatedTotalTime float64
}

func (st *SingleTableExplainAnalyzeTree) String() string {
	txt := fmt.Sprintf(
		"SingleTableExplainAnalyzeTree(TableName: %s, EstimatedTotalTime: %f)",
		st.TableName,
		st.EstimatedTotalTime,
	)
	return txt
}
