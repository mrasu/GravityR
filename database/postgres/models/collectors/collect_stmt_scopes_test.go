package collectors_test

import (
	"github.com/auxten/postgresql-parser/pkg/sql/parser"
	"github.com/mrasu/GravityR/database/db_models"
	"github.com/mrasu/GravityR/database/postgres/models/collectors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCollectStmtScopes_SingleField(t *testing.T) {
	tests := []struct {
		name   string
		query  string
		fields []*db_models.Field
		table  *db_models.Table
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
			fields: []*db_models.Field{{Columns: []*db_models.FieldColumn{{Name: "name", Type: db_models.FieldReference}}}},
			table:  &db_models.Table{AsName: "u", Name: "users"},
		},
		{
			name:   "column reference query with table name",
			query:  "SELECT u.name FROM users AS u",
			fields: []*db_models.Field{{Columns: []*db_models.FieldColumn{{Table: "u", Name: "name", Type: db_models.FieldReference}}}},
			table:  &db_models.Table{AsName: "u", Name: "users"},
		},
		{
			name:   "star query",
			query:  "SELECT * FROM users",
			fields: []*db_models.Field{{Columns: []*db_models.FieldColumn{{Type: db_models.FieldStar}}}},
			table:  &db_models.Table{Name: "users"},
		},
		{
			name:   "star query with table name",
			query:  "SELECT users.* FROM users",
			fields: []*db_models.Field{{Columns: []*db_models.FieldColumn{{Table: "users", Type: db_models.FieldStar}}}},
			table:  &db_models.Table{Name: "users"},
		},
		{
			name:   "aliased query",
			query:  "SELECT name AS name1 FROM users",
			fields: []*db_models.Field{{AsName: "name1", Columns: []*db_models.FieldColumn{{Name: "name", Type: db_models.FieldReference}}}},
			table:  &db_models.Table{Name: "users"},
		},
		{
			name:   "between query",
			query:  "SELECT lunch_at BETWEEN start_at AND end_at FROM users",
			fields: []*db_models.Field{{Columns: []*db_models.FieldColumn{{Name: "lunch_at", Type: db_models.FieldReference}, {Name: "start_at", Type: db_models.FieldReference}, {Name: "end_at", Type: db_models.FieldReference}}}},
			table:  &db_models.Table{Name: "users"},
		},
		{
			name:  "operator query",
			query: "SELECT end_at - start_at FROM users",
			fields: []*db_models.Field{{Columns: []*db_models.FieldColumn{
				{Name: "end_at", Type: db_models.FieldReference},
				{Name: "start_at", Type: db_models.FieldReference},
			}}},
			table: &db_models.Table{Name: "users"},
		},
		{
			name:   "case query",
			query:  "SELECT CASE WHEN name = last_name THEN full_name ELSE first_name END FROM users",
			fields: []*db_models.Field{{Columns: []*db_models.FieldColumn{{Name: "name", Type: db_models.FieldReference}, {Name: "last_name", Type: db_models.FieldReference}, {Name: "full_name", Type: db_models.FieldReference}, {Name: "first_name", Type: db_models.FieldReference}}}},
			table:  &db_models.Table{Name: "users"},
		},
		{
			name:   "in query",
			query:  "SELECT user_id in (1, 2, 3, parent_id) FROM users",
			fields: []*db_models.Field{{Columns: []*db_models.FieldColumn{{Name: "user_id", Type: db_models.FieldReference}, {Name: "parent_id", Type: db_models.FieldReference}}}},
			table:  &db_models.Table{Name: "users"},
		},
		{
			name:   "is null query",
			query:  "SELECT user_id IS NULL FROM users",
			fields: []*db_models.Field{{Columns: []*db_models.FieldColumn{{Name: "user_id", Type: db_models.FieldReference}}}},
			table:  &db_models.Table{Name: "users"},
		},
		{
			name:   "is bool query",
			query:  "SELECT is_valid IS TRUE FROM users",
			fields: []*db_models.Field{{Columns: []*db_models.FieldColumn{{Name: "is_valid", Type: db_models.FieldReference}}}},
			table:  &db_models.Table{Name: "users"},
		},
		{
			name:   "like query",
			query:  "SELECT name LIKE valid_prefix FROM users",
			fields: []*db_models.Field{{Columns: []*db_models.FieldColumn{{Name: "name", Type: db_models.FieldReference}, {Name: "valid_prefix", Type: db_models.FieldReference}}}},
			table:  &db_models.Table{Name: "users"},
		},
		{
			name:   "parenthesis query",
			query:  "SELECT count1 + (count2 * count3) FROM users",
			fields: []*db_models.Field{{Columns: []*db_models.FieldColumn{{Name: "count1", Type: db_models.FieldReference}, {Name: "count2", Type: db_models.FieldReference}, {Name: "count3", Type: db_models.FieldReference}}}},
			table:  &db_models.Table{Name: "users"},
		},
		{
			name:   "similar to query",
			query:  "SELECT name SIMILAR TO valid_pattern FROM users",
			fields: []*db_models.Field{{Columns: []*db_models.FieldColumn{{Name: "name", Type: db_models.FieldReference}, {Name: "valid_pattern", Type: db_models.FieldReference}}}},
			table:  &db_models.Table{Name: "users"},
		},
		{
			name:   "unary operator query",
			query:  "SELECT -user_id FROM users",
			fields: []*db_models.Field{{Columns: []*db_models.FieldColumn{{Name: "user_id", Type: db_models.FieldReference}}}},
			table:  &db_models.Table{Name: "users"},
		},
		{
			name:   "trim query",
			query:  "SELECT TRIM(name) FROM users",
			fields: []*db_models.Field{{Columns: []*db_models.FieldColumn{{Name: "name", Type: db_models.FieldReference}}}},
			table:  &db_models.Table{Name: "users"},
		},
		{
			name:   "trim query with direction and column reference",
			query:  "SELECT TRIM(LEADING ' ' FROM name) FROM users",
			fields: []*db_models.Field{{Columns: []*db_models.FieldColumn{{Name: "name", Type: db_models.FieldReference}}}},
			table:  &db_models.Table{Name: "users"},
		},
		{
			name:   "cast query",
			query:  "SELECT CAST(user_id AS CHAR) FROM users",
			fields: []*db_models.Field{{Columns: []*db_models.FieldColumn{{Name: "user_id", Type: db_models.FieldReference}}}},
			table:  &db_models.Table{Name: "users"},
		},
		{
			name:   "aggregation query",
			query:  "SELECT SUM(id) FROM users",
			fields: []*db_models.Field{{Columns: []*db_models.FieldColumn{{Name: "id", Type: db_models.FieldReference}}}},
			table:  &db_models.Table{Name: "users"},
		},
		{
			name:   "aggregation query with star",
			query:  "SELECT COUNT(*) FROM users",
			fields: nil,
			table:  &db_models.Table{Name: "users"},
		},
		{
			name:   "window function query",
			query:  "SELECT AVG(todo_id) OVER(PARTITION BY user_id ORDER BY created_at ROWS BETWEEN 1 PRECEDING AND 2 FOLLOWING) FROM todos",
			fields: []*db_models.Field{{Columns: []*db_models.FieldColumn{{Name: "todo_id", Type: db_models.FieldReference}, {Name: "user_id", Type: db_models.FieldReference}, {Name: "created_at", Type: db_models.FieldReference}}}},
			table:  &db_models.Table{Name: "todos"},
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
			fields: []*db_models.Field{{Columns: []*db_models.FieldColumn{{Name: "id", Type: db_models.FieldReference}}}},
			table:  &db_models.Table{Name: "users"},
		},
		{
			name:   "nullif query",
			query:  "SELECT NULLIF(id, 1) FROM users",
			fields: []*db_models.Field{{Columns: []*db_models.FieldColumn{{Name: "id", Type: db_models.FieldReference}}}},
			table:  &db_models.Table{Name: "users"},
		},
		{
			name:   "coalesce query",
			query:  "SELECT COALESCE(id, 1) FROM users",
			fields: []*db_models.Field{{Columns: []*db_models.FieldColumn{{Name: "id", Type: db_models.FieldReference}}}},
			table:  &db_models.Table{Name: "users"},
		},
		{
			name:   "collate query",
			query:  "SELECT name COLLATE \"C\" FROM users",
			fields: []*db_models.Field{{Columns: []*db_models.FieldColumn{{Name: "name", Type: db_models.FieldReference}}}},
			table:  &db_models.Table{Name: "users"},
		},
		{
			name:   "tuple star query",
			query:  "SELECT (users).* FROM users",
			fields: []*db_models.Field{{Columns: []*db_models.FieldColumn{{Table: "users", Type: db_models.FieldStar}}}},
			table:  &db_models.Table{Name: "users"},
		},
		{
			name:   "tuple star query",
			query:  "SELECT (users).id FROM users",
			fields: []*db_models.Field{{Columns: []*db_models.FieldColumn{{Table: "users", Name: "id", Type: db_models.FieldReference}}}},
			table:  &db_models.Table{Name: "users"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := toStatement(t, tt.query)
			actualScopes, errs := collectors.CollectStmtScopes(stmt)
			require.Empty(t, errs)

			require.Equal(t, 1, len(actualScopes))
			actualScope := actualScopes[0]

			assert.Equal(t, "<root>", actualScope.Name)
			assert.ElementsMatch(t, tt.fields, actualScope.Fields)

			if len(actualScope.Tables) == 0 {
				assert.Equal(t, 0, len(actualScope.Tables))
			} else {
				assert.Equal(t, 1, len(actualScope.Tables))
				assert.Equal(t, tt.table, actualScope.Tables[0])
			}
		})
	}
}

