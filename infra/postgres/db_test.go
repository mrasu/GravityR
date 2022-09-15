package postgres_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mrasu/GravityR/infra/postgres"
	"github.com/mrasu/GravityR/thelper"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestDB_ExplainWithAnalyze(t *testing.T) {
	tests := []struct {
		name          string
		explainResult string
		expected      []string
	}{
		{
			name: "simple(SELECT name FROM users)",
			explainResult: `
Seq Scan on users  (cost=0.00..1400.36 rows=65536 width=9) (actual time=0.007..4.327 rows=65536 loops=1)
  Buffers: shared hit=745
Planning Time: 0.031 ms
Execution Time: 6.023 ms
`,
			expected: []string{
				"Seq Scan on users  (cost=0.00..1400.36 rows=65536 width=9) (actual time=0.007..4.327 rows=65536 loops=1)",
				"  Buffers: shared hit=745",
				"Planning Time: 0.031 ms",
				"Execution Time: 6.023 ms",
			},
		},
		{
			name:          "Empty",
			explainResult: "",
			expected:      []string{""},
		},
	}

	dummyQuery := "SELECT ... FROM ..."
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			thelper.MockPostgresDB(t, func(db *postgres.DB, mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"QUERY PLAN"})
				for _, line := range strings.Split(strings.Trim(tt.explainResult, "\n"), "\n") {
					rows.AddRow(line)
				}
				mock.ExpectQuery("EXPLAIN \\(ANALYZE, BUFFERS\\) " + dummyQuery).WillReturnRows(rows)

				res, err := db.ExplainWithAnalyze(dummyQuery)
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, tt.expected, res)
			})
		})
	}
}

var tableColumns = []*postgres.ColumnInfo{
	{ColumnName: "id", IsPK: true, TableName: "users"},
	{ColumnName: "name", TableName: "users"},
	{ColumnName: "email", TableName: "users"},
	{ColumnName: "password", TableName: "users"},
	{ColumnName: "created_at", TableName: "users"},
	{ColumnName: "updated_at", TableName: "users"},

	{ColumnName: "id", IsPK: true, TableName: "todos"},
	{ColumnName: "user_id", TableName: "todos"},
	{ColumnName: "title", TableName: "todos"},
	{ColumnName: "description", TableName: "todos"},
	{ColumnName: "status", TableName: "todos"},
	{ColumnName: "created_at", TableName: "todos"},
	{ColumnName: "updated_at", TableName: "todos"},
}

const tableSchemaQuery = "SELECT\\s+pg_class.relname AS table_name,\\s+pg_attribute.attname AS column_name,\\s+COALESCE\\(pg_index.indisprimary, FALSE\\) AS is_pk\\s+FROM\\s+pg_class"

func TestDB_GetTableColumns(t *testing.T) {
	tests := []struct {
		name         string
		tables       []string
		expectedCols []*postgres.ColumnInfo
	}{
		{
			name:         "one table",
			tables:       []string{"users"},
			expectedCols: lo.Filter(tableColumns, func(v *postgres.ColumnInfo, _ int) bool { return v.TableName == "users" }),
		},
		{
			name:         "multiple table",
			tables:       []string{"users", "todos"},
			expectedCols: tableColumns,
		},
		{
			name:         "not existing table",
			tables:       []string{"missing"},
			expectedCols: []*postgres.ColumnInfo{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			thelper.MockPostgresDB(t, func(db *postgres.DB, mock sqlmock.Sqlmock) {
				mockTableSchemaQuery(mock, tt.tables)

				cols, err := db.GetTableColumns(thelper.DBName, tt.tables)
				assert.NoError(t, err)

				assert.ElementsMatch(t, cols, cols)
			})
		})
	}
}

func mockTableSchemaQuery(mock sqlmock.Sqlmock, tables []string) {
	rows := sqlmock.NewRows([]string{"table_name", "column_name", "is_pk"})
	for _, table := range tables {
		for _, c := range tableColumns {
			if c.TableName == table {
				rows.AddRow(c.TableName, c.ColumnName, c.IsPK)
			}
		}
	}
	mock.ExpectQuery(tableSchemaQuery).WillReturnRows(rows)
}
