package hasura

import (
	"fmt"
	"github.com/lib/pq"
	"github.com/mrasu/GravityR/lib"
	"github.com/pkg/errors"
	"strings"
)

const indexFetchQuery = `
SELECT 
	t.relname AS table_name,
	ARRAY(
		SELECT PG_GET_INDEXDEF(pg_index.indexrelid, k + 1, TRUE)
		FROM GENERATE_SUBSCRIPTS(pg_index.indkey, 1) AS k
		ORDER BY k
	) AS column_names
FROM pg_index
	INNER JOIN pg_class AS i ON pg_index.indexrelid = i.oid
	INNER JOIN pg_class AS t ON pg_index.indrelid = t.oid
	INNER JOIN pg_namespace ON i.relnamespace = pg_namespace.oid
WHERE
	pg_namespace.nspname = ? AND
	t.relname IN (?)
ORDER BY
	t.relname,
	i.relname
`

type IndexInfo struct {
	TableName string
	Columns   []string
}

func buildIndexFetchQuery(schema string, tableNames []string) (string, error) {
	if !wordOnlyReg.MatchString(schema) {
		return "", errors.Errorf("schema name can contains only a-z or _: %s", schema)
	}
	for _, t := range tableNames {
		if !wordOnlyReg.MatchString(t) {
			return "", errors.Errorf("table name can contains only a-z or _: %s", t)
		}
	}
	ts := lib.Join(tableNames, ",", func(v string) string { return fmt.Sprintf("'%s'", v) })
	q := strings.Replace(indexFetchQuery, "?", fmt.Sprintf("'%s'", schema), 1)
	q = strings.Replace(q, "?", ts, 1)

	return q, nil
}

func parseIndexFetchResult(res *RunSQLResponse) ([]*IndexInfo, error) {
	var infos []*IndexInfo
	for i, r := range res.Result {
		if i == 0 {
			// first row shows name of columns
			continue
		}

		cols := pq.StringArray{}
		err := cols.Scan(r[1])
		if err != nil {
			return nil, errors.Wrap(err, "unexpected result from hasura to fetch indexes: not array of string for column names")
		}

		infos = append(infos, &IndexInfo{
			TableName: r[0],
			Columns:   cols,
		})
	}

	return infos, nil
}
