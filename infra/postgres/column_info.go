package postgres

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
	pg_attribute.attnum;
`

type ColumnInfo struct {
	ColumnName string `db:"column_name"`
	IsPK       bool   `db:"is_pk"`
	TableName  string `db:"table_name"`
}
