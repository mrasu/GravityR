package models

import (
	"github.com/mrasu/GravityR/database/db_models"
	"github.com/mrasu/GravityR/database/db_models/builders"
	"github.com/mrasu/GravityR/html/viewmodel"
)

type ExplainAnalyzeTree struct {
	Root *ExplainAnalyzeTreeNode
}

func (eat *ExplainAnalyzeTree) ToSingleTableResults() []*db_models.SingleTableExplainResult {
	return builders.BuildSingleTableExplainResults(eat.Root)
}

func (eat *ExplainAnalyzeTree) ToViewModel() []*viewmodel.VmMysqlExplainAnalyzeNode {
	var res []*viewmodel.VmMysqlExplainAnalyzeNode
	for _, n := range eat.Root.Children {
		res = append(res, n.ToViewModel())
	}
	return res
}
