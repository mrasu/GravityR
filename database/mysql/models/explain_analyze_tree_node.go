package models

import "github.com/mrasu/GravityR/html/viewmodel"

type ExplainAnalyzeTreeNode struct {
	AnalyzeResultLine *ExplainAnalyzeResultLine
	Children          []*ExplainAnalyzeTreeNode
}

func (n *ExplainAnalyzeTreeNode) ToViewModel() *viewmodel.VmExplainAnalyzeNode {
	vm := &viewmodel.VmExplainAnalyzeNode{
		Text:                     n.AnalyzeResultLine.Text,
		Title:                    n.AnalyzeResultLine.Title(),
		TableName:                n.AnalyzeResultLine.TableName,
		EstimatedInitCost:        n.AnalyzeResultLine.EstimatedInitCost,
		EstimatedCost:            n.AnalyzeResultLine.EstimatedCost,
		EstimatedReturnedRows:    n.AnalyzeResultLine.EstimatedReturnedRows,
		ActualTimeFirstRow:       n.AnalyzeResultLine.ActualTimeFirstRow,
		ActualTimeAvg:            n.AnalyzeResultLine.ActualTimeAvg,
		ActualReturnedRows:       n.AnalyzeResultLine.ActualReturnedRows,
		ActualLoopCount:          n.AnalyzeResultLine.ActualLoopCount,
		EstimatedActualTotalTime: n.AnalyzeResultLine.EstimatedTotalTime(),
	}
	for _, c := range n.Children {
		vm.Children = append(vm.Children, c.ToViewModel())
	}

	return vm
}
