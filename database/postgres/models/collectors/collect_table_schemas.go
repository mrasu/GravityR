package collectors

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/mrasu/GravityR/database/db_models"
	"github.com/mrasu/GravityR/lib"
	"github.com/pkg/errors"
)

type column struct {
	ColumnName string `db:"column_name"`
	IsPK       bool   `db:"is_pk"`
	TableName  string `db:"table_name"`
}

const schemaFetchQuery = `
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

func CollectTableSchemas(db *sqlx.DB, schema string, tables []string) ([]*db_models.TableSchema, error) {
	query, args, err := sqlx.In(schemaFetchQuery, schema, tables)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build query to get table schema")
	}
	query = db.Rebind(query)

	var cols []column
	err = db.Select(&cols, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute query to get table schema")
	}

	schemas := db_models.CreateTableSchemas(tables, cols, func(c column) (string, string, bool) {
		return c.TableName, c.ColumnName, c.IsPK
	})

	for i, table := range tables {
		if schemas[i] == nil {
			return nil, lib.NewUnsupportedError(fmt.Sprintf("unknown table found. perhaps using VIEW? not supporting: %s", table))
		}
	}

	return schemas, nil
}
