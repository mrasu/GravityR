package suggest

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mrasu/GravityR/infra/mysql"
	"github.com/mrasu/GravityR/thelper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func Test_runMysqlSuggest(t *testing.T) {
	mr := mysqlRunner{
		query: `
SELECT
	name,
	t.description
FROM
	users
	INNER JOIN tasks AS t ON users.id = t.user_id
WHERE
	users.email = 'test31776@example.com'
`,
		runsExamination: true,
		indexTargets:    []string{"users:name"},
	}

	thelper.InjectClientDist()
	thelper.MockMysqlDB(t, func(db *mysql.DB, mock sqlmock.Sqlmock) {
		m := mysqlMock{}
		m.mockAll(mock, mr.query)

		thelper.CreateTemp(t, "tmp.html", func(tmpfile *os.File) {
			err := mr.suggest(tmpfile.Name(), db, "test_db")
			require.NoError(t, err)

			html, err := os.ReadFile(tmpfile.Name())
			require.NoError(t, err)
			assert.Contains(t, string(html), "Nested loop inner join  (cost=5071411.14")
			assert.Contains(t, string(html), "users.email = 'test31776@example.com'")
			assert.Contains(t, string(html), `originalTimeMillis":0,"indexResults":[{"`)
		})
	})
}

type mysqlMock struct{}

func (m *mysqlMock) mockAll(mock sqlmock.Sqlmock, query string) {
	m.mockExplainQuery(mock)
	m.mockTableSchemaQuery(mock)
	m.mockIndexesQuery(mock)
	m.mockSelectQuery(mock, query)
	m.mockAddIndexQuery(mock)
	m.mockSelectQuery(mock, query)
	m.mockDropIndexQuery(mock)
}

func (m *mysqlMock) mockExplainQuery(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{"EXPLAIN"})
	rows.AddRow(`-> Nested loop inner join  (cost=5071411.14 rows=636171) (actual time=0.720..3211.448 rows=100 loops=1)
    -> Table scan on t  (cost=659832.11 rows=6361711) (actual time=0.072..1187.379 rows=6553600 loops=1)
    -> Filter: (users.email = 'test31776@example.com')  (cost=0.59 rows=0) (actual time=0.000..0.000 rows=0 loops=6553600)
        -> Single-row index lookup on users using PRIMARY (id=t.user_id)  (cost=0.59 rows=1) (actual time=0.000..0.000 rows=1 loops=6553600)
`)
	mock.ExpectQuery("EXPLAIN ANALYZE FORMAT=TREE").WillReturnRows(rows)
}

const tableSchemaQuery = "SELECT\\s+COLUMN_NAME,\\s+COLUMN_KEY,\\s+TABLE_NAME\\s+FROM\\s+information_schema.columns"

func (m *mysqlMock) mockTableSchemaQuery(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{"COLUMN_NAME", "COLUMN_KEY", "TABLE_NAME"})
	rows.AddRow("description", "", "tasks")
	rows.AddRow("user_id", "", "tasks")
	rows.AddRow("email", "", "users")
	rows.AddRow("id", "PRI", "users")
	rows.AddRow("name", "", "users")
	mock.ExpectQuery(tableSchemaQuery).WillReturnRows(rows)
}

const mysqlIndexesQuery = "SELECT\\s+TABLE_NAME,\\s+INDEX_NAME,\\s+COLUMN_NAME\\s+FROM\\s+information_schema.STATISTICS"

func (m *mysqlMock) mockIndexesQuery(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{"TABLE_NAME", "INDEX_NAME", "COLUMN_NAME"})
	rows.AddRow("tasks", "tasks_idx", "tdescriptionasks")
	rows.AddRow("users", "users_idx", "email")
	rows.AddRow("users", "users_idx", "id")
	mock.ExpectQuery(mysqlIndexesQuery).WillReturnRows(rows)
}

func (m *mysqlMock) mockSelectQuery(mock sqlmock.Sqlmock, query string) {
	mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(0, 0))
}

func (m *mysqlMock) mockAddIndexQuery(mock sqlmock.Sqlmock) {
	mock.ExpectExec("ALTER TABLE `users` ADD INDEX").WillReturnResult(sqlmock.NewResult(0, 0))
}

func (m *mysqlMock) mockDropIndexQuery(mock sqlmock.Sqlmock) {
	mock.ExpectExec("ALTER TABLE `users` DROP INDEX").WillReturnResult(sqlmock.NewResult(0, 0))
}
