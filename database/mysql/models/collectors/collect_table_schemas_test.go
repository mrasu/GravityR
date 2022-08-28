package collectors_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/mrasu/GravityR/database/db_models"
	"github.com/mrasu/GravityR/database/mysql/models/collectors"
	"github.com/mrasu/GravityR/thelper"
	"github.com/stretchr/testify/assert"
	"testing"
)

var tableColumns = map[string][]*db_models.ColumnSchema{
	"users": {
		{Name: "id"},
		{Name: "name"},
		{Name: "email"},
		{Name: "password"},
		{Name: "created_at"},
		{Name: "updated_at"},
	},
	"todos": {
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
			tables: []string{"users", "todos"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			thelper.MockDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				mockTableSchemaQuery(mock, tt.tables)

				ts, err := collectors.CollectTableSchemas(db, thelper.DbName, tt.tables)
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
	thelper.MockDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
		emptyRows := sqlmock.NewRows([]string{"COLUMN_NAME", "COLUMN_KEY", "TABLE_NAME"})
		mock.ExpectQuery(tableSchemaQuery).WillReturnRows(emptyRows)

		ts, err := collectors.CollectTableSchemas(db, thelper.DbName, []string{"missing"})
		assert.Nil(t, ts)
		assert.Error(t, err)
	})
}
