package model

import (
	"github.com/mrasu/GravityR/database/common_model/builder"
	"github.com/mrasu/GravityR/html/viewmodel"
	"strings"
)

type ExplainAnalyzeTreeNode struct {
	AnalyzeResultNode *ExplainAnalyzeResultNode
	Children          []*ExplainAnalyzeTreeNode
	SpaceSize         int
}

func (n *ExplainAnalyzeTreeNode) ToViewModel() *viewmodel.VmPostgresExplainAnalyzeNode {
	text := strings.Join(n.AnalyzeResultNode.Lines, "\n")

	vm := &viewmodel.VmPostgresExplainAnalyzeNode{
		Text:                  text,
		Title:                 n.AnalyzeResultNode.Title(),
		TableName:             n.AnalyzeResultNode.TableName,
		EstimatedInitCost:     n.AnalyzeResultNode.EstimatedInitCost,
		EstimatedCost:         n.AnalyzeResultNode.EstimatedCost,
		EstimatedReturnedRows: n.AnalyzeResultNode.EstimatedReturnedRows,
		EstimatedWidth:        n.AnalyzeResultNode.EstimatedWidth,
		ActualTimeFirstRow:    n.AnalyzeResultNode.ActualTimeFirstRow,
		ActualTimeAvg:         n.AnalyzeResultNode.ActualTimeAvg,
		ActualReturnedRows:    n.AnalyzeResultNode.ActualReturnedRows,
		ActualLoopCount:       n.AnalyzeResultNode.ActualLoopCount,
	}
	for _, c := range n.Children {
		vm.Children = append(vm.Children, c.ToViewModel())
	}

	return vm
}

func (n *ExplainAnalyzeTreeNode) TableName() string {
	return n.AnalyzeResultNode.TableName
}

func (n *ExplainAnalyzeTreeNode) EstimatedTotalTime() float64 {
	return n.AnalyzeResultNode.EstimatedTotalTime()
}

func (n *ExplainAnalyzeTreeNode) GetChildren() []builder.ExplainNode {
	var res []builder.ExplainNode
	for _, c := range n.Children {
		res = append(res, c)
	}
	return res
}
