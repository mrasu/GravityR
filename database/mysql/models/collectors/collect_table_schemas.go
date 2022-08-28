package collectors

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/mrasu/GravityR/database/db_models"
	"github.com/mrasu/GravityR/lib"
	"github.com/pkg/errors"
)

type column struct {
	ColumnName string `db:"COLUMN_NAME"`
	ColumnKey  string `db:"COLUMN_KEY"`
	TableName  string `db:"TABLE_NAME"`
}

const schemaFetchQuery = `
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

func CollectTableSchemas(db *sqlx.DB, database string, tables []string) ([]*db_models.TableSchema, error) {
	query, args, err := sqlx.In(schemaFetchQuery, database, tables)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build query to get table schema")
	}

	var cols []column
	err = db.Select(&cols, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute query to get table schema")
	}

	schemas := db_models.CreateTableSchemas(tables, cols, func(c column) (string, string, bool) {
		return c.TableName, c.ColumnName, c.ColumnKey == "PRI"
	})

	for i, table := range tables {
		if schemas[i] == nil {
			return nil, lib.NewUnsupportedError(fmt.Sprintf("unknown table found. perhaps using VIEW? not supporting: %s", table))
		}
	}

	return schemas, nil
}
