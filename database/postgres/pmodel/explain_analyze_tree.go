package pmodel

import (
	"github.com/mrasu/GravityR/database/dmodel"
	"github.com/mrasu/GravityR/database/dservice"
	"github.com/mrasu/GravityR/html/viewmodel"
)

type ExplainAnalyzeTree struct {
	Root        *ExplainAnalyzeTreeNode
	SummaryText string
}

func (eat *ExplainAnalyzeTree) ToSingleTableResults() []*dmodel.SingleTableExplainResult {
	return dservice.BuildSingleTableExplainResults(eat.Root)
}

func (eat *ExplainAnalyzeTree) ToViewModel() []*viewmodel.VmPostgresExplainAnalyzeNode {
	var res []*viewmodel.VmPostgresExplainAnalyzeNode
	for _, n := range eat.Root.Children {
		res = append(res, n.ToViewModel())
	}
	return res
}
