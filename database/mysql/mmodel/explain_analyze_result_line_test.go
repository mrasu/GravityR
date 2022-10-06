package mmodel

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"
	"testing"
)

func TestParseExplainAnalyzeResultLine(t *testing.T) {
	tests := []struct {
		name         string
		line         string
		expectedNest int
		expectedLine *ExplainAnalyzeResultLine
	}{
		{
			name:         "simple",
			line:         "-> Table scan on tasks  (cost=659696.88 rows=6361711) (actual time=0.069..1803.459 rows=6553600 loops=1)",
			expectedNest: 0,
			expectedLine: &ExplainAnalyzeResultLine{
				Text:                  "-> Table scan on tasks  (cost=659696.88 rows=6361711) (actual time=0.069..1803.459 rows=6553600 loops=1)",
				TableName:             "tasks",
				EstimatedInitCost:     null.FloatFromPtr(nil),
				EstimatedCost:         null.FloatFrom(659696.88),
				EstimatedReturnedRows: null.IntFrom(6361711),
				ActualTimeFirstRow:    null.FloatFrom(0.069),
				ActualTimeAvg:         null.FloatFrom(1803.459),
				ActualReturnedRows:    null.IntFrom(6553600),
				ActualLoopCount:       null.IntFrom(1),
			},
		},
		{
			name:         "index range scan(where id in (1,2,3))",
			line:         "-> Index range scan on users using PRIMARY  (cost=2.19 rows=3) (actual time=0.026..0.029 rows=3 loops=1)",
			expectedNest: 0,
			expectedLine: &ExplainAnalyzeResultLine{
				Text:                  "-> Index range scan on users using PRIMARY  (cost=2.19 rows=3) (actual time=0.026..0.029 rows=3 loops=1)",
				TableName:             "users",
				EstimatedInitCost:     null.FloatFromPtr(nil),
				EstimatedCost:         null.FloatFrom(2.19),
				EstimatedReturnedRows: null.IntFrom(3),
				ActualTimeFirstRow:    null.FloatFrom(0.026),
				ActualTimeAvg:         null.FloatFrom(0.029),
				ActualReturnedRows:    null.IntFrom(3),
				ActualLoopCount:       null.IntFrom(1),
			},
		},
		{
			name:         "no table",
			line:         "-> Nested loop inner join  (cost=1243547.38 rows=636171) (actual time=0.122..2549.635 rows=2228224 loops=1)",
			expectedNest: 0,
			expectedLine: &ExplainAnalyzeResultLine{
				Text:                  "-> Nested loop inner join  (cost=1243547.38 rows=636171) (actual time=0.122..2549.635 rows=2228224 loops=1)",
				TableName:             "",
				EstimatedInitCost:     null.FloatFromPtr(nil),
				EstimatedCost:         null.FloatFrom(1243547.38),
				EstimatedReturnedRows: null.IntFrom(636171),
				ActualTimeFirstRow:    null.FloatFrom(0.122),
				ActualTimeAvg:         null.FloatFrom(2549.635),
				ActualReturnedRows:    null.IntFrom(2228224),
				ActualLoopCount:       null.IntFrom(1),
			},
		},
		{
			name:         "never executed",
			line:         "-> Single-row index lookup on users using PRIMARY (id=tasks.user_id)  (cost=0.25 rows=1) (never executed)",
			expectedNest: 0,
			expectedLine: &ExplainAnalyzeResultLine{
				Text:                  "-> Single-row index lookup on users using PRIMARY (id=tasks.user_id)  (cost=0.25 rows=1) (never executed)",
				TableName:             "users",
				EstimatedInitCost:     null.FloatFromPtr(nil),
				EstimatedCost:         null.FloatFrom(0.25),
				EstimatedReturnedRows: null.IntFrom(1),
				ActualTimeFirstRow:    null.FloatFromPtr(nil),
				ActualTimeAvg:         null.FloatFromPtr(nil),
				ActualReturnedRows:    null.IntFromPtr(nil),
				ActualLoopCount:       null.IntFromPtr(nil),
			},
		},
		{
			name:         "filtered by index(where id = 1)",
			line:         "-> Rows fetched before execution  (cost=0.00 rows=1) (actual time=0.000..0.000 rows=1 loops=1)",
			expectedNest: 0,
			expectedLine: &ExplainAnalyzeResultLine{
				Text:                  "-> Rows fetched before execution  (cost=0.00 rows=1) (actual time=0.000..0.000 rows=1 loops=1)",
				TableName:             "",
				EstimatedInitCost:     null.FloatFromPtr(nil),
				EstimatedCost:         null.FloatFrom(0.0),
				EstimatedReturnedRows: null.IntFrom(1),
				ActualTimeFirstRow:    null.FloatFrom(0.0),
				ActualTimeAvg:         null.FloatFrom(0.0),
				ActualReturnedRows:    null.IntFrom(1),
				ActualLoopCount:       null.IntFrom(1),
			},
		},
		{
			name:         "filtered by index but no row matches",
			line:         "-> Zero rows (no matching row in const table)  (actual time=0.000..0.000 rows=0 loops=1)",
			expectedNest: 0,
			expectedLine: &ExplainAnalyzeResultLine{
				Text:                  "-> Zero rows (no matching row in const table)  (actual time=0.000..0.000 rows=0 loops=1)",
				TableName:             "",
				EstimatedInitCost:     null.FloatFromPtr(nil),
				EstimatedCost:         null.FloatFromPtr(nil),
				EstimatedReturnedRows: null.IntFromPtr(nil),
				ActualTimeFirstRow:    null.FloatFrom(0.0),
				ActualTimeAvg:         null.FloatFrom(0.0),
				ActualReturnedRows:    null.IntFrom(0),
				ActualLoopCount:       null.IntFrom(1),
			},
		},
		{
			name:         "average cost exists",
			line:         "-> Table scan on tasks  (cost=1..2 rows=3) (actual time=0.069..1803.459 rows=6553600 loops=1)",
			expectedNest: 0,
			expectedLine: &ExplainAnalyzeResultLine{
				Text:                  "-> Table scan on tasks  (cost=1..2 rows=3) (actual time=0.069..1803.459 rows=6553600 loops=1)",
				TableName:             "tasks",
				EstimatedInitCost:     null.FloatFrom(1),
				EstimatedCost:         null.FloatFrom(2.0),
				EstimatedReturnedRows: null.IntFrom(3),
				ActualTimeFirstRow:    null.FloatFrom(0.069),
				ActualTimeAvg:         null.FloatFrom(1803.459),
				ActualReturnedRows:    null.IntFrom(6553600),
				ActualLoopCount:       null.IntFrom(1),
			},
		},
		{
			name:         "filter",
			line:         "-> Filter: (tasks.`status` = 1)  (cost=659696.88 rows=636171) (actual time=0.071..2047.548 rows=2228224 loops=1)",
			expectedNest: 0,
			expectedLine: &ExplainAnalyzeResultLine{
				Text:                  "-> Filter: (tasks.`status` = 1)  (cost=659696.88 rows=636171) (actual time=0.071..2047.548 rows=2228224 loops=1)",
				TableName:             "tasks",
				EstimatedInitCost:     null.FloatFromPtr(nil),
				EstimatedCost:         null.FloatFrom(659696.88),
				EstimatedReturnedRows: null.IntFrom(636171),
				ActualTimeFirstRow:    null.FloatFrom(0.071),
				ActualTimeAvg:         null.FloatFrom(2047.548),
				ActualReturnedRows:    null.IntFrom(2228224),
				ActualLoopCount:       null.IntFrom(1),
			},
		},
		{
			name:         "sort",
			line:         "-> Sort: users.email DESC  (cost=6471.25 rows=63830) (actual time=12.631..12.631 rows=0 loops=1)",
			expectedNest: 0,
			expectedLine: &ExplainAnalyzeResultLine{
				Text:                  "-> Sort: users.email DESC  (cost=6471.25 rows=63830) (actual time=12.631..12.631 rows=0 loops=1)",
				TableName:             "users",
				EstimatedInitCost:     null.FloatFromPtr(nil),
				EstimatedCost:         null.FloatFrom(6471.25),
				EstimatedReturnedRows: null.IntFrom(63830),
				ActualTimeFirstRow:    null.FloatFrom(12.631),
				ActualTimeAvg:         null.FloatFrom(12.631),
				ActualReturnedRows:    null.IntFrom(0),
				ActualLoopCount:       null.IntFrom(1),
			},
		},
		{
			name:         "aggregate",
			line:         "-> Aggregate: count(users.email)  (actual time=16.468..16.468 rows=1 loops=1)",
			expectedNest: 0,
			expectedLine: &ExplainAnalyzeResultLine{
				Text:                  "-> Aggregate: count(users.email)  (actual time=16.468..16.468 rows=1 loops=1)",
				TableName:             "users",
				EstimatedInitCost:     null.FloatFromPtr(nil),
				EstimatedCost:         null.FloatFromPtr(nil),
				EstimatedReturnedRows: null.IntFromPtr(nil),
				ActualTimeFirstRow:    null.FloatFrom(16.468),
				ActualTimeAvg:         null.FloatFrom(16.468),
				ActualReturnedRows:    null.IntFrom(1),
				ActualLoopCount:       null.IntFrom(1),
			},
		},
		{
			name:         "aggregate count(*)",
			line:         "-> Aggregate: count(0)  (actual time=10.124..10.124 rows=1 loops=1)",
			expectedNest: 0,
			expectedLine: &ExplainAnalyzeResultLine{
				Text:                  "-> Aggregate: count(0)  (actual time=10.124..10.124 rows=1 loops=1)",
				TableName:             "",
				EstimatedInitCost:     null.FloatFromPtr(nil),
				EstimatedCost:         null.FloatFromPtr(nil),
				EstimatedReturnedRows: null.IntFromPtr(nil),
				ActualTimeFirstRow:    null.FloatFrom(10.124),
				ActualTimeAvg:         null.FloatFrom(10.124),
				ActualReturnedRows:    null.IntFrom(1),
				ActualLoopCount:       null.IntFrom(1),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nest, parsedLine, err := ParseExplainAnalyzeResultLine(tt.line)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedNest, nest)
			assert.Equal(t, tt.expectedLine, parsedLine)
		})
	}
}

