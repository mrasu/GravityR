package pservice_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
	"github.com/mrasu/GravityR/database/postgres/pmodel"
	"github.com/mrasu/GravityR/database/postgres/pservice"
	"github.com/mrasu/GravityR/infra/postgres"
	"github.com/mrasu/GravityR/thelper"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
	"strings"
	"testing"
)

func TestExplainer_ExplainWithAnalyze(t *testing.T) {
	tests := []struct {
		name          string
		explainResult string
		expectedTree  *pmodel.ExplainAnalyzeTree
	}{
		{
			name: "simple(SELECT name FROM users)",
			explainResult: `
Seq Scan on users  (cost=0.00..1400.36 rows=65536 width=9) (actual time=0.007..4.327 rows=65536 loops=1)
  Buffers: shared hit=745
Planning Time: 0.031 ms
Execution Time: 6.023 ms
`,
			expectedTree: &pmodel.ExplainAnalyzeTree{
				SummaryText: `Planning Time: 0.031 ms
Execution Time: 6.023 ms`,
				Root: &pmodel.ExplainAnalyzeTreeNode{
					SpaceSize:         -1,
					AnalyzeResultNode: &pmodel.ExplainAnalyzeResultNode{},
					Children: []*pmodel.ExplainAnalyzeTreeNode{
						{
							SpaceSize: 0,
							AnalyzeResultNode: &pmodel.ExplainAnalyzeResultNode{
								Lines: []string{
									"Seq Scan on users  (cost=0.00..1400.36 rows=65536 width=9) (actual time=0.007..4.327 rows=65536 loops=1)",
									"  Buffers: shared hit=745",
								},
								TableName:             "users",
								EstimatedInitCost:     null.FloatFrom(0.00),
								EstimatedCost:         null.FloatFrom(1400.36),
								EstimatedReturnedRows: null.IntFrom(65536),
								EstimatedWidth:        null.IntFrom(9),
								ActualTimeFirstRow:    null.FloatFrom(0.007),
								ActualTimeAvg:         null.FloatFrom(4.327),
								ActualReturnedRows:    null.IntFrom(65536),
								ActualLoopCount:       null.IntFrom(1),
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			thelper.MockPostgresDB(t, func(db *postgres.DB, mock sqlmock.Sqlmock) {
				mockExplainAnalyze(tt.explainResult, mock)

				tree, err := pservice.NewExplainer(db).ExplainWithAnalyze("SELECT ...")
				require.NoError(t, err)

				if diff := cmp.Diff(tt.expectedTree, tree); diff != "" {
					t.Errorf(diff)
				}
			})
		})
	}
}

func mockExplainAnalyze(explainResult string, mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{"QUERY PLAN"})
	explainLines := strings.Split(strings.Trim(explainResult, "\n"), "\n")
	for _, ex := range explainLines {
		rows.AddRow(ex)
	}
	mock.ExpectQuery("EXPLAIN \\(ANALYZE, BUFFERS\\) SELECT").WillReturnRows(rows)
}
