package suggest

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mrasu/GravityR/infra/postgres"
	"github.com/mrasu/GravityR/thelper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func Test_runPostgresSuggest(t *testing.T) {
	v := postgresVarS{
		query: `
SELECT
	name,
	t.description
FROM
	users
	INNER JOIN todos AS t ON users.id = t.user_id
WHERE
	users.email = 'test31776@example.com'
`,
		runsExamination: true,
		indexTargets:    []string{"users:name"},
	}

	thelper.InjectClientDist()
	thelper.MockPostgresDB(t, func(db *postgres.DB, mock sqlmock.Sqlmock) {
		m := postgresMock{}
		m.mockAll(mock, v.query)

		thelper.CreateTemp(t, "tmp.html", func(tmpfile *os.File) {
			pr := postgresRunner{}
			err := pr.runSuggest(v, tmpfile.Name(), db, "test_db")
			require.NoError(t, err)

			html, err := os.ReadFile(tmpfile.Name())
			require.NoError(t, err)
			assert.Contains(t, string(html), "Hash Join  (cost=1564.21..103602.35")
			assert.Contains(t, string(html), "Planning Time: 1.652 ms")
			assert.Contains(t, string(html), "users.email = 'test31776@example.com'")
			assert.Contains(t, string(html), `originalTimeMillis":0,"indexResults":[{"`)
		})
	})
}

type postgresMock struct{}

func (m *postgresMock) mockAll(mock sqlmock.Sqlmock, query string) {
	m.mockExplainQuery(mock)
	m.mockTableSchemaQuery(mock)
	m.mockSelectQuery(mock, query)
	m.mockCreateIndexQuery(mock)
	m.mockSelectQuery(mock, query)
	m.mockDropIndexQuery(mock)
}

func (m *postgresMock) mockExplainQuery(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{"QUERY PLAN"})
	explains := []string{
		"Gather  (cost=2564.21..104612.35 rows=100 width=26) (actual time=723.551..729.418 rows=100 loops=1)",
		"  Workers Planned: 2",
		"  Workers Launched: 2",
		"  Buffers: shared hit=1608 read=68308",
		"  ->  Hash Join  (cost=1564.21..103602.35 rows=42 width=26) (actual time=585.084..700.687 rows=33 loops=3)",
		"        Hash Cond: (t.user_id = users.id)",
		"        Buffers: shared hit=1608 read=68308",
		"        ->  Parallel Seq Scan on todos t  (cost=0.00..94869.67 rows=2730667 width=21) (actual time=0.286..529.485 rows=2184533 loops=3)",
		"              Buffers: shared read=67563",
		"        ->  Hash  (cost=1564.20..1564.20 rows=1 width=13) (actual time=20.432..20.433 rows=1 loops=3)",
		"              Buckets: 1024  Batches: 1  Memory Usage: 9kB",
		"              Buffers: shared hit=1490 read=745",
		"              ->  Seq Scan on users  (cost=0.00..1564.20 rows=1 width=13) (actual time=17.609..20.416 rows=1 loops=3)",
		"                    Filter: ((email)::text = 'test31776@example.com'::text)",
		"                    Rows Removed by Filter: 65535",
		"                    Buffers: shared hit=1490 read=745",
		"Planning:",
		"  Buffers: shared hit=178 read=9",
		"Planning Time: 1.652 ms",
		"JIT:",
		"  Functions: 39",
		"  Options: Inlining false, Optimization false, Expressions true, Deforming true",
		"  Timing: Generation 2.460 ms, Inlining 0.000 ms, Optimization 4.860 ms, Emission 38.619 ms, Total 45.939 ms",
		"Execution Time: 796.036 ms",
	}
	for _, ex := range explains {
		rows.AddRow(ex)
	}
	mock.ExpectQuery("EXPLAIN \\(ANALYZE, BUFFERS\\) SELECT").WillReturnRows(rows)
}

const psqlTableSchemaQuery = "SELECT\\s+pg_class.relname AS table_name"

func (m *postgresMock) mockTableSchemaQuery(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{"table_name", "column_name", "is_pk"})
	rows.AddRow("todos", "description", false)
	rows.AddRow("todos", "user_id", false)
	rows.AddRow("users", "email", false)
	rows.AddRow("users", "id", true)
	rows.AddRow("users", "name", false)
	mock.ExpectQuery(psqlTableSchemaQuery).WillReturnRows(rows)
}

func (m *postgresMock) mockSelectQuery(mock sqlmock.Sqlmock, query string) {
	mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(0, 0))
}

func (m *postgresMock) mockCreateIndexQuery(mock sqlmock.Sqlmock) {
	mock.ExpectExec(`CREATE INDEX ".+_name" ON`).WillReturnResult(sqlmock.NewResult(0, 0))
}

func (m *postgresMock) mockDropIndexQuery(mock sqlmock.Sqlmock) {
	mock.ExpectExec(`DROP INDEX ".+_name"`).WillReturnResult(sqlmock.NewResult(0, 0))
}