func TestExplainAnalyzeResultLine_EstimatedTotalTime(t *testing.T) {
	tests := []struct {
		name     string
		avg      null.Float
		loop     null.Int
		expected float64
	}{
		{
			name:     "simple",
			avg:      null.FloatFrom(3.0),
			loop:     null.IntFrom(2),
			expected: 6.0,
		},
		{
			name:     "single loop",
			avg:      null.FloatFrom(3.0),
			loop:     null.IntFrom(1),
			expected: 3.0,
		},
		{
			name:     "massive loop",
			avg:      null.FloatFrom(3.0),
			loop:     null.IntFrom(1000000),
			expected: 3000000.0,
		},
		{
			name: "massive loop but 0 avg",
			avg:  null.FloatFrom(0.0),
			loop: null.IntFrom(1000000),
			// 1000000 * 0.0001
			expected: 100.0,
		},
		{
			name:     "nil",
			avg:      null.FloatFromPtr(nil),
			loop:     null.IntFromPtr(nil),
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := ExplainAnalyzeResultLine{ActualTimeAvg: tt.avg, ActualLoopCount: tt.loop}
			assert.Equal(t, tt.expected, l.EstimatedTotalTime())
		})
	}
}

func TestExplainAnalyzeResultLine_Title(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected string
	}{
		{
			name:     "with colon",
			text:     "-> Aggregate: count(0)  (actual time=10.124..10.124 rows=1 loops=1)",
			expected: "Aggregate",
		},
		{
			name:     "with leading space",
			text:     "   -> Aggregate: count(0)  (actual time=10.124..10.124 rows=1 loops=1)",
			expected: "Aggregate",
		},
		{
			name:     "with `on`",
			text:     "-> Table scan on count_tbl  (cost=238542.78 rows=2120358) (actual time=0.000..0.000 rows=3 loops=1)",
			expected: "Table scan",
		},
		{
			name:     "with `on`2",
			text:     "-> Single-row index lookup on users using PRIMARY (user_id=tasks.user_id)  (cost=0.96 rows=1) (actual time=0.000..0.000 rows=1 loops=6553600)",
			expected: "Single-row index lookup",
		},
		{
			name:     "with `join`",
			text:     "-> Nested loop inner join  (cost=7376073.37 rows=2120358) (actual time=0.107..3196.300 rows=1086000 loops=1)",
			expected: "Nested loop inner join",
		},
		{
			name:     "no marker",
			text:     "-> Materialize CTE count_tbl if needed  (actual time=3371.802..3371.802 rows=3 loops=1)",
			expected: "Materialize",
		},
		{
			name:     "no marker2",
			text:     "-> Select #2 (subquery in projection; run only once)",
			expected: "Select",
		},
		{
			name:     "no arrow",
			text:     "Select #2 (subquery in projection; run only once)",
			expected: "Select",
		},
		{
			name:     "unknown",
			text:     "123 456 789",
			expected: "123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := ExplainAnalyzeResultLine{Text: tt.text}
			assert.Equal(t, tt.expected, l.Title())
		})
	}
}
