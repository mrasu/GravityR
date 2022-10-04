package collector_test

import (
	"github.com/mrasu/GravityR/database/postgres/model"
	"github.com/mrasu/GravityR/database/postgres/model/collector"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"
	"strings"
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
Seq Scan on users  (cost=0.00..1400.36 rows=65536 width=9) (actual time=0.007..4.327 rows=65536 loops=1)
  Buffers: shared hit=745
Planning Time: 0.031 ms
Execution Time: 6.023 ms
`,
			expectedTree: &model.ExplainAnalyzeTree{
				SummaryText: `Planning Time: 0.031 ms
Execution Time: 6.023 ms`,
				Root: &model.ExplainAnalyzeTreeNode{
					SpaceSize:         -1,
					AnalyzeResultNode: &model.ExplainAnalyzeResultNode{},
					Children: []*model.ExplainAnalyzeTreeNode{
						{
							SpaceSize: 0,
							AnalyzeResultNode: &model.ExplainAnalyzeResultNode{
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
		{
			name: "single join(SELECT name FROM users INNER JOIN tasks ON ...)",
			explainResult: `
Hash Join  (cost=2219.56..86842.45 rows=1351260 width=9) (actual time=19.873..1520.060 rows=6553600 loops=1)
  Hash Cond: (tasks.user_id = users.id)
  Buffers: shared hit=12235 read=56073 dirtied=58571 written=55345
  ->  Seq Scan on tasks  (cost=0.00..81075.60 rows=1351260 width=4) (actual time=0.011..551.193 rows=6553600 loops=1)
        Buffers: shared hit=12186 read=55377 dirtied=58571 written=55345
  ->  Hash  (cost=1400.36..1400.36 rows=65536 width=13) (actual time=19.813..19.815 rows=65536 loops=1)
        Buckets: 65536  Batches: 1  Memory Usage: 3581kB
        Buffers: shared hit=49 read=696
        ->  Seq Scan on users  (cost=0.00..1400.36 rows=65536 width=13) (actual time=0.006..9.294 rows=65536 loops=1)
              Buffers: shared hit=49 read=696
Planning Time: 0.119 ms
Execution Time: 1693.383 ms
`,
			expectedTree: &model.ExplainAnalyzeTree{
				SummaryText: `Planning Time: 0.119 ms
Execution Time: 1693.383 ms`,
				Root: &model.ExplainAnalyzeTreeNode{
					SpaceSize:         -1,
					AnalyzeResultNode: &model.ExplainAnalyzeResultNode{},
					Children: []*model.ExplainAnalyzeTreeNode{
						{
							SpaceSize: 0,
							AnalyzeResultNode: &model.ExplainAnalyzeResultNode{
								Lines: []string{
									"Hash Join  (cost=2219.56..86842.45 rows=1351260 width=9) (actual time=19.873..1520.060 rows=6553600 loops=1)",
									"  Hash Cond: (tasks.user_id = users.id)",
									"  Buffers: shared hit=12235 read=56073 dirtied=58571 written=55345",
								},
								TableName:             "",
								EstimatedInitCost:     null.FloatFrom(2219.56),
								EstimatedCost:         null.FloatFrom(86842.45),
								EstimatedReturnedRows: null.IntFrom(1351260),
								EstimatedWidth:        null.IntFrom(9),
								ActualTimeFirstRow:    null.FloatFrom(19.873),
								ActualTimeAvg:         null.FloatFrom(1520.060),
								ActualReturnedRows:    null.IntFrom(6553600),
								ActualLoopCount:       null.IntFrom(1),
							},
							Children: []*model.ExplainAnalyzeTreeNode{
								{
									SpaceSize: 2,
									AnalyzeResultNode: &model.ExplainAnalyzeResultNode{
										Lines: []string{
											"  ->  Seq Scan on tasks  (cost=0.00..81075.60 rows=1351260 width=4) (actual time=0.011..551.193 rows=6553600 loops=1)",
											"        Buffers: shared hit=12186 read=55377 dirtied=58571 written=55345",
										},
										TableName:             "tasks",
										EstimatedInitCost:     null.FloatFrom(0.00),
										EstimatedCost:         null.FloatFrom(81075.60),
										EstimatedReturnedRows: null.IntFrom(1351260),
										EstimatedWidth:        null.IntFrom(4),
										ActualTimeFirstRow:    null.FloatFrom(0.011),
										ActualTimeAvg:         null.FloatFrom(551.193),
										ActualReturnedRows:    null.IntFrom(6553600),
										ActualLoopCount:       null.IntFrom(1),
									},
								},
								{
									SpaceSize: 2,
									AnalyzeResultNode: &model.ExplainAnalyzeResultNode{
										Lines: []string{
											"  ->  Hash  (cost=1400.36..1400.36 rows=65536 width=13) (actual time=19.813..19.815 rows=65536 loops=1)",
											"        Buckets: 65536  Batches: 1  Memory Usage: 3581kB",
											"        Buffers: shared hit=49 read=696",
										},
										TableName:             "",
										EstimatedInitCost:     null.FloatFrom(1400.36),
										EstimatedCost:         null.FloatFrom(1400.36),
										EstimatedReturnedRows: null.IntFrom(65536),
										EstimatedWidth:        null.IntFrom(13),
										ActualTimeFirstRow:    null.FloatFrom(19.813),
										ActualTimeAvg:         null.FloatFrom(19.815),
										ActualReturnedRows:    null.IntFrom(65536),
										ActualLoopCount:       null.IntFrom(1),
									},
									Children: []*model.ExplainAnalyzeTreeNode{
										{
											SpaceSize: 8,
											AnalyzeResultNode: &model.ExplainAnalyzeResultNode{
												Lines: []string{
													"        ->  Seq Scan on users  (cost=0.00..1400.36 rows=65536 width=13) (actual time=0.006..9.294 rows=65536 loops=1)",
													"              Buffers: shared hit=49 read=696",
												},
												TableName:             "users",
												EstimatedInitCost:     null.FloatFrom(0.00),
												EstimatedCost:         null.FloatFrom(1400.36),
												EstimatedReturnedRows: null.IntFrom(65536),
												EstimatedWidth:        null.IntFrom(13),
												ActualTimeFirstRow:    null.FloatFrom(0.006),
												ActualTimeAvg:         null.FloatFrom(9.294),
												ActualReturnedRows:    null.IntFrom(65536),
												ActualLoopCount:       null.IntFrom(1),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "single join(SELECT name FROM users INNER JOIN tasks ON ... WHERE tasks.description = 'a')",
			explainResult: `
Hash Join  (cost=2219.56..168906.66 rows=6553586 width=9) (actual time=24.270..1442.547 rows=6553600 loops=1)
  Hash Cond: (tasks.user_id = users.id)
  Buffers: shared hit=13123 read=55185
  ->  Seq Scan on tasks  (cost=0.00..149482.83 rows=6553586 width=4) (actual time=5.663..630.542 rows=6553600 loops=1)
        Filter: ((description)::text = 'test description'::text)
        Buffers: shared hit=12378 read=55185
  ->  Hash  (cost=1400.36..1400.36 rows=65536 width=13) (actual time=18.456..18.458 rows=65536 loops=1)
        Buckets: 65536  Batches: 1  Memory Usage: 3581kB
        Buffers: shared hit=745
        ->  Seq Scan on users  (cost=0.00..1400.36 rows=65536 width=13) (actual time=0.014..8.005 rows=65536 loops=1)
              Buffers: shared hit=745
Planning:
  Buffers: shared hit=6
Planning Time: 0.142 ms
JIT:
  Functions: 12
  Options: Inlining false, Optimization false, Expressions true, Deforming true
  Timing: Generation 0.918 ms, Inlining 0.000 ms, Optimization 0.303 ms, Emission 5.203 ms, Total 6.424 ms
Execution Time: 1609.370 ms
`,
			expectedTree: &model.ExplainAnalyzeTree{
				SummaryText: `Planning:
  Buffers: shared hit=6
Planning Time: 0.142 ms
JIT:
  Functions: 12
  Options: Inlining false, Optimization false, Expressions true, Deforming true
  Timing: Generation 0.918 ms, Inlining 0.000 ms, Optimization 0.303 ms, Emission 5.203 ms, Total 6.424 ms
Execution Time: 1609.370 ms`,
				Root: &model.ExplainAnalyzeTreeNode{
					SpaceSize:         -1,
					AnalyzeResultNode: &model.ExplainAnalyzeResultNode{},
					Children: []*model.ExplainAnalyzeTreeNode{
						{
							SpaceSize: 0,
							AnalyzeResultNode: &model.ExplainAnalyzeResultNode{
								Lines: []string{
									"Hash Join  (cost=2219.56..168906.66 rows=6553586 width=9) (actual time=24.270..1442.547 rows=6553600 loops=1)",
									"  Hash Cond: (tasks.user_id = users.id)",
									"  Buffers: shared hit=13123 read=55185",
								},
								TableName:             "",
								EstimatedInitCost:     null.FloatFrom(2219.56),
								EstimatedCost:         null.FloatFrom(168906.66),
								EstimatedReturnedRows: null.IntFrom(6553586),
								EstimatedWidth:        null.IntFrom(9),
								ActualTimeFirstRow:    null.FloatFrom(24.270),
								ActualTimeAvg:         null.FloatFrom(1442.547),
								ActualReturnedRows:    null.IntFrom(6553600),
								ActualLoopCount:       null.IntFrom(1),
							},
							Children: []*model.ExplainAnalyzeTreeNode{
								{
									SpaceSize: 2,
									AnalyzeResultNode: &model.ExplainAnalyzeResultNode{
										Lines: []string{
											"  ->  Seq Scan on tasks  (cost=0.00..149482.83 rows=6553586 width=4) (actual time=5.663..630.542 rows=6553600 loops=1)",
											"        Filter: ((description)::text = 'test description'::text)",
											"        Buffers: shared hit=12378 read=55185",
										},
										TableName:             "tasks",
										EstimatedInitCost:     null.FloatFrom(0.00),
										EstimatedCost:         null.FloatFrom(149482.83),
										EstimatedReturnedRows: null.IntFrom(6553586),
										EstimatedWidth:        null.IntFrom(4),
										ActualTimeFirstRow:    null.FloatFrom(5.663),
										ActualTimeAvg:         null.FloatFrom(630.542),
										ActualReturnedRows:    null.IntFrom(6553600),
										ActualLoopCount:       null.IntFrom(1),
									},
								},
								{
									SpaceSize: 2,
									AnalyzeResultNode: &model.ExplainAnalyzeResultNode{
										Lines: []string{
											"  ->  Hash  (cost=1400.36..1400.36 rows=65536 width=13) (actual time=18.456..18.458 rows=65536 loops=1)",
											"        Buckets: 65536  Batches: 1  Memory Usage: 3581kB",
											"        Buffers: shared hit=745",
										},
										TableName:             "",
										EstimatedInitCost:     null.FloatFrom(1400.36),
										EstimatedCost:         null.FloatFrom(1400.36),
										EstimatedReturnedRows: null.IntFrom(65536),
										EstimatedWidth:        null.IntFrom(13),
										ActualTimeFirstRow:    null.FloatFrom(18.456),
										ActualTimeAvg:         null.FloatFrom(18.458),
										ActualReturnedRows:    null.IntFrom(65536),
										ActualLoopCount:       null.IntFrom(1),
									},
									Children: []*model.ExplainAnalyzeTreeNode{
										{
											SpaceSize: 8,
											AnalyzeResultNode: &model.ExplainAnalyzeResultNode{
												Lines: []string{
													"        ->  Seq Scan on users  (cost=0.00..1400.36 rows=65536 width=13) (actual time=0.014..8.005 rows=65536 loops=1)",
													"              Buffers: shared hit=745",
												},
												TableName:             "users",
												EstimatedInitCost:     null.FloatFrom(0.00),
												EstimatedCost:         null.FloatFrom(1400.36),
												EstimatedReturnedRows: null.IntFrom(65536),
												EstimatedWidth:        null.IntFrom(13),
												ActualTimeFirstRow:    null.FloatFrom(0.014),
												ActualTimeAvg:         null.FloatFrom(8.005),
												ActualReturnedRows:    null.IntFrom(65536),
												ActualLoopCount:       null.IntFrom(1),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "CTE explain",
			explainResult: `
Sort  (cost=107094.02..107094.02 rows=1 width=64) (actual time=749.870..755.470 rows=2 loops=1)
  Sort Key: (CASE WHEN (count_tbl.status = 1) THEN 'Pending'::text WHEN (count_tbl.status = 2) THEN 'Doing'::text WHEN (count_tbl.status = 3) THEN 'Done'::text ELSE NULL::text END) DESC
  Sort Method: quicksort  Memory: 25kB
  Buffers: shared hit=1529 read=68308
  CTE count_tbl
    ->  Finalize GroupAggregate  (cost=107093.08..107093.84 rows=3 width=10) (actual time=749.833..755.432 rows=3 loops=1)
          Group Key: tasks.status
          Buffers: shared hit=1526 read=68308
          ->  Gather Merge  (cost=107093.08..107093.78 rows=6 width=10) (actual time=749.826..755.424 rows=9 loops=1)
                Workers Planned: 2
                Workers Launched: 2
                Buffers: shared hit=1526 read=68308
                ->  Sort  (cost=106093.05..106093.06 rows=3 width=10) (actual time=733.773..733.775 rows=3 loops=3)
                      Sort Key: tasks.status
                      Sort Method: quicksort  Memory: 25kB
                      Buffers: shared hit=1526 read=68308
                      Worker 0:  Sort Method: quicksort  Memory: 25kB
                      Worker 1:  Sort Method: quicksort  Memory: 25kB
                      ->  Partial HashAggregate  (cost=106093.00..106093.03 rows=3 width=10) (actual time=733.758..733.760 rows=3 loops=3)
                            Group Key: tasks.status
                            Batches: 1  Memory Usage: 24kB
                            Buffers: shared hit=1512 read=68308
                            Worker 0:  Batches: 1  Memory Usage: 24kB
                            Worker 1:  Batches: 1  Memory Usage: 24kB
                            ->  Hash Join  (cost=2014.32..104052.38 rows=408124 width=2) (actual time=28.492..698.394 rows=325800 loops=3)
                                  Hash Cond: (tasks.user_id = users.id)
                                  Buffers: shared hit=1512 read=68308
                                  ->  Parallel Seq Scan on tasks  (cost=0.00..94869.61 rows=2730661 width=6) (actual time=0.262..489.041 rows=2184533 loops=3)
                                        Buffers: shared read=67563
                                  ->  Hash  (cost=1891.88..1891.88 rows=9795 width=4) (actual time=27.080..27.081 rows=9774 loops=3)
                                        Buckets: 16384  Batches: 1  Memory Usage: 472kB
                                        Buffers: shared hit=1490 read=745
                                        ->  Seq Scan on users  (cost=0.00..1891.88 rows=9795 width=4) (actual time=13.415..25.919 rows=9774 loops=3)
                                              Filter: (created_at > (now() - '2 mons'::interval))
                                              Rows Removed by Filter: 55762
                                              Buffers: shared hit=1490 read=745
  InitPlan 2 (returns $2)
    ->  Aggregate  (cost=0.07..0.08 rows=1 width=32) (actual time=0.005..0.006 rows=1 loops=1)
          ->  CTE Scan on count_tbl count_tbl_1  (cost=0.00..0.06 rows=3 width=8) (actual time=0.000..0.002 rows=3 loops=1)
  ->  CTE Scan on count_tbl  (cost=0.00..0.09 rows=1 width=64) (actual time=749.849..749.850 rows=2 loops=1)
        Filter: ((status > 0) AND (count > 0))
        Rows Removed by Filter: 1
        Buffers: shared hit=301 read=23013
Planning:
  Buffers: shared hit=203 read=12
Planning Time: 2.202 ms
JIT:
  Functions: 61
  Options: Inlining false, Optimization false, Expressions true, Deforming true
  Timing: Generation 3.239 ms, Inlining 0.000 ms, Optimization 4.500 ms, Emission 34.318 ms, Total 42.056 ms
Execution Time: 821.344 ms
`,
			expectedTree: &model.ExplainAnalyzeTree{
				SummaryText: `Planning:
  Buffers: shared hit=203 read=12
Planning Time: 2.202 ms
JIT:
  Functions: 61
  Options: Inlining false, Optimization false, Expressions true, Deforming true
  Timing: Generation 3.239 ms, Inlining 0.000 ms, Optimization 4.500 ms, Emission 34.318 ms, Total 42.056 ms
Execution Time: 821.344 ms`,
				Root: &model.ExplainAnalyzeTreeNode{
					AnalyzeResultNode: &model.ExplainAnalyzeResultNode{},
					SpaceSize:         -1,
					Children: []*model.ExplainAnalyzeTreeNode{
						{
							SpaceSize: 0,
							AnalyzeResultNode: &model.ExplainAnalyzeResultNode{
								Lines: []string{
									"Sort  (cost=107094.02..107094.02 rows=1 width=64) (actual time=749.870..755.470 rows=2 loops=1)",
									"  Sort Key: (CASE WHEN (count_tbl.status = 1) THEN 'Pending'::text WHEN (count_tbl.status = 2) THEN 'Doing'::text WHEN (count_tbl.status = 3) THEN 'Done'::text ELSE NULL::text END) DESC",
									"  Sort Method: quicksort  Memory: 25kB",
									"  Buffers: shared hit=1529 read=68308",
								},
								TableName:             "",
								EstimatedInitCost:     null.FloatFrom(107094.02),
								EstimatedCost:         null.FloatFrom(107094.02),
								EstimatedReturnedRows: null.IntFrom(1),
								EstimatedWidth:        null.IntFrom(64),
								ActualTimeFirstRow:    null.FloatFrom(749.870),
								ActualTimeAvg:         null.FloatFrom(755.470),
								ActualReturnedRows:    null.IntFrom(2),
								ActualLoopCount:       null.IntFrom(1),
							},
							Children: []*model.ExplainAnalyzeTreeNode{
								{
									SpaceSize: 2,
									AnalyzeResultNode: &model.ExplainAnalyzeResultNode{
										Lines: []string{
											"  CTE count_tbl",
											"    ->  Finalize GroupAggregate  (cost=107093.08..107093.84 rows=3 width=10) (actual time=749.833..755.432 rows=3 loops=1)",
											"          Group Key: tasks.status",
											"          Buffers: shared hit=1526 read=68308",
										},
										TableName:             "tasks",
										EstimatedInitCost:     null.FloatFrom(107093.08),
										EstimatedCost:         null.FloatFrom(107093.84),
										EstimatedReturnedRows: null.IntFrom(3),
										EstimatedWidth:        null.IntFrom(10),
										ActualTimeFirstRow:    null.FloatFrom(749.833),
										ActualTimeAvg:         null.FloatFrom(755.432),
										ActualReturnedRows:    null.IntFrom(3),
										ActualLoopCount:       null.IntFrom(1),
									},
									Children: []*model.ExplainAnalyzeTreeNode{
										{
											SpaceSize: 10,
											AnalyzeResultNode: &model.ExplainAnalyzeResultNode{
												Lines: []string{
													"          ->  Gather Merge  (cost=107093.08..107093.78 rows=6 width=10) (actual time=749.826..755.424 rows=9 loops=1)",
													"                Workers Planned: 2",
													"                Workers Launched: 2",
													"                Buffers: shared hit=1526 read=68308",
												},
												TableName:             "",
												EstimatedInitCost:     null.FloatFrom(107093.08),
												EstimatedCost:         null.FloatFrom(107093.78),
												EstimatedReturnedRows: null.IntFrom(6),
												EstimatedWidth:        null.IntFrom(10),
												ActualTimeFirstRow:    null.FloatFrom(749.826),
												ActualTimeAvg:         null.FloatFrom(755.424),
												ActualReturnedRows:    null.IntFrom(9),
												ActualLoopCount:       null.IntFrom(1),
											},
											Children: []*model.ExplainAnalyzeTreeNode{
												{
													SpaceSize: 16,
													AnalyzeResultNode: &model.ExplainAnalyzeResultNode{
														Lines: []string{
															"                ->  Sort  (cost=106093.05..106093.06 rows=3 width=10) (actual time=733.773..733.775 rows=3 loops=3)",
															"                      Sort Key: tasks.status",
															"                      Sort Method: quicksort  Memory: 25kB",
															"                      Buffers: shared hit=1526 read=68308",
															"                      Worker 0:  Sort Method: quicksort  Memory: 25kB",
															"                      Worker 1:  Sort Method: quicksort  Memory: 25kB",
														},
														TableName:             "tasks",
														EstimatedInitCost:     null.FloatFrom(106093.05),
														EstimatedCost:         null.FloatFrom(106093.06),
														EstimatedReturnedRows: null.IntFrom(3),
														EstimatedWidth:        null.IntFrom(10),
														ActualTimeFirstRow:    null.FloatFrom(733.773),
														ActualTimeAvg:         null.FloatFrom(733.775),
														ActualReturnedRows:    null.IntFrom(3),
														ActualLoopCount:       null.IntFrom(3),
													},
													Children: []*model.ExplainAnalyzeTreeNode{
														{
															SpaceSize: 22,
															AnalyzeResultNode: &model.ExplainAnalyzeResultNode{
																Lines: []string{
																	"                      ->  Partial HashAggregate  (cost=106093.00..106093.03 rows=3 width=10) (actual time=733.758..733.760 rows=3 loops=3)",
																	"                            Group Key: tasks.status",
																	"                            Batches: 1  Memory Usage: 24kB",
																	"                            Buffers: shared hit=1512 read=68308",
																	"                            Worker 0:  Batches: 1  Memory Usage: 24kB",
																	"                            Worker 1:  Batches: 1  Memory Usage: 24kB",
																},
																TableName:             "tasks",
																EstimatedInitCost:     null.FloatFrom(106093.00),
																EstimatedCost:         null.FloatFrom(106093.03),
																EstimatedReturnedRows: null.IntFrom(3),
																EstimatedWidth:        null.IntFrom(10),
																ActualTimeFirstRow:    null.FloatFrom(733.758),
																ActualTimeAvg:         null.FloatFrom(733.760),
																ActualReturnedRows:    null.IntFrom(3),
																ActualLoopCount:       null.IntFrom(3),
															},
															Children: []*model.ExplainAnalyzeTreeNode{
																{
																	SpaceSize: 28,
																	AnalyzeResultNode: &model.ExplainAnalyzeResultNode{
																		Lines: []string{
																			"                            ->  Hash Join  (cost=2014.32..104052.38 rows=408124 width=2) (actual time=28.492..698.394 rows=325800 loops=3)",
																			"                                  Hash Cond: (tasks.user_id = users.id)",
																			"                                  Buffers: shared hit=1512 read=68308",
																		},
																		TableName:             "",
																		EstimatedInitCost:     null.FloatFrom(2014.32),
																		EstimatedCost:         null.FloatFrom(104052.38),
																		EstimatedReturnedRows: null.IntFrom(408124),
																		EstimatedWidth:        null.IntFrom(2),
																		ActualTimeFirstRow:    null.FloatFrom(28.492),
																		ActualTimeAvg:         null.FloatFrom(698.394),
																		ActualReturnedRows:    null.IntFrom(325800),
																		ActualLoopCount:       null.IntFrom(3),
																	},
																	Children: []*model.ExplainAnalyzeTreeNode{
																		{
																			SpaceSize: 34,
																			AnalyzeResultNode: &model.ExplainAnalyzeResultNode{
																				Lines: []string{
																					"                                  ->  Parallel Seq Scan on tasks  (cost=0.00..94869.61 rows=2730661 width=6) (actual time=0.262..489.041 rows=2184533 loops=3)",
																					"                                        Buffers: shared read=67563",
																				},
																				TableName:             "tasks",
																				EstimatedInitCost:     null.FloatFrom(0.00),
																				EstimatedCost:         null.FloatFrom(94869.61),
																				EstimatedReturnedRows: null.IntFrom(2730661),
																				EstimatedWidth:        null.IntFrom(6),
																				ActualTimeFirstRow:    null.FloatFrom(0.262),
																				ActualTimeAvg:         null.FloatFrom(489.041),
																				ActualReturnedRows:    null.IntFrom(2184533),
																				ActualLoopCount:       null.IntFrom(3),
																			},
																		},
																		{
																			SpaceSize: 34,
																			AnalyzeResultNode: &model.ExplainAnalyzeResultNode{
																				Lines: []string{
																					"                                  ->  Hash  (cost=1891.88..1891.88 rows=9795 width=4) (actual time=27.080..27.081 rows=9774 loops=3)",
																					"                                        Buckets: 16384  Batches: 1  Memory Usage: 472kB",
																					"                                        Buffers: shared hit=1490 read=745",
																				},
																				TableName:             "",
																				EstimatedInitCost:     null.FloatFrom(1891.88),
																				EstimatedCost:         null.FloatFrom(1891.88),
																				EstimatedReturnedRows: null.IntFrom(9795),
																				EstimatedWidth:        null.IntFrom(4),
																				ActualTimeFirstRow:    null.FloatFrom(27.080),
																				ActualTimeAvg:         null.FloatFrom(27.081),
																				ActualReturnedRows:    null.IntFrom(9774),
																				ActualLoopCount:       null.IntFrom(3),
																			},
																			Children: []*model.ExplainAnalyzeTreeNode{
																				{
																					SpaceSize: 40,
																					AnalyzeResultNode: &model.ExplainAnalyzeResultNode{
																						Lines: []string{
																							"                                        ->  Seq Scan on users  (cost=0.00..1891.88 rows=9795 width=4) (actual time=13.415..25.919 rows=9774 loops=3)",
																							"                                              Filter: (created_at > (now() - '2 mons'::interval))",
																							"                                              Rows Removed by Filter: 55762",
																							"                                              Buffers: shared hit=1490 read=745",
																						},
																						TableName:             "users",
																						EstimatedInitCost:     null.FloatFrom(0.00),
																						EstimatedCost:         null.FloatFrom(1891.88),
																						EstimatedReturnedRows: null.IntFrom(9795),
																						EstimatedWidth:        null.IntFrom(4),
																						ActualTimeFirstRow:    null.FloatFrom(13.415),
																						ActualTimeAvg:         null.FloatFrom(25.919),
																						ActualReturnedRows:    null.IntFrom(9774),
																						ActualLoopCount:       null.IntFrom(3),
																					},
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
								{
									SpaceSize: 2,
									AnalyzeResultNode: &model.ExplainAnalyzeResultNode{
										Lines: []string{
											"  InitPlan 2 (returns $2)",
											"    ->  Aggregate  (cost=0.07..0.08 rows=1 width=32) (actual time=0.005..0.006 rows=1 loops=1)",
										},
										TableName:             "",
										EstimatedInitCost:     null.FloatFrom(0.07),
										EstimatedCost:         null.FloatFrom(0.08),
										EstimatedReturnedRows: null.IntFrom(1),
										EstimatedWidth:        null.IntFrom(32),
										ActualTimeFirstRow:    null.FloatFrom(0.005),
										ActualTimeAvg:         null.FloatFrom(0.006),
										ActualReturnedRows:    null.IntFrom(1),
										ActualLoopCount:       null.IntFrom(1),
									},
									Children: []*model.ExplainAnalyzeTreeNode{
										{
											SpaceSize: 10,
											AnalyzeResultNode: &model.ExplainAnalyzeResultNode{
												Lines: []string{
													"          ->  CTE Scan on count_tbl count_tbl_1  (cost=0.00..0.06 rows=3 width=8) (actual time=0.000..0.002 rows=3 loops=1)",
												},
												TableName:             "count_tbl",
												EstimatedInitCost:     null.FloatFrom(0.00),
												EstimatedCost:         null.FloatFrom(0.06),
												EstimatedReturnedRows: null.IntFrom(3),
												EstimatedWidth:        null.IntFrom(8),
												ActualTimeFirstRow:    null.FloatFrom(0.000),
												ActualTimeAvg:         null.FloatFrom(0.002),
												ActualReturnedRows:    null.IntFrom(3),
												ActualLoopCount:       null.IntFrom(1),
											},
										},
									},
								},
								{
									SpaceSize: 2,
									AnalyzeResultNode: &model.ExplainAnalyzeResultNode{
										Lines: []string{
											"  ->  CTE Scan on count_tbl  (cost=0.00..0.09 rows=1 width=64) (actual time=749.849..749.850 rows=2 loops=1)",
											"        Filter: ((status > 0) AND (count > 0))",
											"        Rows Removed by Filter: 1",
											"        Buffers: shared hit=301 read=23013",
										},
										TableName:             "count_tbl",
										EstimatedInitCost:     null.FloatFrom(0.00),
										EstimatedCost:         null.FloatFrom(0.09),
										EstimatedReturnedRows: null.IntFrom(1),
										EstimatedWidth:        null.IntFrom(64),
										ActualTimeFirstRow:    null.FloatFrom(749.849),
										ActualTimeAvg:         null.FloatFrom(749.850),
										ActualReturnedRows:    null.IntFrom(2),
										ActualLoopCount:       null.IntFrom(1),
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			explainLines := strings.Split(strings.Trim(tt.explainResult, "\n"), "\n")
			tree, err := collector.CollectExplainAnalyzeTree(explainLines)
			assert.NoError(t, err)
			assert.NotNil(t, tree)
			assert.Equal(t, tt.expectedTree, tree)
		})
	}
}

func TestCollectExplainAnalyzeTree_NoAnalyze(t *testing.T) {
	tests := []struct {
		name          string
		explainResult string
		expectedTree  *model.ExplainAnalyzeTree
	}{
		{
			name: "simple(Subquery explain)",
			explainResult: `
Aggregate  (cost=105447.85..105447.86 rows=1 width=32)
  ->  Nested Loop Left Join  (cost=2564.50..105446.35 rows=100 width=55)
        ->  Gather  (cost=2564.21..104612.35 rows=100 width=27)
              Workers Planned: 2
              ->  Hash Join  (cost=1564.21..103602.35 rows=42 width=27)
                    Hash Cond: (tasks.user_id = __be_0_users.id)
                    ->  Parallel Seq Scan on tasks  (cost=0.00..94869.67 rows=2730667 width=27)
                    ->  Hash  (cost=1564.20..1564.20 rows=1 width=4)
                          ->  Seq Scan on users __be_0_users  (cost=0.00..1564.20 rows=1 width=4)
                                Filter: ((email)::text = 'test1111@example.com'::text)
        ->  Subquery Scan on "_root.or.user.base"  (cost=0.29..8.33 rows=1 width=32)
              ->  Limit  (cost=0.29..8.31 rows=1 width=566)
                    ->  Index Scan using users_pkey on users  (cost=0.29..8.31 rows=1 width=566)
                          Index Cond: (id = tasks.user_id)
              SubPlan 1
                ->  Result  (cost=0.00..0.01 rows=1 width=32)
  SubPlan 2
    ->  Result  (cost=0.00..0.01 rows=1 width=32)
JIT:
  Functions: 30
  Options: Inlining false, Optimization false, Expressions true, Deforming true
`,
			expectedTree: &model.ExplainAnalyzeTree{
				SummaryText: `JIT:
  Functions: 30
  Options: Inlining false, Optimization false, Expressions true, Deforming true`,
				Root: &model.ExplainAnalyzeTreeNode{
					SpaceSize:         -1,
					AnalyzeResultNode: &model.ExplainAnalyzeResultNode{},
					Children: []*model.ExplainAnalyzeTreeNode{
						{
							SpaceSize: 0,
							AnalyzeResultNode: &model.ExplainAnalyzeResultNode{
								Lines: []string{
									"Aggregate  (cost=105447.85..105447.86 rows=1 width=32)",
								},
								TableName:             "",
								EstimatedInitCost:     null.FloatFrom(105447.85),
								EstimatedCost:         null.FloatFrom(105447.86),
								EstimatedReturnedRows: null.IntFrom(1),
								EstimatedWidth:        null.IntFrom(32),
								ActualTimeFirstRow:    null.FloatFromPtr(nil),
								ActualTimeAvg:         null.FloatFromPtr(nil),
								ActualReturnedRows:    null.IntFromPtr(nil),
								ActualLoopCount:       null.IntFromPtr(nil),
							},
							Children: []*model.ExplainAnalyzeTreeNode{
								{
									SpaceSize: 2,
									AnalyzeResultNode: &model.ExplainAnalyzeResultNode{
										Lines: []string{
											"  ->  Nested Loop Left Join  (cost=2564.50..105446.35 rows=100 width=55)",
										},
										TableName:             "",
										EstimatedInitCost:     null.FloatFrom(2564.50),
										EstimatedCost:         null.FloatFrom(105446.35),
										EstimatedReturnedRows: null.IntFrom(100),
										EstimatedWidth:        null.IntFrom(55),
										ActualTimeFirstRow:    null.FloatFromPtr(nil),
										ActualTimeAvg:         null.FloatFromPtr(nil),
										ActualReturnedRows:    null.IntFromPtr(nil),
										ActualLoopCount:       null.IntFromPtr(nil),
									},
									Children: []*model.ExplainAnalyzeTreeNode{
										{
											SpaceSize: 8,
											AnalyzeResultNode: &model.ExplainAnalyzeResultNode{
												Lines: []string{
													"        ->  Gather  (cost=2564.21..104612.35 rows=100 width=27)",
													"              Workers Planned: 2",
												},
												TableName:             "",
												EstimatedInitCost:     null.FloatFrom(2564.21),
												EstimatedCost:         null.FloatFrom(104612.35),
												EstimatedReturnedRows: null.IntFrom(100),
												EstimatedWidth:        null.IntFrom(27),
												ActualTimeFirstRow:    null.FloatFromPtr(nil),
												ActualTimeAvg:         null.FloatFromPtr(nil),
												ActualReturnedRows:    null.IntFromPtr(nil),
												ActualLoopCount:       null.IntFromPtr(nil),
											},
											Children: []*model.ExplainAnalyzeTreeNode{
												{
													SpaceSize: 14,
													AnalyzeResultNode: &model.ExplainAnalyzeResultNode{
														Lines: []string{
															"              ->  Hash Join  (cost=1564.21..103602.35 rows=42 width=27)",
															"                    Hash Cond: (tasks.user_id = __be_0_users.id)",
														},
														TableName:             "",
														EstimatedInitCost:     null.FloatFrom(1564.21),
														EstimatedCost:         null.FloatFrom(103602.35),
														EstimatedReturnedRows: null.IntFrom(42),
														EstimatedWidth:        null.IntFrom(27),
														ActualTimeFirstRow:    null.FloatFromPtr(nil),
														ActualTimeAvg:         null.FloatFromPtr(nil),
														ActualReturnedRows:    null.IntFromPtr(nil),
														ActualLoopCount:       null.IntFromPtr(nil),
													},
													Children: []*model.ExplainAnalyzeTreeNode{
														{
															SpaceSize: 20,
															AnalyzeResultNode: &model.ExplainAnalyzeResultNode{
																Lines: []string{
																	"                    ->  Parallel Seq Scan on tasks  (cost=0.00..94869.67 rows=2730667 width=27)",
																},
																TableName:             "tasks",
																EstimatedInitCost:     null.FloatFrom(0.00),
																EstimatedCost:         null.FloatFrom(94869.67),
																EstimatedReturnedRows: null.IntFrom(2730667),
																EstimatedWidth:        null.IntFrom(27),
																ActualTimeFirstRow:    null.FloatFromPtr(nil),
																ActualTimeAvg:         null.FloatFromPtr(nil),
																ActualReturnedRows:    null.IntFromPtr(nil),
																ActualLoopCount:       null.IntFromPtr(nil),
															},
														},
														{
															SpaceSize: 20,
															AnalyzeResultNode: &model.ExplainAnalyzeResultNode{
																Lines: []string{
																	"                    ->  Hash  (cost=1564.20..1564.20 rows=1 width=4)",
																},
																TableName:             "",
																EstimatedInitCost:     null.FloatFrom(1564.20),
																EstimatedCost:         null.FloatFrom(1564.20),
																EstimatedReturnedRows: null.IntFrom(1),
																EstimatedWidth:        null.IntFrom(4),
																ActualTimeFirstRow:    null.FloatFromPtr(nil),
																ActualTimeAvg:         null.FloatFromPtr(nil),
																ActualReturnedRows:    null.IntFromPtr(nil),
																ActualLoopCount:       null.IntFromPtr(nil),
															},
															Children: []*model.ExplainAnalyzeTreeNode{
																{
																	SpaceSize: 26,
																	AnalyzeResultNode: &model.ExplainAnalyzeResultNode{
																		Lines: []string{
																			"                          ->  Seq Scan on users __be_0_users  (cost=0.00..1564.20 rows=1 width=4)",
																			"                                Filter: ((email)::text = 'test1111@example.com'::text)",
																		},
																		TableName:             "users",
																		EstimatedInitCost:     null.FloatFrom(0.00),
																		EstimatedCost:         null.FloatFrom(1564.20),
																		EstimatedReturnedRows: null.IntFrom(1),
																		EstimatedWidth:        null.IntFrom(4),
																		ActualTimeFirstRow:    null.FloatFromPtr(nil),
																		ActualTimeAvg:         null.FloatFromPtr(nil),
																		ActualReturnedRows:    null.IntFromPtr(nil),
																		ActualLoopCount:       null.IntFromPtr(nil),
																	},
																},
															},
														},
													},
												},
											},
										},
										{
											SpaceSize: 8,
											AnalyzeResultNode: &model.ExplainAnalyzeResultNode{
												Lines: []string{
													"        ->  Subquery Scan on \"_root.or.user.base\"  (cost=0.29..8.33 rows=1 width=32)",
												},
												TableName:             "_root.or.user.base",
												EstimatedInitCost:     null.FloatFrom(0.29),
												EstimatedCost:         null.FloatFrom(8.33),
												EstimatedReturnedRows: null.IntFrom(1),
												EstimatedWidth:        null.IntFrom(32),
												ActualTimeFirstRow:    null.FloatFromPtr(nil),
												ActualTimeAvg:         null.FloatFromPtr(nil),
												ActualReturnedRows:    null.IntFromPtr(nil),
												ActualLoopCount:       null.IntFromPtr(nil),
											},
											Children: []*model.ExplainAnalyzeTreeNode{
												{
													SpaceSize: 14,
													AnalyzeResultNode: &model.ExplainAnalyzeResultNode{
														Lines: []string{
															"              ->  Limit  (cost=0.29..8.31 rows=1 width=566)",
														},
														TableName:             "",
														EstimatedInitCost:     null.FloatFrom(0.29),
														EstimatedCost:         null.FloatFrom(8.31),
														EstimatedReturnedRows: null.IntFrom(1),
														EstimatedWidth:        null.IntFrom(566),
														ActualTimeFirstRow:    null.FloatFromPtr(nil),
														ActualTimeAvg:         null.FloatFromPtr(nil),
														ActualReturnedRows:    null.IntFromPtr(nil),
														ActualLoopCount:       null.IntFromPtr(nil),
													},
													Children: []*model.ExplainAnalyzeTreeNode{
														{
															SpaceSize: 20,
															AnalyzeResultNode: &model.ExplainAnalyzeResultNode{
																Lines: []string{
																	"                    ->  Index Scan using users_pkey on users  (cost=0.29..8.31 rows=1 width=566)",
																	"                          Index Cond: (id = tasks.user_id)",
																},
																TableName:             "users",
																EstimatedInitCost:     null.FloatFrom(0.29),
																EstimatedCost:         null.FloatFrom(8.31),
																EstimatedReturnedRows: null.IntFrom(1),
																EstimatedWidth:        null.IntFrom(566),
																ActualTimeFirstRow:    null.FloatFromPtr(nil),
																ActualTimeAvg:         null.FloatFromPtr(nil),
																ActualReturnedRows:    null.IntFromPtr(nil),
																ActualLoopCount:       null.IntFromPtr(nil),
															},
														},
													},
												},
												{
													SpaceSize: 14,
													AnalyzeResultNode: &model.ExplainAnalyzeResultNode{
														Lines: []string{
															"              SubPlan 1",
															"                ->  Result  (cost=0.00..0.01 rows=1 width=32)",
														},
														TableName:             "",
														EstimatedInitCost:     null.FloatFrom(0.00),
														EstimatedCost:         null.FloatFrom(0.01),
														EstimatedReturnedRows: null.IntFrom(1),
														EstimatedWidth:        null.IntFrom(32),
														ActualTimeFirstRow:    null.FloatFromPtr(nil),
														ActualTimeAvg:         null.FloatFromPtr(nil),
														ActualReturnedRows:    null.IntFromPtr(nil),
														ActualLoopCount:       null.IntFromPtr(nil),
													},
												},
											},
										},
									},
								},
								{
									SpaceSize: 2,
									AnalyzeResultNode: &model.ExplainAnalyzeResultNode{
										Lines: []string{
											"  SubPlan 2",
											"    ->  Result  (cost=0.00..0.01 rows=1 width=32)",
										},
										TableName:             "",
										EstimatedInitCost:     null.FloatFrom(0.00),
										EstimatedCost:         null.FloatFrom(0.01),
										EstimatedReturnedRows: null.IntFrom(1),
										EstimatedWidth:        null.IntFrom(32),
										ActualTimeFirstRow:    null.FloatFromPtr(nil),
										ActualTimeAvg:         null.FloatFromPtr(nil),
										ActualReturnedRows:    null.IntFromPtr(nil),
										ActualLoopCount:       null.IntFromPtr(nil),
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			explainLines := strings.Split(strings.Trim(tt.explainResult, "\n"), "\n")
			tree, err := collector.CollectExplainAnalyzeTree(explainLines)
			assert.NoError(t, err)
			assert.NotNil(t, tree)
			assert.Equal(t, tt.expectedTree, tree)
		})
	}
}
