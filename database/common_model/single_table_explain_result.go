package common_model

import (
	"fmt"
)

type SingleTableExplainResult struct {
	TableName string

	EstimatedTotalTime float64
}

func (st *SingleTableExplainResult) String() string {
	txt := fmt.Sprintf(
		"SingleTableExplainResult(TableName: %s, EstimatedTotalTime: %f)",
		st.TableName,
		st.EstimatedTotalTime,
	)
	return txt
}