func TestCollectStmtScopes_Joins(t *testing.T) {
	tests := []struct {
		name   string
		query  string
		fields []*db_models.Field
		tables []*db_models.Table
	}{
		{
			name:   "no join query",
			query:  "SELECT u.name FROM users AS u",
			fields: []*db_models.Field{{Columns: []*db_models.FieldColumn{{Table: "u", Name: "name", Type: db_models.FieldReference}}}},
			tables: []*db_models.Table{{AsName: "u", Name: "users"}},
		},
		{
			name: "single join query",
			query: `
				SELECT todos.id
				FROM users AS u
				INNER JOIN todos ON u.id = todos.user_id
			`,
			fields: []*db_models.Field{
				{Columns: []*db_models.FieldColumn{{Table: "todos", Name: "id", Type: db_models.FieldReference}}},
				{Columns: []*db_models.FieldColumn{{Table: "u", Name: "id", Type: db_models.FieldCondition}, {Table: "todos", Name: "user_id", Type: db_models.FieldCondition}}},
			},
			tables: []*db_models.Table{{AsName: "u", Name: "users"}, {Name: "todos"}},
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
			fields: []*db_models.Field{
				{Columns: []*db_models.FieldColumn{{Table: "todos", Name: "id", Type: db_models.FieldReference}}},
				{Columns: []*db_models.FieldColumn{{Table: "u", Name: "id", Type: db_models.FieldCondition}, {Table: "todos", Name: "user_id", Type: db_models.FieldCondition}}},
				{Columns: []*db_models.FieldColumn{{Table: "u", Name: "id", Type: db_models.FieldCondition}, {Table: "t2", Name: "user_id", Type: db_models.FieldCondition}}},
				{Columns: []*db_models.FieldColumn{{Table: "u", Name: "id", Type: db_models.FieldCondition}, {Table: "t3", Name: "user_id", Type: db_models.FieldCondition}}},
			},
			tables: []*db_models.Table{{AsName: "u", Name: "users"}, {Name: "todos"}, {AsName: "t2", Name: "todos"}, {AsName: "t3", Name: "todos"}},
		},
		{
			name: "row join query",
			query: `
				SELECT todos.id
				FROM users AS u
				INNER JOIN todos ON ROW(u.id, u.id2) = ROW(todos.user_id, todos.user_id2)
			`,
			fields: []*db_models.Field{
				{Columns: []*db_models.FieldColumn{{Table: "todos", Name: "id", Type: db_models.FieldReference}}},
				{Columns: []*db_models.FieldColumn{
					{Table: "u", Name: "id", Type: db_models.FieldCondition},
					{Table: "u", Name: "id2", Type: db_models.FieldCondition},
					{Table: "todos", Name: "user_id", Type: db_models.FieldCondition},
					{Table: "todos", Name: "user_id2", Type: db_models.FieldCondition},
				}},
			},
			tables: []*db_models.Table{{AsName: "u", Name: "users"}, {Name: "todos"}},
		},
		{
			name: "subquery join query",
			query: `
				SELECT todos.id
				FROM users AS u
				INNER JOIN (SELECT * FROM todos) AS t ON u.id = t.user_id
			`,
			fields: []*db_models.Field{
				{Columns: []*db_models.FieldColumn{{Table: "todos", Name: "id", Type: db_models.FieldReference}}},
				{Columns: []*db_models.FieldColumn{
					{Table: "u", Name: "id", Type: db_models.FieldCondition},
					{Table: "t", Name: "user_id", Type: db_models.FieldCondition},
				}},
			},
			tables: []*db_models.Table{{AsName: "u", Name: "users"}, {AsName: "t"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := toStatement(t, tt.query)
			actualScopes, errs := collectors.CollectStmtScopes(stmt)
			require.Empty(t, errs)

			require.Equal(t, 1, len(actualScopes))
			actualScope := actualScopes[0]

			assert.Equal(t, "<root>", actualScope.Name)
			assert.ElementsMatch(t, tt.fields, actualScope.Fields)
			assert.ElementsMatch(t, tt.tables, actualScope.Tables)
		})
	}
}

func TestCollectStmtScopes_NestedScope(t *testing.T) {
	tests := []struct {
		name   string
		query  string
		fields []*db_models.Field
		tables []string
	}{
		{
			name:   "subquery",
			query:  "SELECT id, (SELECT COUNT(status) FROM todos) AS status_count FROM users",
			fields: []*db_models.Field{{Columns: []*db_models.FieldColumn{{Name: "", Type: db_models.FieldSubquery}}}},
			tables: []string{},
		},
		{
			name:   "subquery with comparison",
			query:  "SELECT (SELECT COUNT(status) FROM todos) - (SELECT COUNT(description) FROM todos) AS no_desc_count",
			fields: []*db_models.Field{{Columns: []*db_models.FieldColumn{{Name: "", Type: db_models.FieldSubquery}}}},
			tables: []string{},
		},
		{
			name:   "exists subquery",
			query:  "SELECT id, EXISTS (SELECT * FROM todos WHERE user_id = users.user_id) FROM users",
			fields: []*db_models.Field{{Columns: []*db_models.FieldColumn{{Name: "", Type: db_models.FieldSubquery}}}},
			tables: []string{},
		},
		{
			name:   "in subquery",
			query:  "SELECT user_id IN (SELECT user_id FROM todos) FROM users",
			fields: []*db_models.Field{{Columns: []*db_models.FieldColumn{{Name: "", Type: db_models.FieldSubquery}}}},
			tables: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Skipf("Nest not implemented yet")
			stmt := toStatement(t, tt.query)
			actualScopes, errs := collectors.CollectStmtScopes(stmt)
			require.Empty(t, errs)

			require.Equal(t, 1, len(actualScopes))
			assert.ElementsMatch(t, tt.fields, actualScopes[0].Fields)

			// TODO
			//assert.ElementsMatch(t, tt.tables, actualScopes[0].Tables)
		})
	}
}

