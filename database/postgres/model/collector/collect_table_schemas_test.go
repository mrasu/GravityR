package collector_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mrasu/GravityR/database/common_model"
	"github.com/mrasu/GravityR/database/postgres/model/collector"
	"github.com/mrasu/GravityR/infra/postgres"
	"github.com/mrasu/GravityR/thelper"
	"github.com/stretchr/testify/assert"
	"testing"
)

var tableColumns = map[string][]*common_model.ColumnSchema{
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

const tableSchemaQuery = "SELECT\\s+pg_class.relname AS table_name,\\s+pg_attribute.attname AS column_name,\\s+COALESCE\\(pg_index.indisprimary, FALSE\\) AS is_pk\\s+FROM\\s+pg_class"

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
			thelper.MockPostgresDB(t, func(db *postgres.DB, mock sqlmock.Sqlmock) {
				mockTableSchemaQuery(mock, tt.tables)

				ts, err := collector.CollectTableSchemas(db, thelper.DBName, tt.tables)
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
	rows := sqlmock.NewRows([]string{"table_name", "column_name", "is_pk"})
	for _, table := range tables {
		for _, c := range tableColumns[table] {
			rows.AddRow(table, c.Name, c.Name == "id")
		}
	}
	mock.ExpectQuery(tableSchemaQuery).WillReturnRows(rows)
}

func TestCollectTableSchemas_Error(t *testing.T) {
	thelper.MockPostgresDB(t, func(db *postgres.DB, mock sqlmock.Sqlmock) {
		emptyRows := sqlmock.NewRows([]string{"table_name", "column_name", "is_pk"})
		mock.ExpectQuery(tableSchemaQuery).WillReturnRows(emptyRows)

		ts, err := collector.CollectTableSchemas(db, thelper.DBName, []string{"missing"})
		assert.Nil(t, ts)
		assert.Error(t, err)
	})
}
