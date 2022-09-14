package collectors

import (
	"fmt"
	"github.com/mrasu/GravityR/database/db_models"
	"github.com/mrasu/GravityR/infra/hasura"
	"github.com/mrasu/GravityR/lib"
)

func CollectTableSchemas(cli *hasura.Client, schema string, tables []string) ([]*db_models.TableSchema, error) {
	cols, err := cli.GetTableColumns(schema, tables)
	if err != nil {
		return nil, err
	}

	schemas := db_models.CreateTableSchemas(tables, cols, func(c *hasura.ColumnInfo) (string, string, bool) {
		return c.TableName, c.ColumnName, c.IsPK
	})

	for i, table := range tables {
		if schemas[i] == nil {
			return nil, lib.NewUnsupportedError(fmt.Sprintf("unknown table found. perhaps using VIEW? not supporting: %s", table))
		}
	}

	return schemas, nil
}