func TestCollectStmtScopes_Where(t *testing.T) {
	tests := []struct {
		name   string
		query  string
		fields []*db_models.Field
		table  *db_models.Table
	}{
		{
			name:  "simple query(column at left)",
			query: "SELECT foo FROM users WHERE name = 'hello'",
			fields: []*db_models.Field{
				{Columns: []*db_models.FieldColumn{{Name: "foo", Type: db_models.FieldReference}}},
				{Columns: []*db_models.FieldColumn{{Name: "name", Type: db_models.FieldCondition}}},
			},
			table: &db_models.Table{Name: "users"},
		},
		{
			name:  "simple query(column at right",
			query: "SELECT foo FROM users WHERE 'hello' = name",
			fields: []*db_models.Field{
				{Columns: []*db_models.FieldColumn{{Name: "foo", Type: db_models.FieldReference}}},
				{Columns: []*db_models.FieldColumn{{Name: "name", Type: db_models.FieldCondition}}},
			},
			table: &db_models.Table{Name: "users"},
		},
		{
			name:  "same name with selected column",
			query: "SELECT name FROM users WHERE name = 'hello'",
			fields: []*db_models.Field{
				{Columns: []*db_models.FieldColumn{{Name: "name", Type: db_models.FieldReference}}},
				{Columns: []*db_models.FieldColumn{{Name: "name", Type: db_models.FieldCondition}}},
			},
			table: &db_models.Table{Name: "users"},
		},
		{
			name:  "compare with column",
			query: "SELECT foo FROM users WHERE name = first_name",
			fields: []*db_models.Field{
				{Columns: []*db_models.FieldColumn{{Name: "foo", Type: db_models.FieldReference}}},
				{Columns: []*db_models.FieldColumn{
					{Name: "name", Type: db_models.FieldCondition},
					{Name: "first_name", Type: db_models.FieldCondition},
				}},
			},
			table: &db_models.Table{Name: "users"},
		},
		{
			name:  "multiple comparison",
			query: "SELECT foo FROM users WHERE first_name = 'hello' AND middle_name = 'nice' AND last_name = 'world'",
			fields: []*db_models.Field{
				{Columns: []*db_models.FieldColumn{{Name: "foo", Type: db_models.FieldReference}}},
				{Columns: []*db_models.FieldColumn{
					{Name: "first_name", Type: db_models.FieldCondition},
					{Name: "middle_name", Type: db_models.FieldCondition},
					{Name: "last_name", Type: db_models.FieldCondition},
				}},
			},
			table: &db_models.Table{Name: "users"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := toStatement(t, tt.query)
			actualScopes, errs := collectors.CollectStmtScopes(stmt)
			require.Empty(t, errs)

			require.Equal(t, 1, len(actualScopes))
			actualScope := actualScopes[0]

			assert.Equal(t, "<root>", actualScope.Name)
			assert.ElementsMatch(t, tt.fields, actualScope.Fields)
			assert.ElementsMatch(t, []*db_models.Table{tt.table}, actualScope.Tables)
		})
	}
}

