package pservice_test

import (
	"github.com/auxten/postgresql-parser/pkg/sql/parser"
	"github.com/mrasu/GravityR/database/postgres/pservice"
	"github.com/mrasu/GravityR/lib"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCollectTableNames(t *testing.T) {
	tests := []struct {
		name   string
		query  string
		tables []string
	}{
		{
			name:   "simple query",
			query:  "SELECT * FROM userS",
			tables: []string{"users"},
		},
		{
			name:   "simple aliased query",
			query:  "SELECT * FROM userS AS uuU",
			tables: []string{"users"},
		},
		{
			name:   "joined query",
			query:  "SELECT * FROM users AS u INNER JOIN tasks AS t ON users.id = t.user_id INNER JOIN labels AS l ON t.id = l.task_id",
			tables: []string{"users", "tasks", "labels"},
		},
		{
			name:   "subquery",
			query:  "SELECT * FROM users AS u INNER JOIN (SELECT * FROM tasks AS t) AS t2 ON u.id = t2.user_id",
			tables: []string{"users", "tasks"},
		},
		{
			name:   "CTE query",
			query:  "WITH t2 AS (SELECT * FROM tasks AS t) SELECT * FROM users AS u INNER JOIN t2 ON u.id = t2.user_id",
			tables: []string{"tasks", "users"},
		},
		{
			name:   "CTE name duplicated query",
			query:  "WITH users AS (SELECT * FROM tasks AS t) SELECT * FROM users AS u INNER JOIN users ON u.id = users.user_id",
			tables: []string{"tasks"},
		},
		{
			name:   "CTE name duplicated query",
			query:  "WITH users AS (SELECT * FROM tasks) SELECT * FROM users",
			tables: []string{"tasks"},
		},
		{
			name:   "union query",
			query:  "SELECT * FROM users UNION ALL SELECT * FROM tasks",
			tables: []string{"users", "tasks"},
		},
		{
			name:   "union query",
			query:  "SELECT * FROM users INNER JOIN (SELECT * FROM users UNION ALL SELECT * FROM tasks) AS u2 ON u2.id = users.id;",
			tables: []string{"users", "tasks"},
		},
		{
			name: "CTE name duplicated query",
			query: `
with count_tbl as
( select
  status,
  count(status) as count
  from users
  inner join tasks on users.user_id = tasks.user_id
  inner join tasks2 on users.user_id = tasks2.user_id
  where users.created_at > CURRENT_TIMESTAMP - interval '2 month'
  group by status
),
haha as (
 select users.name, status from users inner join (select status, user_id from tasks) as a on users.user_id = a.user_id
 where email like 'test%'
)
select
  case when status = 1 then 'Pending'
       when status = 2 then 'Doing'
       when status = 3 then 'Done' end as status,
  count / ( select sum(count) from count_tbl ) * 100 as percent
from  count_tbl
where status > 0 and count > 0
order by status desc`,
			tables: []string{"users", "tasks", "tasks2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmtNodes, err := parser.Parse(tt.query)
			assert.NoError(t, err)
			actualTables, errs := pservice.CollectTableNames(&stmtNodes[0])
			assert.Empty(t, errs)

			assert.ElementsMatch(t, tt.tables, actualTables)
		})
	}
}

func TestCollectTableNames_Errors(t *testing.T) {
	tests := []struct {
		name  string
		query string
		errs  []error
	}{
		{
			name:  "insert query",
			query: "INSERT INTO users(name) VALUES('hello')",
			errs:  []error{&lib.UnsupportedError{}},
		},
		{
			name:  "update query",
			query: "UPDATE users SET name = 'hello'",
			errs:  []error{&lib.UnsupportedError{}},
		},
		{
			name:  "reference different schema",
			query: "SELECT * FROM schema1.users",
			errs:  []error{&lib.UnsupportedError{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmtNodes, err := parser.Parse(tt.query)
			assert.NoError(t, err)
			actualTables, errs := pservice.CollectTableNames(&stmtNodes[0])
			assert.Nil(t, actualTables)

			assert.Equal(t, len(tt.errs), len(errs))
			for i, err := range errs {
				exErr := tt.errs[i]
				assert.ErrorAs(t, err, &exErr)
			}
		})
	}
}