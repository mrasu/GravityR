package mysql

import (
	"fmt"
	"github.com/pkg/errors"
	"regexp"
	"strconv"
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

var tableReg = regexp.MustCompile(` on (\w+)`)
var analyzeEstimationReg = regexp.MustCompile(`\(cost=([\d.]+?)(\.\.([\d.]+?))? rows=(\d+)\)`)
var analyzeActualReg = regexp.MustCompile(`\(actual time=([\d.]+?)\.\.([\d.]+?) rows=(\d+) loops=(\d+)\)`)

func ParseAnalyzeResult(line string) (*AnalyzeResultLine, error) {
	res := &AnalyzeResultLine{Text: line}

	tm := tableReg.FindStringSubmatch(line)
	if tm != nil {
		res.TableName = tm[1]
	}

	// c.f.) https://github.com/mysql/mysql-server/blob/8d8c986e5716e38cb776b627a8eee9e92241b4ce/sql/join_optimizer/explain_access_path.cc#L1076
	em := analyzeEstimationReg.FindStringSubmatch(line)
	if em != nil {
		if em[2] == "" {
			cost, err := strconv.ParseFloat(em[1], 64)
			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("failed to parse analyze(%s): %s", em[1], line))
			}
			res.EstimatedCost = cost
		} else {
			initCost, err := strconv.ParseFloat(em[1], 64)
			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("failed to parse analyze(%s): %s", em[1], line))
			}
			res.EstimatedInitCost = initCost
			cost, err := strconv.ParseFloat(em[3], 64)
			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("failed to parse analyze(%s): %s", em[3], line))
			}
			res.EstimatedCost = cost
		}

		eRowNum, err := strconv.Atoi(em[4])
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to parse analyze(%s): %s", em[4], line))
		}
		res.EstimatedReturnedRows = eRowNum
	}

	// c.f.) https://github.com/mysql/mysql-server/blob/6846e6b2f72931991cc9fd589dc9946ea2ab58c9/sql/iterators/timing_iterator.h#L159
	am := analyzeActualReg.FindStringSubmatch(line)
	if am != nil {
		timeFirst, err := strconv.ParseFloat(am[1], 64)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to parse analyze(%s): %s", am[1], line))
		}
		res.ActualTimeFirstRow = timeFirst

		time, err := strconv.ParseFloat(am[2], 64)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to parse analyze(%s): %s", am[2], line))
		}
		res.ActualTimeAvg = time

		rowNum, err := strconv.Atoi(am[3])
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to parse analyze(%s): %s", am[3], line))
		}
		res.ActualReturnedRows = rowNum

		loopNum, err := strconv.Atoi(am[4])
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to parse analyze(%s): %s", am[4], line))
		}
		res.ActualLoopCount = loopNum
	}

	return res, nil
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
