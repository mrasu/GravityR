package mysql

const columnFetchQuery = `
SELECT
	COLUMN_NAME,
	COLUMN_KEY,
	TABLE_NAME 
FROM
	information_schema.columns
WHERE
  	TABLE_SCHEMA = ? AND
  	TABLE_NAME IN (?)
ORDER BY TABLE_NAME, ORDINAL_POSITION
`

type ColumnInfo struct {
	ColumnName string `db:"COLUMN_NAME"`
	ColumnKey  string `db:"COLUMN_KEY"`
	TableName  string `db:"TABLE_NAME"`
}
