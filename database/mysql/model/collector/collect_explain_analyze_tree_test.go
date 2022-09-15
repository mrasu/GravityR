package collector_test

import (
	"github.com/mrasu/GravityR/database/mysql/model"
	"github.com/mrasu/GravityR/database/mysql/model/collector"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"
	"testing"
)

func TestCollectExplainAnalyzeTree(t *testing.T) {
	tests := []struct {
		name          string
		explainResult string
		expectedTree  *model.ExplainAnalyzeTree
	}{
		{
			name: "simple(SELECT name FROM users)",
			explainResult: `
-> Index scan on users using id  (cost=6735.18 rows=63830) (actual time=0.062..9.176 rows=65536 loops=1)
`,
			expectedTree: &model.ExplainAnalyzeTree{
				Root: &model.ExplainAnalyzeTreeNode{
					AnalyzeResultLine: &model.ExplainAnalyzeResultLine{},
					Children: []*model.ExplainAnalyzeTreeNode{
						{
							AnalyzeResultLine: &model.ExplainAnalyzeResultLine{
								Text:                  "-> Index scan on users using id  (cost=6735.18 rows=63830) (actual time=0.062..9.176 rows=65536 loops=1)",
								TableName:             "users",
								EstimatedInitCost:     null.FloatFromPtr(nil),
								EstimatedCost:         null.FloatFrom(6735.18),
								EstimatedReturnedRows: null.IntFrom(63830),
								ActualTimeFirstRow:    null.FloatFrom(0.062),
								ActualTimeAvg:         null.FloatFrom(9.176),
								ActualReturnedRows:    null.IntFrom(65536),
								ActualLoopCount:       null.IntFrom(1),
							},
						},
					},
				},
			},
		},
		{
			name: "single join(SELECT name FROM users INNER JOIN todos ON ...)",
			explainResult: `
-> Nested loop inner join  (cost=6378854.77 rows=6361711) (actual time=0.147..2483.756 rows=6553600 loops=1)
    -> Table scan on todos  (cost=659260.39 rows=6361711) (actual time=0.074..1138.883 rows=6553600 loops=1)
    -> Single-row index lookup on users using PRIMARY (id=todos.user_id)  (cost=0.80 rows=1) (actual time=0.000..0.000 rows=1 loops=6553600)
`,
			expectedTree: &model.ExplainAnalyzeTree{
				Root: &model.ExplainAnalyzeTreeNode{
					AnalyzeResultLine: &model.ExplainAnalyzeResultLine{},
					Children: []*model.ExplainAnalyzeTreeNode{
						{
							AnalyzeResultLine: &model.ExplainAnalyzeResultLine{
								Text:                  "-> Nested loop inner join  (cost=6378854.77 rows=6361711) (actual time=0.147..2483.756 rows=6553600 loops=1)",
								TableName:             "",
								EstimatedInitCost:     null.FloatFromPtr(nil),
								EstimatedCost:         null.FloatFrom(6378854.77),
								EstimatedReturnedRows: null.IntFrom(6361711),
								ActualTimeFirstRow:    null.FloatFrom(0.147),
								ActualTimeAvg:         null.FloatFrom(2483.756),
								ActualReturnedRows:    null.IntFrom(6553600),
								ActualLoopCount:       null.IntFrom(1),
							},
							Children: []*model.ExplainAnalyzeTreeNode{
								{
									AnalyzeResultLine: &model.ExplainAnalyzeResultLine{
										Text:                  "    -> Table scan on todos  (cost=659260.39 rows=6361711) (actual time=0.074..1138.883 rows=6553600 loops=1)",
										TableName:             "todos",
										EstimatedInitCost:     null.FloatFromPtr(nil),
										EstimatedCost:         null.FloatFrom(659260.39),
										EstimatedReturnedRows: null.IntFrom(6361711),
										ActualTimeFirstRow:    null.FloatFrom(0.074),
										ActualTimeAvg:         null.FloatFrom(1138.883),
										ActualReturnedRows:    null.IntFrom(6553600),
										ActualLoopCount:       null.IntFrom(1),
									},
								},
								{
									AnalyzeResultLine: &model.ExplainAnalyzeResultLine{
										Text:                  "    -> Single-row index lookup on users using PRIMARY (id=todos.user_id)  (cost=0.80 rows=1) (actual time=0.000..0.000 rows=1 loops=6553600)",
										TableName:             "users",
										EstimatedInitCost:     null.FloatFromPtr(nil),
										EstimatedCost:         null.FloatFrom(0.80),
										EstimatedReturnedRows: null.IntFrom(1),
										ActualTimeFirstRow:    null.FloatFrom(0.0),
										ActualTimeAvg:         null.FloatFrom(0.0),
										ActualReturnedRows:    null.IntFrom(1),
										ActualLoopCount:       null.IntFrom(6553600),
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "single join(SELECT name FROM users INNER JOIN todos ON ... WHERE todos.description = 'a')",
			explainResult: `
-> Nested loop inner join  (cost=1360547.06 rows=636171) (actual time=0.115..3426.457 rows=6553600 loops=1)
    -> Filter: (todos.` + "`" + `description` + "`" + ` = 'a')  (cost=662245.22 rows=636171) (actual time=0.046..2095.384 rows=6553600 loops=1)
        -> Table scan on todos  (cost=662245.22 rows=6361711) (actual time=0.043..1191.708 rows=6553600 loops=1)
    -> Single-row index lookup on users using PRIMARY (id=todos.user_id)  (cost=1.00 rows=1) (actual time=0.000..0.000 rows=1 loops=6553600)
`,
			expectedTree: &model.ExplainAnalyzeTree{
				Root: &model.ExplainAnalyzeTreeNode{
					AnalyzeResultLine: &model.ExplainAnalyzeResultLine{},
					Children: []*model.ExplainAnalyzeTreeNode{
						{
							AnalyzeResultLine: &model.ExplainAnalyzeResultLine{
								Text:                  "-> Nested loop inner join  (cost=1360547.06 rows=636171) (actual time=0.115..3426.457 rows=6553600 loops=1)",
								TableName:             "",
								EstimatedInitCost:     null.FloatFromPtr(nil),
								EstimatedCost:         null.FloatFrom(1360547.06),
								EstimatedReturnedRows: null.IntFrom(636171),
								ActualTimeFirstRow:    null.FloatFrom(0.115),
								ActualTimeAvg:         null.FloatFrom(3426.457),
								ActualReturnedRows:    null.IntFrom(6553600),
								ActualLoopCount:       null.IntFrom(1),
							},
							Children: []*model.ExplainAnalyzeTreeNode{
								{
									AnalyzeResultLine: &model.ExplainAnalyzeResultLine{
										Text:                  "    -> Filter: (todos.`description` = 'a')  (cost=662245.22 rows=636171) (actual time=0.046..2095.384 rows=6553600 loops=1)",
										TableName:             "todos",
										EstimatedInitCost:     null.FloatFromPtr(nil),
										EstimatedCost:         null.FloatFrom(662245.22),
										EstimatedReturnedRows: null.IntFrom(636171),
										ActualTimeFirstRow:    null.FloatFrom(0.046),
										ActualTimeAvg:         null.FloatFrom(2095.384),
										ActualReturnedRows:    null.IntFrom(6553600),
										ActualLoopCount:       null.IntFrom(1),
									},
									Children: []*model.ExplainAnalyzeTreeNode{
										{
											AnalyzeResultLine: &model.ExplainAnalyzeResultLine{
												Text:                  "        -> Table scan on todos  (cost=662245.22 rows=6361711) (actual time=0.043..1191.708 rows=6553600 loops=1)",
												TableName:             "todos",
												EstimatedInitCost:     null.FloatFromPtr(nil),
												EstimatedCost:         null.FloatFrom(662245.22),
												EstimatedReturnedRows: null.IntFrom(6361711),
												ActualTimeFirstRow:    null.FloatFrom(0.043),
												ActualTimeAvg:         null.FloatFrom(1191.708),
												ActualReturnedRows:    null.IntFrom(6553600),
												ActualLoopCount:       null.IntFrom(1),
											},
										},
									},
								},
								{
									AnalyzeResultLine: &model.ExplainAnalyzeResultLine{
										Text:                  "    -> Single-row index lookup on users using PRIMARY (id=todos.user_id)  (cost=1.00 rows=1) (actual time=0.000..0.000 rows=1 loops=6553600)",
										TableName:             "users",
										EstimatedInitCost:     null.FloatFromPtr(nil),
										EstimatedCost:         null.FloatFrom(1.0),
										EstimatedReturnedRows: null.IntFrom(1),
										ActualTimeFirstRow:    null.FloatFrom(0.0),
										ActualTimeAvg:         null.FloatFrom(0.0),
										ActualReturnedRows:    null.IntFrom(1),
										ActualLoopCount:       null.IntFrom(6553600),
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "deep nest",
			explainResult: `
-> Nested loop inner join  (cost=1)
    -> Nested loop inner join  (cost=2)
        -> Nested loop inner join  (cost=3)
            -> Nested loop inner join  (cost=4)
                -> Nested loop inner join  (cost=5)
-> Nested loop inner join  (cost=6)
`,
			expectedTree: &model.ExplainAnalyzeTree{
				Root: &model.ExplainAnalyzeTreeNode{
					AnalyzeResultLine: &model.ExplainAnalyzeResultLine{},
					Children: []*model.ExplainAnalyzeTreeNode{
						{
							AnalyzeResultLine: &model.ExplainAnalyzeResultLine{
								Text:          "-> Nested loop inner join  (cost=1)",
								EstimatedCost: null.FloatFrom(1.0),
							},
							Children: []*model.ExplainAnalyzeTreeNode{{
								AnalyzeResultLine: &model.ExplainAnalyzeResultLine{
									Text:          "    -> Nested loop inner join  (cost=2)",
									EstimatedCost: null.FloatFrom(2.0),
								},
								Children: []*model.ExplainAnalyzeTreeNode{{
									AnalyzeResultLine: &model.ExplainAnalyzeResultLine{
										Text:          "        -> Nested loop inner join  (cost=3)",
										EstimatedCost: null.FloatFrom(3.0),
									},
									Children: []*model.ExplainAnalyzeTreeNode{{
										AnalyzeResultLine: &model.ExplainAnalyzeResultLine{
											Text:          "            -> Nested loop inner join  (cost=4)",
											EstimatedCost: null.FloatFrom(4.0),
										},
										Children: []*model.ExplainAnalyzeTreeNode{{
											AnalyzeResultLine: &model.ExplainAnalyzeResultLine{
												Text:          "                -> Nested loop inner join  (cost=5)",
												EstimatedCost: null.FloatFrom(5.0),
											},
										},
										},
									}},
								}},
							}},
						},
						{
							AnalyzeResultLine: &model.ExplainAnalyzeResultLine{
								Text:          "-> Nested loop inner join  (cost=6)",
								EstimatedCost: null.FloatFrom(6.0),
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree, err := collector.CollectExplainAnalyzeTree(tt.explainResult)
			assert.NoError(t, err)
			assert.NotNil(t, tree)
			assert.Equal(t, tt.expectedTree, tree)
		})
	}
}
