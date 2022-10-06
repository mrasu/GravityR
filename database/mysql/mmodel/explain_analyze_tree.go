package mmodel

import (
	"github.com/mrasu/GravityR/database/dmodel"
	"github.com/mrasu/GravityR/database/dservice"
	"github.com/mrasu/GravityR/html/viewmodel"
)

type ExplainAnalyzeTree struct {
	Root *ExplainAnalyzeTreeNode
}

func (eat *ExplainAnalyzeTree) ToSingleTableResults() []*dmodel.SingleTableExplainResult {
	return dservice.BuildSingleTableExplainResults(eat.Root)
}

func (eat *ExplainAnalyzeTree) ToViewModel() []*viewmodel.VmMysqlExplainAnalyzeNode {
	var res []*viewmodel.VmMysqlExplainAnalyzeNode
	for _, n := range eat.Root.Children {
		res = append(res, n.ToViewModel())
	}
	return res
}
