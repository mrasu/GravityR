package collector

import (
	"fmt"
	"github.com/mrasu/GravityR/database"
	"github.com/mrasu/GravityR/infra/mysql"
	"github.com/mrasu/GravityR/lib"
	"github.com/rs/zerolog/log"
)

func CollectTableSchemas(db *mysql.DB, dbName string, tables []string) ([]*database.TableSchema, error) {
	cols, err := db.GetTableColumns(dbName, tables)
	if err != nil {
		return nil, err
	}

	schemas := database.CreateTableSchemas(tables, cols, func(c *mysql.ColumnInfo) (string, string, bool) {
		return c.TableName, c.ColumnName, c.ColumnKey == "PRI"
	})

	for i, table := range tables {
		if schemas[i] == nil {
			return nil, lib.NewUnsupportedError(fmt.Sprintf("unknown table found. perhaps using VIEW? not supporting: %s", table))
		}
	}

	log.Debug().Msg("Table schemas:")
	for i, s := range schemas {
		log.Printf("\t%d. %s", i, s.TableDescription())
	}

	return schemas, nil
}
