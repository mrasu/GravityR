package model

import (
	"github.com/mrasu/GravityR/database/common_model"
	"github.com/mrasu/GravityR/database/common_model/builder"
	"github.com/mrasu/GravityR/html/viewmodel"
)

type ExplainAnalyzeTree struct {
	Root        *ExplainAnalyzeTreeNode
	SummaryText string
}

func (eat *ExplainAnalyzeTree) ToSingleTableResults() []*common_model.SingleTableExplainResult {
	return builder.BuildSingleTableExplainResults(eat.Root)
}

func (eat *ExplainAnalyzeTree) ToViewModel() []*viewmodel.VmPostgresExplainAnalyzeNode {
	var res []*viewmodel.VmPostgresExplainAnalyzeNode
	for _, n := range eat.Root.Children {
		res = append(res, n.ToViewModel())
	}
	return res
}
