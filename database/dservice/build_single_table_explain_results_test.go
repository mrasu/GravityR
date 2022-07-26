package dservice_test

import (
	"github.com/mrasu/GravityR/database/dmodel"
	"github.com/mrasu/GravityR/database/dservice"
	"github.com/stretchr/testify/assert"
	"testing"
)

type explainNode struct {
	tableName          string
	estimatedTotalTime float64

	children []*explainNode
}

func (en *explainNode) TableName() string {
	return en.tableName
}
func (en *explainNode) EstimatedTotalTime() float64 {
	return en.estimatedTotalTime
}
func (en *explainNode) GetChildren() []dservice.ExplainNode {
	var res []dservice.ExplainNode
	for _, c := range en.children {
		res = append(res, c)
	}
	return res
}

func TestExplainAnalyzeTree_ToSingleTableTrees(t *testing.T) {
	tests := []struct {
		name          string
		nodes         []*explainNode
		expectedTrees []*dmodel.SingleTableExplainResult
	}{
		{
			name: "one line",
			nodes: []*explainNode{
				{
					tableName:          "tasks",
					estimatedTotalTime: 0.4,
				},
			},
			expectedTrees: []*dmodel.SingleTableExplainResult{
				{
					TableName:          "tasks",
					EstimatedTotalTime: 0.4,
				},
			},
		},
		{
			name: "multiple lines",
			nodes: []*explainNode{
				{
					tableName:          "tasks",
					estimatedTotalTime: 0.4,
					children: []*explainNode{
						{
							tableName:          "tasks",
							estimatedTotalTime: 0.9,
						},
					},
				},
			},
			expectedTrees: []*dmodel.SingleTableExplainResult{
				{
					TableName:          "tasks",
					EstimatedTotalTime: 0.2 * 2,
				},
			},
		},
		{
			name: "multiple refTables",
			nodes: []*explainNode{
				{
					tableName:          "",
					estimatedTotalTime: 0.1,
					children: []*explainNode{
						{
							tableName:          "tasks",
							estimatedTotalTime: 0.4,
							children: []*explainNode{
								{
									tableName:          "tasks",
									estimatedTotalTime: 0.9,
								},
							},
						},
						{
							children: []*explainNode{
								{
									tableName:          "users",
									estimatedTotalTime: 1.6,
								},
							},
						},
					},
				},
			},
			expectedTrees: []*dmodel.SingleTableExplainResult{
				{
					TableName:          "tasks",
					EstimatedTotalTime: 0.2 * 2,
				},
				{
					TableName:          "users",
					EstimatedTotalTime: 0.4 * 4,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := &explainNode{children: tt.nodes}
			trees := dservice.BuildSingleTableExplainResults(root)

			assert.Equal(t, tt.expectedTrees, trees)
		})
	}
}
