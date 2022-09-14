package collectors

import (
	"fmt"
	"github.com/mrasu/GravityR/database/db_models"
	"github.com/mrasu/GravityR/infra/mysql"
	"github.com/mrasu/GravityR/lib"
)

func CollectTableSchemas(db *mysql.DB, database string, tables []string) ([]*db_models.TableSchema, error) {
	cols, err := db.GetTableColumns(database, tables)
	if err != nil {
		return nil, err
	}

	schemas := db_models.CreateTableSchemas(tables, cols, func(c *mysql.ColumnInfo) (string, string, bool) {
		return c.TableName, c.ColumnName, c.ColumnKey == "PRI"
	})

	for i, table := range tables {
		if schemas[i] == nil {
			return nil, lib.NewUnsupportedError(fmt.Sprintf("unknown table found. perhaps using VIEW? not supporting: %s", table))
		}
	}

	return schemas, nil
}
