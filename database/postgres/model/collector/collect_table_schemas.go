package collector

import (
	"fmt"
	"github.com/mrasu/GravityR/database/common_model"
	"github.com/mrasu/GravityR/infra/postgres"
	"github.com/mrasu/GravityR/lib"
)

func CollectTableSchemas(db *postgres.DB, schema string, tables []string) ([]*common_model.TableSchema, error) {
	cols, err := db.GetTableColumns(schema, tables)
	if err != nil {
		return nil, err
	}

	schemas := common_model.CreateTableSchemas(tables, cols, func(c *postgres.ColumnInfo) (string, string, bool) {
		return c.TableName, c.ColumnName, c.IsPK
	})

	for i, table := range tables {
		if schemas[i] == nil {
			return nil, lib.NewUnsupportedError(fmt.Sprintf("unknown table found. perhaps using VIEW? not supporting: %s", table))
		}
	}

	return schemas, nil
}