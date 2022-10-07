package mysql_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mrasu/GravityR/infra/mysql"
	"github.com/mrasu/GravityR/thelper"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDB_Explain(t *testing.T) {
	tests := []struct {
		name          string
		explainResult string
		expected      string
	}{
		{
			name: "single join(SELECT name FROM users INNER JOIN tasks ON ...)",
			explainResult: `
-> Nested loop inner join  (cost=6378854.77 rows=6361711) (actual time=0.147..2483.756 rows=6553600 loops=1)
    -> Table scan on tasks  (cost=659260.39 rows=6361711) (actual time=0.074..1138.883 rows=6553600 loops=1)
    -> Single-row index lookup on users using PRIMARY (id=tasks.user_id)  (cost=0.80 rows=1) (actual time=0.000..0.000 rows=1 loops=6553600)
`,
			expected: `
-> Nested loop inner join  (cost=6378854.77 rows=6361711) (actual time=0.147..2483.756 rows=6553600 loops=1)
    -> Table scan on tasks  (cost=659260.39 rows=6361711) (actual time=0.074..1138.883 rows=6553600 loops=1)
    -> Single-row index lookup on users using PRIMARY (id=tasks.user_id)  (cost=0.80 rows=1) (actual time=0.000..0.000 rows=1 loops=6553600)
`,
		},
		{
			name:          "Empty",
			explainResult: "",
			expected:      "",
		},
	}

	dummyQuery := "SELECT ... FROM ..."
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			thelper.MockMysqlDB(t, func(db *mysql.DB, mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"EXPLAIN"}).AddRow(tt.explainResult)
				mock.ExpectQuery("EXPLAIN ANALYZE FORMAT=TREE " + dummyQuery).WillReturnRows(rows)

				res, err := db.ExplainWithAnalyze(dummyQuery)
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, tt.expected, res)
			})
		})
	}
}

var tableColumns = []*mysql.ColumnInfo{
	{ColumnName: "id", ColumnKey: "PRI", TableName: "users"},
	{ColumnName: "name", TableName: "users"},
	{ColumnName: "email", TableName: "users"},
	{ColumnName: "password", TableName: "users"},
	{ColumnName: "created_at", TableName: "users"},
	{ColumnName: "updated_at", TableName: "users"},

	{ColumnName: "id", ColumnKey: "PRI", TableName: "tasks"},
	{ColumnName: "user_id", TableName: "tasks"},
	{ColumnName: "title", TableName: "tasks"},
	{ColumnName: "description", TableName: "tasks"},
	{ColumnName: "status", TableName: "tasks"},
	{ColumnName: "created_at", TableName: "tasks"},
	{ColumnName: "updated_at", TableName: "tasks"},
}

const tableSchemaQuery = "SELECT\\s+COLUMN_NAME,\\s+COLUMN_KEY,\\s+TABLE_NAME\\s+FROM\\s+information_schema.columns"

func TestDB_GetTableColumns(t *testing.T) {
	tests := []struct {
		name         string
		tables       []string
		expectedCols []*mysql.ColumnInfo
	}{
		{
			name:         "one table",
			tables:       []string{"users"},
			expectedCols: lo.Filter(tableColumns, func(v *mysql.ColumnInfo, _ int) bool { return v.TableName == "users" }),
		},
		{
			name:         "multiple table",
			tables:       []string{"users", "tasks"},
			expectedCols: tableColumns,
		},
		{
			name:         "not existing table",
			tables:       []string{"missing"},
			expectedCols: []*mysql.ColumnInfo{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			thelper.MockMysqlDB(t, func(db *mysql.DB, mock sqlmock.Sqlmock) {
				mockTableSchemaQuery(mock, tt.tables)

				cols, err := db.GetTableColumns(thelper.DBName, tt.tables)
				assert.NoError(t, err)

				assert.ElementsMatch(t, cols, cols)
			})
		})
	}
}

func mockTableSchemaQuery(mock sqlmock.Sqlmock, tables []string) {
	rows := sqlmock.NewRows([]string{"COLUMN_NAME", "COLUMN_KEY", "TABLE_NAME"})
	for _, table := range tables {
		for _, c := range tableColumns {
			if c.TableName == table {
				rows.AddRow(c.ColumnName, c.ColumnKey, table)
			}
		}
	}
	mock.ExpectQuery(tableSchemaQuery).WillReturnRows(rows)
}
