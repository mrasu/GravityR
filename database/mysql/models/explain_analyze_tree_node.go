package models

import (
	"github.com/mrasu/GravityR/database/db_models/builders"
	"github.com/mrasu/GravityR/html/viewmodel"
)

type ExplainAnalyzeTreeNode struct {
	AnalyzeResultLine *ExplainAnalyzeResultLine
	Children          []*ExplainAnalyzeTreeNode
}

func (n *ExplainAnalyzeTreeNode) ToViewModel() *viewmodel.VmMysqlExplainAnalyzeNode {
	vm := &viewmodel.VmMysqlExplainAnalyzeNode{
		Text:                  n.AnalyzeResultLine.Text,
		Title:                 n.AnalyzeResultLine.Title(),
		TableName:             n.AnalyzeResultLine.TableName,
		EstimatedInitCost:     n.AnalyzeResultLine.EstimatedInitCost,
		EstimatedCost:         n.AnalyzeResultLine.EstimatedCost,
		EstimatedReturnedRows: n.AnalyzeResultLine.EstimatedReturnedRows,
		ActualTimeFirstRow:    n.AnalyzeResultLine.ActualTimeFirstRow,
		ActualTimeAvg:         n.AnalyzeResultLine.ActualTimeAvg,
		ActualReturnedRows:    n.AnalyzeResultLine.ActualReturnedRows,
		ActualLoopCount:       n.AnalyzeResultLine.ActualLoopCount,
	}
	for _, c := range n.Children {
		vm.Children = append(vm.Children, c.ToViewModel())
	}

	return vm
}

func (n *ExplainAnalyzeTreeNode) TableName() string {
	return n.AnalyzeResultLine.TableName
}

func (n *ExplainAnalyzeTreeNode) EstimatedTotalTime() float64 {
	return n.AnalyzeResultLine.EstimatedTotalTime()
}

func (n *ExplainAnalyzeTreeNode) GetChildren() []builders.ExplainNode {
	var res []builders.ExplainNode
	for _, c := range n.Children {
		res = append(res, c)
	}
	return res
}
