package pmodel_test

import (
	"github.com/mrasu/GravityR/database/postgres/pmodel"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"
	"testing"
)

func TestParseExplainAnalyzeResultLineNode(t *testing.T) {
	tests := []struct {
		name         string
		lines        []string
		expectedNest int
		expectedLine *pmodel.ExplainAnalyzeResultNode
	}{
		{
			name: "simple",
			lines: []string{
				"Seq Scan on users  (cost=0.00..1400.36 rows=65536 width=9) (actual time=0.007..4.327 rows=65536 loops=1)",
				"  Buffers: shared hit=745",
			},
			expectedNest: 0,
			expectedLine: &pmodel.ExplainAnalyzeResultNode{
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
		{
			name: "index range scan(where id in (1,2,3))",
			lines: []string{
				"Index Scan using tasks_pkey on tasks  (cost=0.43..17.35 rows=3 width=54) (actual time=0.039..0.042 rows=3 loops=1)",
				"  Index Cond: (id = ANY ('{1,2,3}'::integer[]))",
				"  Buffers: shared hit=9 read=4",
			},
			expectedNest: 0,
			expectedLine: &pmodel.ExplainAnalyzeResultNode{
				Lines: []string{
					"Index Scan using tasks_pkey on tasks  (cost=0.43..17.35 rows=3 width=54) (actual time=0.039..0.042 rows=3 loops=1)",
					"  Index Cond: (id = ANY ('{1,2,3}'::integer[]))",
					"  Buffers: shared hit=9 read=4",
				},
				TableName:             "tasks",
				EstimatedInitCost:     null.FloatFrom(0.43),
				EstimatedCost:         null.FloatFrom(17.35),
				EstimatedReturnedRows: null.IntFrom(3),
				EstimatedWidth:        null.IntFrom(54),
				ActualTimeFirstRow:    null.FloatFrom(0.039),
				ActualTimeAvg:         null.FloatFrom(0.042),
				ActualReturnedRows:    null.IntFrom(3),
				ActualLoopCount:       null.IntFrom(1),
			},
		},
		{
			name: "no table",
			lines: []string{
				"->  Hash Join  (cost=2219.56..3791.96 rows=65536 width=17) (actual time=14.154..34.793 rows=65536 loops=1)",
				"  Hash Cond: (users.id = u2.id)",
				"  Buffers: shared hit=1490",
			},
			expectedNest: 0,
			expectedLine: &pmodel.ExplainAnalyzeResultNode{
				Lines: []string{
					"->  Hash Join  (cost=2219.56..3791.96 rows=65536 width=17) (actual time=14.154..34.793 rows=65536 loops=1)",
					"  Hash Cond: (users.id = u2.id)",
					"  Buffers: shared hit=1490",
				},
				TableName:             "",
				EstimatedInitCost:     null.FloatFrom(2219.56),
				EstimatedCost:         null.FloatFrom(3791.96),
				EstimatedReturnedRows: null.IntFrom(65536),
				EstimatedWidth:        null.IntFrom(17),
				ActualTimeFirstRow:    null.FloatFrom(14.154),
				ActualTimeAvg:         null.FloatFrom(34.793),
				ActualReturnedRows:    null.IntFrom(65536),
				ActualLoopCount:       null.IntFrom(1),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nest, parsedLine, err := pmodel.ParseExplainAnalyzeResultLineNode(tt.lines)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedNest, nest)
			assert.Equal(t, tt.expectedLine, parsedLine)
		})
	}
}

func TestExplainAnalyzeResultLineNode_EstimatedTotalTime(t *testing.T) {
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
			expected: 3.0,
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
			expected: 3.0,
		},
		{
			name:     "massive loop but 0 avg",
			avg:      null.FloatFrom(0.0),
			loop:     null.IntFrom(1000000),
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := pmodel.ExplainAnalyzeResultNode{ActualTimeAvg: tt.avg, ActualLoopCount: tt.loop}
			assert.Equal(t, tt.expected, l.EstimatedTotalTime())
		})
	}
}

func TestExplainAnalyzeResultLineNode_Title(t *testing.T) {
	tests := []struct {
		name     string
		lines    []string
		expected string
	}{
		{
			name: "with `->`",
			lines: []string{
				"  ->  Hash Join  (cost=4439.12..6183.57 rows=65536 width=21) (actual time=28.373..64.680 rows=65536 loops=1)",
				"        Hash Cond: (users.id = u3.id)",
				"        Buffers: shared hit=2235",
			},
			expected: "Hash Join",
		},
		{
			name: "no `->`",
			lines: []string{
				"Hash Join  (cost=6658.68..8575.17 rows=65536 width=9) (actual time=42.307..93.050 rows=65536 loops=1)",
				"  Hash Cond: (users.id = u4.id)",
				"  Buffers: shared hit=2980",
			},
			expected: "Hash Join",
		},
		{
			name: "with `on`",
			lines: []string{
				"Seq Scan on users  (cost=0.00..1400.36 rows=65536 width=9) (actual time=0.007..4.327 rows=65536 loops=1)",
				"  Buffers: shared hit=745",
			},
			expected: "Seq Scan",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := pmodel.ExplainAnalyzeResultNode{Lines: tt.lines}
			assert.Equal(t, tt.expected, l.Title())
		})
	}
}
