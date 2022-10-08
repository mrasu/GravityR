package parser_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mrasu/GravityR/database/dservice/dparser"
	"github.com/mrasu/GravityR/database/mysql/mservice/parser"
	"github.com/mrasu/GravityR/infra/mysql"
	"github.com/mrasu/GravityR/thelper"
	"github.com/stretchr/testify/assert"
	"testing"
)

var tableColumns = map[string][]*dparser.ColumnSchema{
	"users": {
		{Name: "id"},
		{Name: "name"},
		{Name: "email"},
		{Name: "password"},
		{Name: "created_at"},
		{Name: "updated_at"},
	},
	"tasks": {
		{Name: "id"},
		{Name: "user_id"},
		{Name: "title"},
		{Name: "description"},
		{Name: "status"},
		{Name: "created_at"},
		{Name: "updated_at"},
	},
}

const tableSchemaQuery = "SELECT\\s+COLUMN_NAME,\\s+COLUMN_KEY,\\s+TABLE_NAME\\s+FROM\\s+information_schema.columns"

func TestCollectTableSchemas(t *testing.T) {
	tests := []struct {
		name   string
		tables []string
	}{
		{
			name:   "one table",
			tables: []string{"users"},
		},
		{
			name:   "multiple table",
			tables: []string{"users", "tasks"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			thelper.MockMysqlDB(t, func(db *mysql.DB, mock sqlmock.Sqlmock) {
				mockTableSchemaQuery(mock, tt.tables)

				ts, err := parser.CollectTableSchemas(db, thelper.DBName, tt.tables)
				assert.NoError(t, err)

				assert.Equal(t, len(ts), len(tt.tables))
				for i, name := range tt.tables {
					actual := ts[i]
					assert.Equal(t, actual.Name, name)
					assert.ElementsMatch(t, actual.Columns, tableColumns[name])
					assert.ElementsMatch(t, actual.PrimaryKeys, []string{"id"})
				}
			})
		})
	}
}

func mockTableSchemaQuery(mock sqlmock.Sqlmock, tables []string) {
	rows := sqlmock.NewRows([]string{"COLUMN_NAME", "COLUMN_KEY", "TABLE_NAME"})
	for _, table := range tables {
		for _, c := range tableColumns[table] {
			pk := ""
			if c.Name == "id" {
				pk = "PRI"
			}
			rows.AddRow(c.Name, pk, table)
		}
	}
	mock.ExpectQuery(tableSchemaQuery).WillReturnRows(rows)
}

func TestCollectTableSchemas_Error(t *testing.T) {
	thelper.MockMysqlDB(t, func(db *mysql.DB, mock sqlmock.Sqlmock) {
		emptyRows := sqlmock.NewRows([]string{"COLUMN_NAME", "COLUMN_KEY", "TABLE_NAME"})
		mock.ExpectQuery(tableSchemaQuery).WillReturnRows(emptyRows)

		ts, err := parser.CollectTableSchemas(db, thelper.DBName, []string{"missing"})
		assert.Nil(t, ts)
		assert.Error(t, err)
	})
}
