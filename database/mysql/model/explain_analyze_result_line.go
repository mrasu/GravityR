package model

import (
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v4"
	"regexp"
	"strconv"
	"strings"
)

// ExplainAnalyzeResultLine is the parsed information from `EXPLAIN ANALYZE` query
// c.f)
//	https://dev.mysql.com/doc/refman/8.0/en/explain.html#explain-analyze
//	https://github.com/mysql/mysql-server/blob/8d8c986e5716e38cb776b627a8eee9e92241b4ce/sql/join_optimizer/explain_access_path.cc#L475
type ExplainAnalyzeResultLine struct {
	Text      string
	TableName string

	EstimatedInitCost     null.Float
	EstimatedCost         null.Float
	EstimatedReturnedRows null.Int

	ActualTimeFirstRow null.Float
	ActualTimeAvg      null.Float
	ActualReturnedRows null.Int
	ActualLoopCount    null.Int
}

var prefixIndentReg = regexp.MustCompile(`^\s*`)
var onReg = regexp.MustCompile(` on (\w+)`)
var filterReg = regexp.MustCompile(` Filter: \((\w+)\.`)
var sortReg = regexp.MustCompile(` Sort: (\w+)\.`)
var aggregateReg = regexp.MustCompile(` Aggregate: \w+\((\w+)\.`)
var analyzeEstimationReg = regexp.MustCompile(`\(cost=([\d.]+?)(\.\.([\d.]+?))?( rows=(\d+))?\)`)
var analyzeActualReg = regexp.MustCompile(`\(actual time=([\d.]+?)\.\.([\d.]+?) rows=(\d+) loops=(\d+)\)`)

func ParseExplainAnalyzeResultLine(line string) (int, *ExplainAnalyzeResultLine, error) {
	resLine := &ExplainAnalyzeResultLine{
		Text:      line,
		TableName: extractTableFromExplainAnalyzeResultLine(line),
	}

	// c.f.) https://github.com/mysql/mysql-server/blob/8d8c986e5716e38cb776b627a8eee9e92241b4ce/sql/join_optimizer/explain_access_path.cc#L1076
	em := analyzeEstimationReg.FindStringSubmatch(line)
	if em != nil {
		if em[2] == "" {
			cost, err := strconv.ParseFloat(em[1], 64)
			if err != nil {
				return 0, nil, errors.Wrap(err, fmt.Sprintf("failed to parse analyze(%s): %s", em[1], line))
			}
			resLine.EstimatedCost = null.FloatFrom(cost)
		} else {
			initCost, err := strconv.ParseFloat(em[1], 64)
			if err != nil {
				return 0, nil, errors.Wrap(err, fmt.Sprintf("failed to parse analyze(%s): %s", em[1], line))
			}
			resLine.EstimatedInitCost = null.FloatFrom(initCost)
			cost, err := strconv.ParseFloat(em[3], 64)
			if err != nil {
				return 0, nil, errors.Wrap(err, fmt.Sprintf("failed to parse analyze(%s): %s", em[3], line))
			}
			resLine.EstimatedCost = null.FloatFrom(cost)
		}

		if em[4] != "" {
			eRowNum, err := strconv.ParseInt(em[5], 10, 64)
			if err != nil {
				return 0, nil, errors.Wrap(err, fmt.Sprintf("failed to parse analyze(%s): %s", em[4], line))
			}
			resLine.EstimatedReturnedRows = null.IntFrom(eRowNum)
		}
	}

	// c.f.) https://github.com/mysql/mysql-server/blob/6846e6b2f72931991cc9fd589dc9946ea2ab58c9/sql/iterators/timing_iterator.h#L159
	am := analyzeActualReg.FindStringSubmatch(line)
	if am != nil {
		timeFirst, err := strconv.ParseFloat(am[1], 64)
		if err != nil {
			return 0, nil, errors.Wrap(err, fmt.Sprintf("failed to parse analyze(%s): %s", am[1], line))
		}
		resLine.ActualTimeFirstRow = null.FloatFrom(timeFirst)

		time, err := strconv.ParseFloat(am[2], 64)
		if err != nil {
			return 0, nil, errors.Wrap(err, fmt.Sprintf("failed to parse analyze(%s): %s", am[2], line))
		}
		resLine.ActualTimeAvg = null.FloatFrom(time)

		rowNum, err := strconv.ParseInt(am[3], 10, 64)
		if err != nil {
			return 0, nil, errors.Wrap(err, fmt.Sprintf("failed to parse analyze(%s): %s", am[3], line))
		}
		resLine.ActualReturnedRows = null.IntFrom(rowNum)

		loopNum, err := strconv.ParseInt(am[4], 10, 64)
		if err != nil {
			return 0, nil, errors.Wrap(err, fmt.Sprintf("failed to parse analyze(%s): %s", am[4], line))
		}
		resLine.ActualLoopCount = null.IntFrom(loopNum)
	}

	pim := prefixIndentReg.FindString(line)
	nestSize := len(pim) / 4
	return nestSize, resLine, nil
}

func extractTableFromExplainAnalyzeResultLine(line string) string {
	m := onReg.FindStringSubmatch(line)
	if m != nil {
		return m[1]
	}

	m = filterReg.FindStringSubmatch(line)
	if m != nil {
		return m[1]
	}

	m = sortReg.FindStringSubmatch(line)
	if m != nil {
		return m[1]
	}

	m = aggregateReg.FindStringSubmatch(line)
	if m != nil {
		return m[1]
	}

	return ""
}

func (l *ExplainAnalyzeResultLine) EstimatedTotalTime() float64 {
	if l.ActualTimeAvg.ValueOrZero() == 0 && l.ActualLoopCount.ValueOrZero() > 1000 {
		// ActualTimeAvg doesn't show the value less than 0.000, however,
		// when the number of loop is much bigger, time can be meaningful even the result of multiplication is zero.
		// To handle the problem, assume ActualTimeAvg is some less than 0.000
		return 0.0001 * float64(l.ActualLoopCount.ValueOrZero())
	} else {
		return l.ActualTimeAvg.ValueOrZero() * float64(l.ActualLoopCount.ValueOrZero())
	}
}

func (l *ExplainAnalyzeResultLine) String() string {
	return fmt.Sprintf(
		"[%s] estimate_cost=%v, estimate_rows=%v, time=%v..%v, totalTime=%f, rows=%v, loops=%v\n%s",
		l.TableName,
		l.EstimatedCost, l.EstimatedReturnedRows,
		l.ActualTimeFirstRow, l.ActualTimeAvg, l.EstimatedTotalTime(),
		l.ActualReturnedRows, l.ActualLoopCount,
		l.Text,
	)
}

var titleReg = regexp.MustCompile(`->\s+(.+)\s+\(`)
var titleWordReg = regexp.MustCompile(`\w+`)

func (l *ExplainAnalyzeResultLine) Title() string {
	m := titleReg.FindStringSubmatch(l.Text)
	if m == nil {
		wm := titleWordReg.FindStringSubmatch(l.Text)
		if wm == nil {
			return ""
		} else {
			return wm[0]
		}
	}

	txt := m[1]
	if strings.Contains(txt, ":") {
		t := strings.Split(txt, ":")
		return strings.TrimSpace(t[0])
	}
	if strings.Contains(txt, " join ") {
		t := strings.Split(txt, " join ")
		return strings.TrimSpace(t[0]) + " join"
	}
	if strings.Contains(txt, " on ") {
		t := strings.Split(txt, " on ")
		return strings.TrimSpace(t[0])
	}
	return strings.Split(txt, " ")[0]
}
