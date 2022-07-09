package mysql

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"strings"
)

func Analyze(db *sqlx.DB, query string) ([]*AnalyzeResultLine, error) {
	var lines []*AnalyzeResultLine

	rows, err := db.Query("EXPLAIN ANALYZE FORMAT=TREE " + query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select")
	}
	for rows.Next() {
		var tree string
		if err := rows.Scan(&tree); err != nil {
			return nil, errors.Wrap(err, "failed to Scan")
		}

		treeLines := strings.Split(tree, "\n")
		for _, line := range treeLines {
			if line == "" {
				continue
			}

			l, err := ParseAnalyzeResult(line)
			if err != nil {
				return nil, err
			}

			lines = append(lines, l)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "failed to Next")
	}

	return lines, nil
}
