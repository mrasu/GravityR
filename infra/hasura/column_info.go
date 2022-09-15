package hasura

import (
	"fmt"
	"github.com/mrasu/GravityR/lib"
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

const columnFetchQuery = `
SELECT
	pg_class.relname AS table_name,
	pg_attribute.attname AS column_name,
	COALESCE(pg_index.indisprimary, FALSE) AS is_pk
FROM
	pg_class
	INNER JOIN pg_attribute ON pg_attribute.attrelid = pg_class.oid
	INNER JOIN pg_namespace ON pg_class.relnamespace = pg_namespace.oid
	LEFT OUTER JOIN pg_index ON pg_class.oid = pg_index.indrelid AND pg_attribute.attnum = ANY(pg_index.indkey)
WHERE
	pg_namespace.nspname = ? AND
	pg_class.relname IN (?) AND
	pg_class.relkind = 'r' AND
	pg_attribute.attnum >= 0
ORDER BY
	pg_class.relname,
	pg_attribute.attnum
`

type ColumnInfo struct {
	ColumnName string
	IsPK       bool
	TableName  string
}

var wordOnlyReg = regexp.MustCompile(`\A\w+\z`)

func buildColumnFetchQuery(schema string, tableNames []string) (string, error) {
	if !wordOnlyReg.MatchString(schema) {
		return "", errors.Errorf("schema name can contains only a-z or _: %s", schema)
	}
	for _, t := range tableNames {
		if !wordOnlyReg.MatchString(t) {
			return "", errors.Errorf("table name can contains only a-z or _: %s", t)
		}
	}
	ts := lib.Join(tableNames, ",", func(v string) string { return fmt.Sprintf("'%s'", v) })
	q := strings.Replace(columnFetchQuery, "?", fmt.Sprintf("'%s'", schema), 1)
	q = strings.Replace(q, "?", ts, 1)

	return q, nil
}

func parseColumnFetchResult(res *RunSQLResponse) []*ColumnInfo {
	var cols []*ColumnInfo
	for i, r := range res.Result {
		if i == 0 {
			// first row shows name of columns
			continue
		}

		isPk := false
		if r[2] == "t" {
			isPk = true
		}

		cols = append(cols, &ColumnInfo{
			TableName:  r[0],
			ColumnName: r[1],
			IsPK:       isPk,
		})
	}

	return cols
}
