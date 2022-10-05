package postgres

import (
	"github.com/lib/pq"
	"github.com/samber/lo"
)

const indexFetchQuery = `
SELECT 
	i.relname AS index_name,
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

type queryIndexInfo struct {
	IndexName   string         `db:"index_name"`
	TableName   string         `db:"table_name"`
	ColumnNames pq.StringArray `db:"column_names"`
}

func toIndexInfo(queryInfos []*queryIndexInfo) []*IndexInfo {
	return lo.Map(queryInfos, func(qi *queryIndexInfo, _ int) *IndexInfo {
		return &IndexInfo{
			TableName: qi.TableName,
			Columns:   qi.ColumnNames,
		}
	})
}
