package mmodel

import (
	"fmt"
)

// AnalyzeResultLine is the parsed information from `EXPLAIN ANALYZE` query
// c.f) https://dev.mysql.com/doc/refman/8.0/en/explain.html#explain-analyze
type AnalyzeResultLine struct {
	Text      string
	TableName string

	EstimatedInitCost     float64
	EstimatedCost         float64
	EstimatedReturnedRows int

	ActualTimeFirstRow float64
	ActualTimeAvg      float64
	ActualReturnedRows int
	ActualLoopCount    int
}

func (l *AnalyzeResultLine) ActualTotalTime() float64 {
	return l.ActualTimeAvg * float64(l.ActualLoopCount)
}

func (l *AnalyzeResultLine) String() string {
	return fmt.Sprintf(
		"[%s] estimate_cost=%f, estimate_rows=%d, time=%f..%f, totalTime=%f, rows=%d, loops=%d\n%s",
		l.TableName,
		l.EstimatedCost, l.EstimatedReturnedRows,
		l.ActualTimeFirstRow, l.ActualTimeAvg, l.ActualTotalTime(),
		l.ActualReturnedRows, l.ActualLoopCount,
		l.Text,
	)
}
