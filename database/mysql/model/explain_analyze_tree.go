package model

import (
	"github.com/mrasu/GravityR/database"
	"github.com/mrasu/GravityR/database/service"
	"github.com/mrasu/GravityR/html/viewmodel"
)

type ExplainAnalyzeTree struct {
	Root *ExplainAnalyzeTreeNode
}

func (eat *ExplainAnalyzeTree) ToSingleTableResults() []*database.SingleTableExplainResult {
	return service.BuildSingleTableExplainResults(eat.Root)
}

func (eat *ExplainAnalyzeTree) ToViewModel() []*viewmodel.VmMysqlExplainAnalyzeNode {
	var res []*viewmodel.VmMysqlExplainAnalyzeNode
	for _, n := range eat.Root.Children {
		res = append(res, n.ToViewModel())
	}
	return res
}
