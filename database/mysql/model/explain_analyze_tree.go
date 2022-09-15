package model

import (
	"github.com/mrasu/GravityR/database/common_model"
	"github.com/mrasu/GravityR/database/common_model/builder"
	"github.com/mrasu/GravityR/html/viewmodel"
)

type ExplainAnalyzeTree struct {
	Root *ExplainAnalyzeTreeNode
}

func (eat *ExplainAnalyzeTree) ToSingleTableResults() []*common_model.SingleTableExplainResult {
	return builder.BuildSingleTableExplainResults(eat.Root)
}

func (eat *ExplainAnalyzeTree) ToViewModel() []*viewmodel.VmMysqlExplainAnalyzeNode {
	var res []*viewmodel.VmMysqlExplainAnalyzeNode
	for _, n := range eat.Root.Children {
		res = append(res, n.ToViewModel())
	}
	return res
}
