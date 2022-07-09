package inspectors

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/mrasu/GravityR/database/mysql/models"
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

func CollectTableSchema(db *sqlx.DB, database string, tables []string) ([]*models.TableSchema, error) {
	query, args, err := sqlx.In(schemaFetchQuery, database, tables)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build query to get table schema")
	}
	rows, err := db.Queryx(query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get table schema")
	}

	ts := map[string]*models.TableSchema{}
	for rows.Next() {
		var c column
		err = rows.StructScan(&c)
		if err != nil {
			return nil, errors.Wrap(err, "failed to wrap result from information_schema.columns")
		}
		var t *models.TableSchema
		if et, ok := ts[c.TableName]; ok {
			t = et
		} else {
			t = &models.TableSchema{
				Name: c.TableName,
			}
			ts[c.TableName] = t
		}

		t.Columns = append(t.Columns, &models.ColumnSchema{
			Name: c.ColumnName,
		})

		if c.ColumnKey == "PRI" {
			t.PrimaryKeys = append(t.PrimaryKeys, c.ColumnName)
		}
	}

	var res []*models.TableSchema
	for _, table := range tables {
		if t, ok := ts[table]; ok {
			res = append(res, t)
		} else {
			return nil, lib.NewUnsupportedError(fmt.Sprintf("unknown table. perhaps using VIEW? not supporting: %s", table))
		}
	}
	return res, nil
}