func TestCollectStmtScopes_GroupBy(t *testing.T) {
	tests := []struct {
		name   string
		query  string
		fields []*db_models.Field
		tables []*db_models.Table
	}{
		{
			name:  "simple query",
			query: "SELECT COUNT(*) FROM users GROUP BY name",
			fields: []*db_models.Field{
				{Columns: []*db_models.FieldColumn{{Name: "name", Type: db_models.FieldReference}}},
			},
			tables: []*db_models.Table{{Name: "users"}},
		},
		{
			name:   "group by value",
			query:  "SELECT COUNT(*) FROM users GROUP BY 1",
			fields: []*db_models.Field{},
			tables: []*db_models.Table{{Name: "users"}},
		},
		{
			name:   "group by null",
			query:  "SELECT COUNT(*) FROM users GROUP BY null",
			fields: []*db_models.Field{},
			tables: []*db_models.Table{{Name: "users"}},
		},
		{
			name:  "multiple columns",
			query: "SELECT COUNT(*) FROM users GROUP BY first_name, middle_name, last_name",
			fields: []*db_models.Field{
				{Columns: []*db_models.FieldColumn{{Name: "first_name", Type: db_models.FieldReference}}},
				{Columns: []*db_models.FieldColumn{{Name: "middle_name", Type: db_models.FieldReference}}},
				{Columns: []*db_models.FieldColumn{{Name: "last_name", Type: db_models.FieldReference}}},
			},
			tables: []*db_models.Table{{Name: "users"}},
		},
		{
			name:  "rollup",
			query: "SELECT COUNT(*) FROM users GROUP BY ROLLUP(name)",
			fields: []*db_models.Field{
				{Columns: []*db_models.FieldColumn{{Name: "name", Type: db_models.FieldReference}}},
			},
			tables: []*db_models.Table{{Name: "users"}},
		},
		{
			name:  "cube",
			query: "SELECT COUNT(*) FROM users GROUP BY CUBE(name)",
			fields: []*db_models.Field{
				{Columns: []*db_models.FieldColumn{{Name: "name", Type: db_models.FieldReference}}},
			},
			tables: []*db_models.Table{{Name: "users"}},
		},
		/*
			Not supporting GROUPING SETS because cockroachdb doesn't support it
			https://github.com/cockroachdb/cockroach/issues/46280
			{
				name:  "grouping sets",
				query: "SELECT COUNT(*) FROM users GROUP BY GROUPING SETS(name)",
				fields: []*db_models.Field{
					{Columns: []*db_models.FieldColumn{{Name: "name", Type: db_models.FieldReference}}},
				},
				tables: []*db_models.Table{{Name: "users"}},
			},
		*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := toStatement(t, tt.query)
			actualScopes, errs := collectors.CollectStmtScopes(stmt)
			require.Empty(t, errs)

			require.Equal(t, 1, len(actualScopes))
			actualScope := actualScopes[0]

			assert.Equal(t, "<root>", actualScope.Name)
			assert.ElementsMatch(t, tt.fields, actualScope.Fields)
			assert.ElementsMatch(t, tt.tables, actualScope.Tables)
		})
	}
}

func TestCollectStmtScopes_Having(t *testing.T) {
	tests := []struct {
		name   string
		query  string
		fields []*db_models.Field
		tables []*db_models.Table
	}{
		{
			name:  "simple query",
			query: "SELECT count(*) FROM users HAVING COUNT(name) > 1",
			fields: []*db_models.Field{
				{Columns: []*db_models.FieldColumn{{Name: "name", Type: db_models.FieldReference}}},
			},
			tables: []*db_models.Table{{Name: "users"}},
		},
		{
			name:  "comparison",
			query: "SELECT count(*) FROM users HAVING COUNT(first_name) > COUNT(last_name)",
			fields: []*db_models.Field{
				{Columns: []*db_models.FieldColumn{
					{Name: "first_name", Type: db_models.FieldReference},
					{Name: "last_name", Type: db_models.FieldReference},
				}},
			},
			tables: []*db_models.Table{{Name: "users"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := toStatement(t, tt.query)
			actualScopes, errs := collectors.CollectStmtScopes(stmt)
			require.Empty(t, errs)

			require.Equal(t, 1, len(actualScopes))
			actualScope := actualScopes[0]

			assert.Equal(t, "<root>", actualScope.Name)
			assert.ElementsMatch(t, tt.fields, actualScope.Fields)
			assert.ElementsMatch(t, tt.tables, actualScope.Tables)
		})
	}
}

func TestCollectStmtScopes_OrderBy(t *testing.T) {
	tests := []struct {
		name   string
		query  string
		fields []*db_models.Field
		tables []*db_models.Table
	}{
		{
			name:  "simple query",
			query: "SELECT 1 FROM users ORDER BY name",
			fields: []*db_models.Field{
				{Columns: []*db_models.FieldColumn{{Name: "name", Type: db_models.FieldReference}}},
			},
			tables: []*db_models.Table{{Name: "users"}},
		},
		{
			name:   "order by value",
			query:  "SELECT COUNT(*) FROM users ORDER BY 1",
			fields: []*db_models.Field{},
			tables: []*db_models.Table{{Name: "users"}},
		},
		{
			name:   "group by null",
			query:  "SELECT COUNT(*) FROM users ORDER BY null",
			fields: []*db_models.Field{},
			tables: []*db_models.Table{{Name: "users"}},
		},
		{
			name:  "multiple columns",
			query: "SELECT COUNT(*) FROM users ORDER BY first_name, middle_name, last_name",
			fields: []*db_models.Field{
				{Columns: []*db_models.FieldColumn{{Name: "first_name", Type: db_models.FieldReference}}},
				{Columns: []*db_models.FieldColumn{{Name: "middle_name", Type: db_models.FieldReference}}},
				{Columns: []*db_models.FieldColumn{{Name: "last_name", Type: db_models.FieldReference}}},
			},
			tables: []*db_models.Table{{Name: "users"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := toStatement(t, tt.query)
			actualScopes, errs := collectors.CollectStmtScopes(stmt)
			require.Empty(t, errs)

			require.Equal(t, 1, len(actualScopes))
			actualScope := actualScopes[0]

			assert.Equal(t, "<root>", actualScope.Name)
			assert.ElementsMatch(t, tt.fields, actualScope.Fields)
			assert.ElementsMatch(t, tt.tables, actualScope.Tables)
		})
	}
}

func TestCollectStmtScopes_Limit(t *testing.T) {
	tests := []struct {
		name   string
		query  string
		fields []*db_models.Field
		tables []*db_models.Table
	}{
		{
			name:   "simple query",
			query:  "SELECT 1 FROM users LIMIT 1",
			fields: []*db_models.Field{},
			tables: []*db_models.Table{{Name: "users"}},
		},
		{
			name:   "offset exists",
			query:  "SELECT 1 FROM users LIMIT 1 OFFSET 1",
			fields: []*db_models.Field{},
			tables: []*db_models.Table{{Name: "users"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := toStatement(t, tt.query)
			actualScopes, errs := collectors.CollectStmtScopes(stmt)
			require.Empty(t, errs)

			require.Equal(t, 1, len(actualScopes))
			actualScope := actualScopes[0]

			assert.Equal(t, "<root>", actualScope.Name)
			assert.ElementsMatch(t, tt.fields, actualScope.Fields)
			assert.ElementsMatch(t, tt.tables, actualScope.Tables)
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