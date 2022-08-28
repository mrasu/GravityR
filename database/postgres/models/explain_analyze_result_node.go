package models

import (
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v4"
	"regexp"
	"strconv"
)

// ExplainAnalyzeResultNode is the parsed information from `EXPLAIN (ANALYZE, BUFFERS)` query
// c.f)
//	https://www.postgresql.org/docs/14/using-explain.html
//	https://github.com/postgres/postgres/blob/REL_14_5/src/backend/commands/explain.c#L1144
type ExplainAnalyzeResultNode struct {
	Lines     []string
	TableName string

	EstimatedInitCost     null.Float
	EstimatedCost         null.Float
	EstimatedReturnedRows null.Int
	EstimatedWidth        null.Int

	ActualTimeFirstRow null.Float
	ActualTimeAvg      null.Float
	ActualReturnedRows null.Int
	ActualLoopCount    null.Int
}

var prefixSpaceReg = regexp.MustCompile(`^\s*`)
var onReg = regexp.MustCompile(` on (\w+)`)
var keyReg = regexp.MustCompile(` Key: (\w+)\.`)
var metricLineReg = regexp.MustCompile(`\(cost=.+\) (\(actual .+\))`)
var analyzeEstimationReg = regexp.MustCompile(`\(cost=([\d.]+?)\.\.([\d.]+?) rows=(\d+) width=(\d+)\)`)
var analyzeActualReg = regexp.MustCompile(`\(actual time=([\d.]+?)\.\.([\d.]+?) rows=(\d+) loops=(\d+)\)`)

func ParseExplainAnalyzeResultLineNode(lines []string) (int, *ExplainAnalyzeResultNode, error) {
	resLine := &ExplainAnalyzeResultNode{
		Lines:     lines,
		TableName: extractTableFromExplainAnalyzeResultLine(lines),
	}

	mLine := getMetricLine(lines)
	// c.f.) https://github.com/postgres/postgres/blob/REL_14_5/src/backend/commands/explain.c#L1583
	em := analyzeEstimationReg.FindStringSubmatch(mLine)
	if em != nil {
		initCost, err := strconv.ParseFloat(em[1], 10)
		if err != nil {
			return 0, nil, errors.Wrap(err, fmt.Sprintf("failed to parse analyze(%s): %s", em[1], mLine))
		}
		resLine.EstimatedInitCost = null.FloatFrom(initCost)
		cost, err := strconv.ParseFloat(em[2], 10)
		if err != nil {
			return 0, nil, errors.Wrap(err, fmt.Sprintf("failed to parse analyze(%s): %s", em[2], mLine))
		}
		resLine.EstimatedCost = null.FloatFrom(cost)

		eRowNum, err := strconv.ParseInt(em[3], 10, 0)
		if err != nil {
			return 0, nil, errors.Wrap(err, fmt.Sprintf("failed to parse analyze(%s): %s", em[3], mLine))
		}
		resLine.EstimatedReturnedRows = null.IntFrom(eRowNum)

		eWidth, err := strconv.ParseInt(em[4], 10, 0)
		if err != nil {
			return 0, nil, errors.Wrap(err, fmt.Sprintf("failed to parse analyze(%s): %s", em[4], mLine))
		}
		resLine.EstimatedWidth = null.IntFrom(eWidth)
	}

	// c.f.)
	//	https://github.com/postgres/postgres/blob/REL_14_5/src/backend/commands/explain.c#L1624
	am := analyzeActualReg.FindStringSubmatch(mLine)
	if am != nil {
		timeFirst, err := strconv.ParseFloat(am[1], 10)
		if err != nil {
			return 0, nil, errors.Wrap(err, fmt.Sprintf("failed to parse analyze(%s): %s", am[1], mLine))
		}
		resLine.ActualTimeFirstRow = null.FloatFrom(timeFirst)

		time, err := strconv.ParseFloat(am[2], 10)
		if err != nil {
			return 0, nil, errors.Wrap(err, fmt.Sprintf("failed to parse analyze(%s): %s", am[2], mLine))
		}
		resLine.ActualTimeAvg = null.FloatFrom(time)

		rowNum, err := strconv.ParseInt(am[3], 10, 0)
		if err != nil {
			return 0, nil, errors.Wrap(err, fmt.Sprintf("failed to parse analyze(%s): %s", am[3], mLine))
		}
		resLine.ActualReturnedRows = null.IntFrom(rowNum)

		loopNum, err := strconv.ParseInt(am[4], 10, 0)
		if err != nil {
			return 0, nil, errors.Wrap(err, fmt.Sprintf("failed to parse analyze(%s): %s", am[4], mLine))
		}
		resLine.ActualLoopCount = null.IntFrom(loopNum)
	}

	spaces := prefixSpaceReg.FindString(mLine)
	return len(spaces), resLine, nil
}

func extractTableFromExplainAnalyzeResultLine(lines []string) string {
	for _, line := range lines {
		m := onReg.FindStringSubmatch(line)
		if m != nil {
			return m[1]
		}

		m = keyReg.FindStringSubmatch(line)
		if m != nil {
			return m[1]
		}
	}
	return ""
}

func getMetricLine(lines []string) string {
	for _, line := range lines {
		if metricLineReg.MatchString(line) {
			return line
		}
	}

	return ""
}

func (n *ExplainAnalyzeResultNode) EstimatedTotalTime() float64 {
	// Return ActualTimeAvg even node has "loop"
	// 	because ActualTimeAvg seems not the average time of loops for PostgreSQL even doc says so
	return n.ActualTimeAvg.ValueOrZero()
}

func (n *ExplainAnalyzeResultNode) String() string {
	return fmt.Sprintf(
		"[%s] estimate_cost=%v, estimate_rows=%v, width=%v, time=%v..%v, totalTime=%f, rows=%v, loops=%v\n%s",
		n.TableName,
		n.EstimatedCost, n.EstimatedReturnedRows, n.EstimatedWidth,
		n.ActualTimeFirstRow, n.ActualTimeAvg, n.EstimatedTotalTime(),
		n.ActualReturnedRows, n.ActualLoopCount,
		n.Lines,
	)
}

var titleReg = regexp.MustCompile(`\s*(->)?\s*(.+?)\s+(on |\()`)

func (n *ExplainAnalyzeResultNode) Title() string {
	mLine := getMetricLine(n.Lines)
	m := titleReg.FindStringSubmatch(mLine)
	if m == nil {
		return ""
	}
	return m[2]
}
