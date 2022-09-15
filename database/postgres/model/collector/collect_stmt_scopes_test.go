package collector_test

import (
	"github.com/auxten/postgresql-parser/pkg/sql/parser"
	"github.com/google/go-cmp/cmp"
	"github.com/mrasu/GravityR/database/common_model"
	"github.com/mrasu/GravityR/database/postgres/model/collector"
	"github.com/mrasu/GravityR/thelper/tdata"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCollectStmtScopes_SingleField(t *testing.T) {
	tests := []struct {
		name   string
		query  string
		fields []*common_model.Field
		table  *common_model.Table
	}{
		{
			name:   "value query",
			query:  "SELECT 1",
			fields: nil,
			table:  nil,
		},
		{
			name:   "column reference query",
			query:  "SELECT name FROM users AS u",
			fields: []*common_model.Field{{Columns: []*common_model.FieldColumn{{Name: "name", Type: common_model.FieldReference}}}},
			table:  &common_model.Table{AsName: "u", Name: "users"},
		},
		{
			name:   "column reference query with table name",
			query:  "SELECT u.name FROM users AS u",
			fields: []*common_model.Field{{Columns: []*common_model.FieldColumn{{Table: "u", Name: "name", Type: common_model.FieldReference}}}},
			table:  &common_model.Table{AsName: "u", Name: "users"},
		},
		{
			name:   "star query",
			query:  "SELECT * FROM users",
			fields: []*common_model.Field{{Columns: []*common_model.FieldColumn{{Type: common_model.FieldStar}}}},
			table:  &common_model.Table{Name: "users"},
		},
		{
			name:   "star query with table name",
			query:  "SELECT users.* FROM users",
			fields: []*common_model.Field{{Columns: []*common_model.FieldColumn{{Table: "users", Type: common_model.FieldStar}}}},
			table:  &common_model.Table{Name: "users"},
		},
		{
			name:   "aliased query",
			query:  "SELECT name AS name1 FROM users",
			fields: []*common_model.Field{{AsName: "name1", Columns: []*common_model.FieldColumn{{Name: "name", Type: common_model.FieldReference}}}},
			table:  &common_model.Table{Name: "users"},
		},
		{
			name:   "between query",
			query:  "SELECT lunch_at BETWEEN start_at AND end_at FROM users",
			fields: []*common_model.Field{{Columns: []*common_model.FieldColumn{{Name: "lunch_at", Type: common_model.FieldReference}, {Name: "start_at", Type: common_model.FieldReference}, {Name: "end_at", Type: common_model.FieldReference}}}},
			table:  &common_model.Table{Name: "users"},
		},
		{
			name:  "operator query",
			query: "SELECT end_at - start_at FROM users",
			fields: []*common_model.Field{{Columns: []*common_model.FieldColumn{
				{Name: "end_at", Type: common_model.FieldReference},
				{Name: "start_at", Type: common_model.FieldReference},
			}}},
			table: &common_model.Table{Name: "users"},
		},
		{
			name:   "case query",
			query:  "SELECT CASE WHEN name = last_name THEN full_name ELSE first_name END FROM users",
			fields: []*common_model.Field{{Columns: []*common_model.FieldColumn{{Name: "name", Type: common_model.FieldReference}, {Name: "last_name", Type: common_model.FieldReference}, {Name: "full_name", Type: common_model.FieldReference}, {Name: "first_name", Type: common_model.FieldReference}}}},
			table:  &common_model.Table{Name: "users"},
		},
		{
			name:   "in query",
			query:  "SELECT user_id in (1, 2, 3, parent_id) FROM users",
			fields: []*common_model.Field{{Columns: []*common_model.FieldColumn{{Name: "user_id", Type: common_model.FieldReference}, {Name: "parent_id", Type: common_model.FieldReference}}}},
			table:  &common_model.Table{Name: "users"},
		},
		{
			name:   "is null query",
			query:  "SELECT user_id IS NULL FROM users",
			fields: []*common_model.Field{{Columns: []*common_model.FieldColumn{{Name: "user_id", Type: common_model.FieldReference}}}},
			table:  &common_model.Table{Name: "users"},
		},
		{
			name:   "is bool query",
			query:  "SELECT is_valid IS TRUE FROM users",
			fields: []*common_model.Field{{Columns: []*common_model.FieldColumn{{Name: "is_valid", Type: common_model.FieldReference}}}},
			table:  &common_model.Table{Name: "users"},
		},
		{
			name:   "like query",
			query:  "SELECT name LIKE valid_prefix FROM users",
			fields: []*common_model.Field{{Columns: []*common_model.FieldColumn{{Name: "name", Type: common_model.FieldReference}, {Name: "valid_prefix", Type: common_model.FieldReference}}}},
			table:  &common_model.Table{Name: "users"},
		},
		{
			name:   "parenthesis query",
			query:  "SELECT count1 + (count2 * count3) FROM users",
			fields: []*common_model.Field{{Columns: []*common_model.FieldColumn{{Name: "count1", Type: common_model.FieldReference}, {Name: "count2", Type: common_model.FieldReference}, {Name: "count3", Type: common_model.FieldReference}}}},
			table:  &common_model.Table{Name: "users"},
		},
		{
			name:   "similar to query",
			query:  "SELECT name SIMILAR TO valid_pattern FROM users",
			fields: []*common_model.Field{{Columns: []*common_model.FieldColumn{{Name: "name", Type: common_model.FieldReference}, {Name: "valid_pattern", Type: common_model.FieldReference}}}},
			table:  &common_model.Table{Name: "users"},
		},
		{
			name:   "unary operator query",
			query:  "SELECT -user_id FROM users",
			fields: []*common_model.Field{{Columns: []*common_model.FieldColumn{{Name: "user_id", Type: common_model.FieldReference}}}},
			table:  &common_model.Table{Name: "users"},
		},
		{
			name:   "trim query",
			query:  "SELECT TRIM(name) FROM users",
			fields: []*common_model.Field{{Columns: []*common_model.FieldColumn{{Name: "name", Type: common_model.FieldReference}}}},
			table:  &common_model.Table{Name: "users"},
		},
		{
			name:   "trim query with direction and column reference",
			query:  "SELECT TRIM(LEADING ' ' FROM name) FROM users",
			fields: []*common_model.Field{{Columns: []*common_model.FieldColumn{{Name: "name", Type: common_model.FieldReference}}}},
			table:  &common_model.Table{Name: "users"},
		},
		{
			name:   "cast query",
			query:  "SELECT CAST(user_id AS CHAR) FROM users",
			fields: []*common_model.Field{{Columns: []*common_model.FieldColumn{{Name: "user_id", Type: common_model.FieldReference}}}},
			table:  &common_model.Table{Name: "users"},
		},
		{
			name:   "aggregation query",
			query:  "SELECT SUM(id) FROM users",
			fields: []*common_model.Field{{Columns: []*common_model.FieldColumn{{Name: "id", Type: common_model.FieldReference}}}},
			table:  &common_model.Table{Name: "users"},
		},
		{
			name:   "aggregation query with star",
			query:  "SELECT COUNT(*) FROM users",
			fields: nil,
			table:  &common_model.Table{Name: "users"},
		},
		{
			name:   "window function query",
			query:  "SELECT AVG(todo_id) OVER(PARTITION BY user_id ORDER BY created_at ROWS BETWEEN 1 PRECEDING AND 2 FOLLOWING) FROM todos",
			fields: []*common_model.Field{{Columns: []*common_model.FieldColumn{{Name: "todo_id", Type: common_model.FieldReference}, {Name: "user_id", Type: common_model.FieldReference}, {Name: "created_at", Type: common_model.FieldReference}}}},
			table:  &common_model.Table{Name: "todos"},
		},
		{
			name:   "prepared statement query",
			query:  "SELECT $1",
			fields: nil,
			table:  nil,
		},
		{
			name:   "array query",
			query:  "SELECT ARRAY[1, id] FROM users",
			fields: []*common_model.Field{{Columns: []*common_model.FieldColumn{{Name: "id", Type: common_model.FieldReference}}}},
			table:  &common_model.Table{Name: "users"},
		},
		{
			name:   "nullif query",
			query:  "SELECT NULLIF(id, 1) FROM users",
			fields: []*common_model.Field{{Columns: []*common_model.FieldColumn{{Name: "id", Type: common_model.FieldReference}}}},
			table:  &common_model.Table{Name: "users"},
		},
		{
			name:   "coalesce query",
			query:  "SELECT COALESCE(id, 1) FROM users",
			fields: []*common_model.Field{{Columns: []*common_model.FieldColumn{{Name: "id", Type: common_model.FieldReference}}}},
			table:  &common_model.Table{Name: "users"},
		},
		{
			name:   "collate query",
			query:  "SELECT name COLLATE \"C\" FROM users",
			fields: []*common_model.Field{{Columns: []*common_model.FieldColumn{{Name: "name", Type: common_model.FieldReference}}}},
			table:  &common_model.Table{Name: "users"},
		},
		{
			name:   "tuple star query",
			query:  "SELECT (users).* FROM users",
			fields: []*common_model.Field{{Columns: []*common_model.FieldColumn{{Table: "users", Type: common_model.FieldStar}}}},
			table:  &common_model.Table{Name: "users"},
		},
		{
			name:   "tuple star query",
			query:  "SELECT (users).id FROM users",
			fields: []*common_model.Field{{Columns: []*common_model.FieldColumn{{Table: "users", Name: "id", Type: common_model.FieldReference}}}},
			table:  &common_model.Table{Name: "users"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := toStatement(t, tt.query)
			actualScopes, errs := collector.CollectStmtScopes(stmt)
			require.Empty(t, errs)

			expectedScopes := []*common_model.StmtScope{{
				Name:   "<root>",
				Fields: tt.fields,
			}}
			if tt.table != nil {
				expectedScopes[0].Tables = []*common_model.Table{tt.table}
			}
			if diff := cmp.Diff(expectedScopes, actualScopes); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestCollectStmtScopes_Joins(t *testing.T) {
	tests := []struct {
		name   string
		query  string
		fields []*common_model.Field
		tables []*common_model.Table
	}{
		{
			name:   "no join query",
			query:  "SELECT u.name FROM users AS u",
			fields: []*common_model.Field{{Columns: []*common_model.FieldColumn{{Table: "u", Name: "name", Type: common_model.FieldReference}}}},
			tables: []*common_model.Table{{AsName: "u", Name: "users"}},
		},
		{
			name: "single join query",
			query: `
				SELECT todos.id
				FROM users AS u
				INNER JOIN todos ON u.id = todos.user_id
			`,
			fields: []*common_model.Field{
				{Columns: []*common_model.FieldColumn{{Table: "todos", Name: "id", Type: common_model.FieldReference}}},
				{Columns: []*common_model.FieldColumn{{Table: "u", Name: "id", Type: common_model.FieldCondition}, {Table: "todos", Name: "user_id", Type: common_model.FieldCondition}}},
			},
			tables: []*common_model.Table{{AsName: "u", Name: "users"}, {Name: "todos"}},
		},
		{
			name: "multiple join query",
			query: `
				SELECT todos.id
				FROM users AS u
				INNER JOIN todos ON u.id = todos.user_id
				INNER JOIN todos AS t2 ON u.id = t2.user_id
				INNER JOIN todos AS t3 ON u.id = t3.user_id
			`,
			fields: []*common_model.Field{
				{Columns: []*common_model.FieldColumn{{Table: "todos", Name: "id", Type: common_model.FieldReference}}},
				{Columns: []*common_model.FieldColumn{{Table: "u", Name: "id", Type: common_model.FieldCondition}, {Table: "todos", Name: "user_id", Type: common_model.FieldCondition}}},
				{Columns: []*common_model.FieldColumn{{Table: "u", Name: "id", Type: common_model.FieldCondition}, {Table: "t2", Name: "user_id", Type: common_model.FieldCondition}}},
				{Columns: []*common_model.FieldColumn{{Table: "u", Name: "id", Type: common_model.FieldCondition}, {Table: "t3", Name: "user_id", Type: common_model.FieldCondition}}},
			},
			tables: []*common_model.Table{{AsName: "u", Name: "users"}, {Name: "todos"}, {AsName: "t2", Name: "todos"}, {AsName: "t3", Name: "todos"}},
		},
		{
			name: "row join query",
			query: `
				SELECT todos.id
				FROM users AS u
				INNER JOIN todos ON ROW(u.id, u.id2) = ROW(todos.user_id, todos.user_id2)
			`,
			fields: []*common_model.Field{
				{Columns: []*common_model.FieldColumn{{Table: "todos", Name: "id", Type: common_model.FieldReference}}},
				{Columns: []*common_model.FieldColumn{
					{Table: "u", Name: "id", Type: common_model.FieldCondition},
					{Table: "u", Name: "id2", Type: common_model.FieldCondition},
					{Table: "todos", Name: "user_id", Type: common_model.FieldCondition},
					{Table: "todos", Name: "user_id2", Type: common_model.FieldCondition},
				}},
			},
			tables: []*common_model.Table{{AsName: "u", Name: "users"}, {Name: "todos"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := toStatement(t, tt.query)
			actualScopes, errs := collector.CollectStmtScopes(stmt)
			require.Empty(t, errs)

			expectedScopes := []*common_model.StmtScope{{
				Name:   "<root>",
				Fields: tt.fields,
				Tables: tt.tables,
			}}
			if diff := cmp.Diff(expectedScopes, actualScopes); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestCollectStmtScopes_Subquery(t *testing.T) {
	tests := []struct {
		name   string
		query  string
		scopes []*common_model.StmtScope
	}{
		{
			name:  "subquery in SELECT",
			query: "SELECT (SELECT id FROM todos limit 1)",
			scopes: []*common_model.StmtScope{
				{
					Name: "<root>",
					Fields: []*common_model.Field{
						{Columns: []*common_model.FieldColumn{{ReferenceName: "<field0>", Type: common_model.FieldSubquery}}},
					},
					FieldScopes: []*common_model.StmtScope{
						{
							Name: "<field0>",
							Fields: []*common_model.Field{
								{Columns: []*common_model.FieldColumn{{Name: "id", Type: common_model.FieldReference}}},
							},
							Tables: []*common_model.Table{
								{Name: "todos"},
							},
						},
					},
				},
			},
		},
		{
			name:  "subquery in SELECT's aliased field",
			query: "SELECT id, (SELECT COUNT(status) FROM todos) AS status_count FROM users",
			scopes: []*common_model.StmtScope{{
				Name: "<root>",
				Fields: []*common_model.Field{
					{Columns: []*common_model.FieldColumn{{Name: "id", Type: common_model.FieldReference}}},
					{AsName: "status_count", Columns: []*common_model.FieldColumn{{ReferenceName: "<field0>", Type: common_model.FieldSubquery}}},
				},
				FieldScopes: []*common_model.StmtScope{
					{
						Name: "<field0>",
						Fields: []*common_model.Field{
							{Columns: []*common_model.FieldColumn{{Name: "status", Type: common_model.FieldReference}}},
						},
						Tables: []*common_model.Table{
							{Name: "todos"},
						},
					},
				},
				Tables: []*common_model.Table{
					{Name: "users"},
				},
			}},
		},
		{
			name:  "subquery with comparison",
			query: "SELECT (SELECT COUNT(status) FROM todos) - (SELECT COUNT(description) FROM todos) AS no_desc_count",
			scopes: []*common_model.StmtScope{{
				Name: "<root>",
				Fields: []*common_model.Field{
					{AsName: "no_desc_count", Columns: []*common_model.FieldColumn{
						{ReferenceName: "<field0>", Type: common_model.FieldSubquery},
						{ReferenceName: "<field1>", Type: common_model.FieldSubquery},
					}},
				},
				FieldScopes: []*common_model.StmtScope{
					{
						Name: "<field0>",
						Fields: []*common_model.Field{
							{Columns: []*common_model.FieldColumn{{Name: "status", Type: common_model.FieldReference}}},
						},
						Tables: []*common_model.Table{
							{Name: "todos"},
						},
					},
					{
						Name: "<field1>",
						Fields: []*common_model.Field{
							{Columns: []*common_model.FieldColumn{{Name: "description", Type: common_model.FieldReference}}},
						},
						Tables: []*common_model.Table{
							{Name: "todos"},
						},
					},
				},
			}},
		},
		{
			name:  "subquery with EXIST",
			query: "SELECT id, EXISTS (SELECT * FROM todos WHERE user_id = users.id) FROM users",
			scopes: []*common_model.StmtScope{{
				Name: "<root>",
				Fields: []*common_model.Field{
					{Columns: []*common_model.FieldColumn{{Name: "id", Type: common_model.FieldReference}}},
					{Columns: []*common_model.FieldColumn{{ReferenceName: "<field0>", Type: common_model.FieldSubquery}}},
				},
				FieldScopes: []*common_model.StmtScope{
					{
						Name: "<field0>",
						Fields: []*common_model.Field{
							{Columns: []*common_model.FieldColumn{{Type: common_model.FieldStar}}},
							{Columns: []*common_model.FieldColumn{
								{Name: "user_id", Type: common_model.FieldCondition},
								{Table: "users", Name: "id", Type: common_model.FieldCondition},
							}},
						},
						Tables: []*common_model.Table{
							{Name: "todos"},
						},
					},
				},
				Tables: []*common_model.Table{
					{Name: "users"},
				},
			}},
		},
		{
			name:  "subquery with IN",
			query: "SELECT user_id IN (SELECT user_id FROM todos) FROM users",
			scopes: []*common_model.StmtScope{{
				Name: "<root>",
				Fields: []*common_model.Field{
					{Columns: []*common_model.FieldColumn{
						{Name: "user_id", Type: common_model.FieldReference},
						{ReferenceName: "<field0>", Type: common_model.FieldSubquery},
					}},
				},
				FieldScopes: []*common_model.StmtScope{
					{
						Name: "<field0>",
						Fields: []*common_model.Field{
							{Columns: []*common_model.FieldColumn{{Name: "user_id", Type: common_model.FieldReference}}},
						},
						Tables: []*common_model.Table{
							{Name: "todos"},
						},
					},
				},
				Tables: []*common_model.Table{
					{Name: "users"},
				},
			}},
		},
		{
			name: "subquery in SELECT's function",
			query: `
SELECT row_to_json(
  (
    SELECT t
    FROM (SELECT id, description FROM todos LIMIT 1) AS t
  )
)
`,
			scopes: []*common_model.StmtScope{
				{
					Name: "<root>",
					Fields: []*common_model.Field{
						{Columns: []*common_model.FieldColumn{{ReferenceName: "<field0>", Type: common_model.FieldSubquery}}},
					},
					FieldScopes: []*common_model.StmtScope{
						{
							Name: "<field0>",
							Fields: []*common_model.Field{
								{Columns: []*common_model.FieldColumn{{Name: "t", Type: common_model.FieldReference}}},
							},
							Tables: []*common_model.Table{
								{Name: "<select0>", AsName: "t"},
							},
							Scopes: []*common_model.StmtScope{
								{
									Name: "<select0>",
									Fields: []*common_model.Field{
										{Columns: []*common_model.FieldColumn{{Name: "id", Type: common_model.FieldReference}}},
										{Columns: []*common_model.FieldColumn{{Name: "description", Type: common_model.FieldReference}}},
									},
									Tables: []*common_model.Table{
										{Name: "todos"},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "subquery in FROM",
			query: `
SELECT *
FROM (SELECT id, user_id FROM todos)
`,
			scopes: []*common_model.StmtScope{
				{
					Name: "<root>",
					Fields: []*common_model.Field{
						{Columns: []*common_model.FieldColumn{{Type: common_model.FieldStar}}},
					},
					Tables: []*common_model.Table{
						{Name: "<select0>"},
					},
					Scopes: []*common_model.StmtScope{
						{
							Name: "<select0>",
							Fields: []*common_model.Field{
								{Columns: []*common_model.FieldColumn{{Name: "id", Type: common_model.FieldReference}}},
								{Columns: []*common_model.FieldColumn{{Name: "user_id", Type: common_model.FieldReference}}},
							},
							Tables: []*common_model.Table{
								{Name: "todos"},
							},
						},
					},
				},
			},
		},
		{
			name: "subquery in JOIN",
			query: `
SELECT *
FROM users
INNER JOIN (SELECT id, user_id FROM todos) AS t ON users.id = t.user_id
`,
			scopes: []*common_model.StmtScope{
				{
					Name: "<root>",
					Fields: []*common_model.Field{
						{Columns: []*common_model.FieldColumn{{Type: common_model.FieldStar}}},
						{Columns: []*common_model.FieldColumn{
							{Table: "users", Name: "id", Type: common_model.FieldCondition},
							{Table: "t", Name: "user_id", Type: common_model.FieldCondition},
						}},
					},
					Tables: []*common_model.Table{
						{Name: "users"},
						{Name: "<select0>", AsName: "t"},
					},
					Scopes: []*common_model.StmtScope{
						{
							Name: "<select0>",
							Fields: []*common_model.Field{
								{Columns: []*common_model.FieldColumn{{Name: "id", Type: common_model.FieldReference}}},
								{Columns: []*common_model.FieldColumn{{Name: "user_id", Type: common_model.FieldReference}}},
							},
							Tables: []*common_model.Table{
								{Name: "todos"},
							},
						},
					},
				},
			},
		},
		{
			name: "subquery using star",
			query: `
SELECT description
FROM
	(SELECT * FROM todos) AS t
`,
			scopes: []*common_model.StmtScope{
				{
					Name: "<root>",
					Fields: []*common_model.Field{
						{Columns: []*common_model.FieldColumn{
							{Name: "description", Type: common_model.FieldReference},
						}},
					},
					Tables: []*common_model.Table{
						{Name: "<select0>", AsName: "t"},
					},
					Scopes: []*common_model.StmtScope{
						{
							Name: "<select0>",
							Fields: []*common_model.Field{
								{Columns: []*common_model.FieldColumn{{Type: common_model.FieldStar}}},
							},
							Tables: []*common_model.Table{
								{Name: "todos"},
							},
						},
					},
				},
			},
		},
		{
			name: "subquery with LATERAL JOIN ",
			query: `
SELECT description
FROM
	(SELECT id FROM users) AS u
	LEFT OUTER JOIN LATERAL (SELECT description FROM todos WHERE user_id = u.id) AS t ON u.id = t.user_id
`,
			scopes: []*common_model.StmtScope{
				{
					Name: "<root>",
					Fields: []*common_model.Field{
						{Columns: []*common_model.FieldColumn{
							{Name: "description", Type: common_model.FieldReference},
						}},
						{Columns: []*common_model.FieldColumn{
							{Table: "u", Name: "id", Type: common_model.FieldCondition},
							{Table: "t", Name: "user_id", Type: common_model.FieldCondition},
						}},
					},
					Tables: []*common_model.Table{
						{Name: "<select0>", AsName: "u"},
						{Name: "<select1>", AsName: "t", IsLateral: true},
					},
					Scopes: []*common_model.StmtScope{
						{
							Name: "<select0>",
							Fields: []*common_model.Field{
								{Columns: []*common_model.FieldColumn{{Name: "id", Type: common_model.FieldReference}}},
							},
							Tables: []*common_model.Table{
								{Name: "users"},
							},
						},
						{
							Name: "<select1>",
							Fields: []*common_model.Field{
								{Columns: []*common_model.FieldColumn{{Name: "description", Type: common_model.FieldReference}}},
								{Columns: []*common_model.FieldColumn{
									{Name: "user_id", Type: common_model.FieldCondition},
									{Table: "u", Name: "id", Type: common_model.FieldCondition},
								}},
							},
							Tables: []*common_model.Table{
								{Name: "todos"},
							},
						},
					},
				},
			},
		},
		{
			name: "subquery with aggregate function",
			query: `
SELECT json_agg(t)
FROM (SELECT id FROM todos LIMIT 3) AS t;
`,
			scopes: []*common_model.StmtScope{
				{
					Name: "<root>",
					Fields: []*common_model.Field{
						{Columns: []*common_model.FieldColumn{{Name: "t", Type: common_model.FieldReference}}},
					},
					Tables: []*common_model.Table{
						{AsName: "t", Name: "<select0>"},
					},
					Scopes: []*common_model.StmtScope{
						{
							Name: "<select0>",
							Fields: []*common_model.Field{
								{Columns: []*common_model.FieldColumn{{Name: "id", Type: common_model.FieldReference}}},
							},
							Tables: []*common_model.Table{
								{Name: "todos"},
							},
						},
					},
				},
			},
		},
		{
			name:   "hasura query",
			query:  tdata.Hasura1.SQL,
			scopes: tdata.Hasura1.Scopes,
		},
		{
			name:   "hasura query2",
			query:  tdata.Hasura2.SQL,
			scopes: tdata.Hasura2.Scopes,
		},
		{
			name:   "hasura query3",
			query:  tdata.Hasura3.SQL,
			scopes: tdata.Hasura3.Scopes,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := toStatement(t, tt.query)
			actualScopes, errs := collector.CollectStmtScopes(stmt)
			require.Empty(t, errs)

			if diff := cmp.Diff(tt.scopes, actualScopes); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestCollectStmtScopes_Where(t *testing.T) {
	tests := []struct {
		name   string
		query  string
		fields []*common_model.Field
		table  *common_model.Table
	}{
		{
			name:  "simple query(column at left)",
			query: "SELECT foo FROM users WHERE name = 'hello'",
			fields: []*common_model.Field{
				{Columns: []*common_model.FieldColumn{{Name: "foo", Type: common_model.FieldReference}}},
				{Columns: []*common_model.FieldColumn{{Name: "name", Type: common_model.FieldCondition}}},
			},
			table: &common_model.Table{Name: "users"},
		},
		{
			name:  "simple query(column at right",
			query: "SELECT foo FROM users WHERE 'hello' = name",
			fields: []*common_model.Field{
				{Columns: []*common_model.FieldColumn{{Name: "foo", Type: common_model.FieldReference}}},
				{Columns: []*common_model.FieldColumn{{Name: "name", Type: common_model.FieldCondition}}},
			},
			table: &common_model.Table{Name: "users"},
		},
		{
			name:  "same name with selected column",
			query: "SELECT name FROM users WHERE name = 'hello'",
			fields: []*common_model.Field{
				{Columns: []*common_model.FieldColumn{{Name: "name", Type: common_model.FieldReference}}},
				{Columns: []*common_model.FieldColumn{{Name: "name", Type: common_model.FieldCondition}}},
			},
			table: &common_model.Table{Name: "users"},
		},
		{
			name:  "compare with column",
			query: "SELECT foo FROM users WHERE name = first_name",
			fields: []*common_model.Field{
				{Columns: []*common_model.FieldColumn{{Name: "foo", Type: common_model.FieldReference}}},
				{Columns: []*common_model.FieldColumn{
					{Name: "name", Type: common_model.FieldCondition},
					{Name: "first_name", Type: common_model.FieldCondition},
				}},
			},
			table: &common_model.Table{Name: "users"},
		},
		{
			name:  "multiple comparison",
			query: "SELECT foo FROM users WHERE first_name = 'hello' AND middle_name = 'nice' AND last_name = 'world'",
			fields: []*common_model.Field{
				{Columns: []*common_model.FieldColumn{{Name: "foo", Type: common_model.FieldReference}}},
				{Columns: []*common_model.FieldColumn{
					{Name: "first_name", Type: common_model.FieldCondition},
					{Name: "middle_name", Type: common_model.FieldCondition},
					{Name: "last_name", Type: common_model.FieldCondition},
				}},
			},
			table: &common_model.Table{Name: "users"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := toStatement(t, tt.query)
			actualScopes, errs := collector.CollectStmtScopes(stmt)
			require.Empty(t, errs)

			expectedScopes := []*common_model.StmtScope{{
				Name:   "<root>",
				Fields: tt.fields,
				Tables: []*common_model.Table{tt.table},
			}}
			if diff := cmp.Diff(expectedScopes, actualScopes); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestCollectStmtScopes_GroupBy(t *testing.T) {
	tests := []struct {
		name   string
		query  string
		fields []*common_model.Field
		tables []*common_model.Table
	}{
		{
			name:  "simple query",
			query: "SELECT COUNT(*) FROM users GROUP BY name",
			fields: []*common_model.Field{
				{Columns: []*common_model.FieldColumn{{Name: "name", Type: common_model.FieldReference}}},
			},
			tables: []*common_model.Table{{Name: "users"}},
		},
		{
			name:   "group by value",
			query:  "SELECT COUNT(*) FROM users GROUP BY 1",
			fields: nil,
			tables: []*common_model.Table{{Name: "users"}},
		},
		{
			name:   "group by null",
			query:  "SELECT COUNT(*) FROM users GROUP BY null",
			fields: nil,
			tables: []*common_model.Table{{Name: "users"}},
		},
		{
			name:  "multiple columns",
			query: "SELECT COUNT(*) FROM users GROUP BY first_name, middle_name, last_name",
			fields: []*common_model.Field{
				{Columns: []*common_model.FieldColumn{{Name: "first_name", Type: common_model.FieldReference}}},
				{Columns: []*common_model.FieldColumn{{Name: "middle_name", Type: common_model.FieldReference}}},
				{Columns: []*common_model.FieldColumn{{Name: "last_name", Type: common_model.FieldReference}}},
			},
			tables: []*common_model.Table{{Name: "users"}},
		},
		{
			name:  "rollup",
			query: "SELECT COUNT(*) FROM users GROUP BY ROLLUP(name)",
			fields: []*common_model.Field{
				{Columns: []*common_model.FieldColumn{{Name: "name", Type: common_model.FieldReference}}},
			},
			tables: []*common_model.Table{{Name: "users"}},
		},
		{
			name:  "cube",
			query: "SELECT COUNT(*) FROM users GROUP BY CUBE(name)",
			fields: []*common_model.Field{
				{Columns: []*common_model.FieldColumn{{Name: "name", Type: common_model.FieldReference}}},
			},
			tables: []*common_model.Table{{Name: "users"}},
		},
		/*
			Not supporting GROUPING SETS because cockroachdb doesn't support it
			https://github.com/cockroachdb/cockroach/issues/46280
			{
				name:  "grouping sets",
				query: "SELECT COUNT(*) FROM users GROUP BY GROUPING SETS(name)",
				fields: []*common_model.Field{
					{Columns: []*common_model.FieldColumn{{Name: "name", Type: common_model.FieldReference}}},
				},
				tables: []*common_model.Table{{Name: "users"}},
			},
		*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := toStatement(t, tt.query)
			actualScopes, errs := collector.CollectStmtScopes(stmt)
			require.Empty(t, errs)

			expectedScopes := []*common_model.StmtScope{{
				Name:   "<root>",
				Fields: tt.fields,
				Tables: tt.tables,
			}}
			if diff := cmp.Diff(expectedScopes, actualScopes); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestCollectStmtScopes_Having(t *testing.T) {
	tests := []struct {
		name   string
		query  string
		fields []*common_model.Field
		tables []*common_model.Table
	}{
		{
			name:  "simple query",
			query: "SELECT count(*) FROM users HAVING COUNT(name) > 1",
			fields: []*common_model.Field{
				{Columns: []*common_model.FieldColumn{{Name: "name", Type: common_model.FieldReference}}},
			},
			tables: []*common_model.Table{{Name: "users"}},
		},
		{
			name:  "comparison",
			query: "SELECT count(*) FROM users HAVING COUNT(first_name) > COUNT(last_name)",
			fields: []*common_model.Field{
				{Columns: []*common_model.FieldColumn{
					{Name: "first_name", Type: common_model.FieldReference},
					{Name: "last_name", Type: common_model.FieldReference},
				}},
			},
			tables: []*common_model.Table{{Name: "users"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := toStatement(t, tt.query)
			actualScopes, errs := collector.CollectStmtScopes(stmt)
			require.Empty(t, errs)

			expectedScopes := []*common_model.StmtScope{{
				Name:   "<root>",
				Fields: tt.fields,
				Tables: tt.tables,
			}}
			if diff := cmp.Diff(expectedScopes, actualScopes); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestCollectStmtScopes_OrderBy(t *testing.T) {
	tests := []struct {
		name   string
		query  string
		fields []*common_model.Field
		tables []*common_model.Table
	}{
		{
			name:  "simple query",
			query: "SELECT 1 FROM users ORDER BY name",
			fields: []*common_model.Field{
				{Columns: []*common_model.FieldColumn{{Name: "name", Type: common_model.FieldReference}}},
			},
			tables: []*common_model.Table{{Name: "users"}},
		},
		{
			name:   "order by value",
			query:  "SELECT COUNT(*) FROM users ORDER BY 1",
			fields: nil,
			tables: []*common_model.Table{{Name: "users"}},
		},
		{
			name:   "group by null",
			query:  "SELECT COUNT(*) FROM users ORDER BY null",
			fields: nil,
			tables: []*common_model.Table{{Name: "users"}},
		},
		{
			name:  "multiple columns",
			query: "SELECT COUNT(*) FROM users ORDER BY first_name, middle_name, last_name",
			fields: []*common_model.Field{
				{Columns: []*common_model.FieldColumn{{Name: "first_name", Type: common_model.FieldReference}}},
				{Columns: []*common_model.FieldColumn{{Name: "middle_name", Type: common_model.FieldReference}}},
				{Columns: []*common_model.FieldColumn{{Name: "last_name", Type: common_model.FieldReference}}},
			},
			tables: []*common_model.Table{{Name: "users"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := toStatement(t, tt.query)
			actualScopes, errs := collector.CollectStmtScopes(stmt)
			require.Empty(t, errs)

			expectedScopes := []*common_model.StmtScope{{
				Name:   "<root>",
				Fields: tt.fields,
				Tables: tt.tables,
			}}
			if diff := cmp.Diff(expectedScopes, actualScopes); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestCollectStmtScopes_Limit(t *testing.T) {
	tests := []struct {
		name   string
		query  string
		fields []*common_model.Field
		tables []*common_model.Table
	}{
		{
			name:   "simple query",
			query:  "SELECT 1 FROM users LIMIT 1",
			fields: nil,
			tables: []*common_model.Table{{Name: "users"}},
		},
		{
			name:   "offset exists",
			query:  "SELECT 1 FROM users LIMIT 1 OFFSET 1",
			fields: nil,
			tables: []*common_model.Table{{Name: "users"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := toStatement(t, tt.query)
			actualScopes, errs := collector.CollectStmtScopes(stmt)
			require.Empty(t, errs)

			expectedScopes := []*common_model.StmtScope{{
				Name:   "<root>",
				Fields: tt.fields,
				Tables: tt.tables,
			}}
			if diff := cmp.Diff(expectedScopes, actualScopes); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func toStatement(t *testing.T, query string) *parser.Statement {
	t.Helper()
	stmts, err := parser.Parse(query)
	require.NoError(t, err)
	require.Equal(t, 1, len(stmts))

	return &stmts[0]
}
