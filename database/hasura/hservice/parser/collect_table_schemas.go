package parser

import (
	"fmt"
	"github.com/mrasu/GravityR/database/dservice/dparser"
	"github.com/mrasu/GravityR/infra/hasura"
	"github.com/mrasu/GravityR/lib"
	"github.com/rs/zerolog/log"
)

func CollectTableSchemas(cli *hasura.Client, schema string, tables []string) ([]*dparser.TableSchema, error) {
	cols, err := cli.GetTableColumns(schema, tables)
	if err != nil {
		return nil, err
	}

	schemas := dparser.CreateTableSchemas(tables, cols, func(c *hasura.ColumnInfo) (string, string, bool) {
		return c.TableName, c.ColumnName, c.IsPK
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
