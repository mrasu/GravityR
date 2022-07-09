package models_test

import (
	"github.com/mrasu/GravityR/database/mysql/models"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"
	"testing"
)

func TestExplainAnalyzeTree_ToSingleTableTrees(t *testing.T) {
	tests := []struct {
		name          string
		nodes         []*models.ExplainAnalyzeTreeNode
		expectedTrees []*models.SingleTableExplainAnalyzeTree
	}{
		{
			name: "one line",
			nodes: []*models.ExplainAnalyzeTreeNode{
				{
					AnalyzeResultLine: buildLine(
						"-> Table scan on todos  (cost=0 rows=0) (actual time=0..0.2 rows=0 loops=2)",
						"todos", 0.2, 2,
					),
				},
			},
			expectedTrees: []*models.SingleTableExplainAnalyzeTree{
				{
					TableName:          "todos",
					EstimatedTotalTime: 0.2 * 2,
				},
			},
		},
		{
			name: "multiple lines",
			nodes: []*models.ExplainAnalyzeTreeNode{
				{
					AnalyzeResultLine: buildLine(
						"    -> Filter: (todos.`description` = 'a')  (cost=0 rows=0) (actual time=0..0.2 rows=0 loops=2)",
						"todos", 0.2, 2,
					),
					Children: []*models.ExplainAnalyzeTreeNode{
						{
							AnalyzeResultLine: buildLine(
								"        -> Table scan on todos  (cost=0 rows=0) (actual time=0..0.3 rows=0 loops=3)",
								"todos", 0.3, 3,
							),
						},
					},
				},
			},
			expectedTrees: []*models.SingleTableExplainAnalyzeTree{
				{
					TableName:          "todos",
					EstimatedTotalTime: 0.2 * 2,
				},
			},
		},
		{
			name: "multiple tables",
			nodes: []*models.ExplainAnalyzeTreeNode{
				{
					AnalyzeResultLine: buildLine(
						"-> Nested loop inner join  (cost=0 rows=0) (actual time=0..0.1 rows=0 loops=1)",
						"", 0.1, 1,
					),
					Children: []*models.ExplainAnalyzeTreeNode{
						{
							AnalyzeResultLine: buildLine(
								"    -> Filter: (todos.`description` = 'a')  (cost=0 rows=0) (actual time=0..0.2 rows=0 loops=2)",
								"todos", 0.2, 2,
							),
							Children: []*models.ExplainAnalyzeTreeNode{
								{
									AnalyzeResultLine: buildLine(
										"        -> Table scan on todos  (cost=0 rows=0) (actual time=0..0.3 rows=0 loops=3)",
										"todos", 0.3, 3,
									),
								},
							},
						},
						{
							AnalyzeResultLine: buildLine(
								"    -> Single-row index lookup on users using PRIMARY (id=todos.user_id)  (cost=0 rows=0) (actual time=0..0.4 rows=0 loops=4)",
								"users", 0.4, 4,
							),
						},
					},
				},
			},
			expectedTrees: []*models.SingleTableExplainAnalyzeTree{
				{
					TableName:          "todos",
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
			eat := &models.ExplainAnalyzeTree{
				Root: &models.ExplainAnalyzeTreeNode{
					AnalyzeResultLine: &models.ExplainAnalyzeResultLine{},
					Children:          tt.nodes,
				},
			}

			trees := eat.ToSingleTableTrees()
			assert.Equal(t, tt.expectedTrees, trees)
		})
	}
}

func buildLine(text, table string, avg float64, loop int64) *models.ExplainAnalyzeResultLine {
	return &models.ExplainAnalyzeResultLine{
		Text:            text,
		TableName:       table,
		ActualTimeAvg:   null.FloatFrom(avg),
		ActualLoopCount: null.IntFrom(loop),
	}
}
