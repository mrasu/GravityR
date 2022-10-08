package parser_test

import (
	"github.com/google/go-cmp/cmp"
	"github.com/mrasu/GravityR/database/dservice/dparser"
	"github.com/mrasu/GravityR/database/mysql/mservice/parser"
	tParser "github.com/pingcap/tidb/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCollectStmtScopes_SingleField(t *testing.T) {
	tests := []struct {
		name   string
		query  string
		fields []*dparser.Field
		table  *dparser.Table
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
			fields: []*dparser.Field{{Columns: []*dparser.FieldColumn{{Name: "name", Type: dparser.FieldReference}}}},
			table:  &dparser.Table{AsName: "u", Name: "users"},
		},
		{
			name:   "column reference query with table name",
			query:  "SELECT u.name FROM users AS u",
			fields: []*dparser.Field{{Columns: []*dparser.FieldColumn{{Table: "u", Name: "name", Type: dparser.FieldReference}}}},
			table:  &dparser.Table{AsName: "u", Name: "users"},
		},
		{
			name:   "star query",
			query:  "SELECT * FROM users",
			fields: []*dparser.Field{{Columns: []*dparser.FieldColumn{{Type: dparser.FieldStar}}}},
			table:  &dparser.Table{Name: "users"},
		},
		{
			name:   "star query with table name",
			query:  "SELECT users.* FROM users",
			fields: []*dparser.Field{{Columns: []*dparser.FieldColumn{{Table: "users", Type: dparser.FieldStar}}}},
			table:  &dparser.Table{Name: "users"},
		},
		{
			name:   "aliased query",
			query:  "SELECT name AS name1 FROM users",
			fields: []*dparser.Field{{AsName: "name1", Columns: []*dparser.FieldColumn{{Name: "name", Type: dparser.FieldReference}}}},
			table:  &dparser.Table{Name: "users"},
		},
		{
			name:   "between query",
			query:  "SELECT lunch_at BETWEEN start_at AND end_at FROM users",
			fields: []*dparser.Field{{Columns: []*dparser.FieldColumn{{Name: "lunch_at", Type: dparser.FieldReference}, {Name: "start_at", Type: dparser.FieldReference}, {Name: "end_at", Type: dparser.FieldReference}}}},
			table:  &dparser.Table{Name: "users"},
		},
		{
			name:  "operator query",
			query: "SELECT end_at - start_at FROM users",
			fields: []*dparser.Field{{Columns: []*dparser.FieldColumn{
				{Name: "end_at", Type: dparser.FieldReference},
				{Name: "start_at", Type: dparser.FieldReference},
			}}},
			table: &dparser.Table{Name: "users"},
		},
		{
			name:   "case query",
			query:  "SELECT CASE WHEN name = last_name THEN full_name ELSE first_name END FROM users",
			fields: []*dparser.Field{{Columns: []*dparser.FieldColumn{{Name: "name", Type: dparser.FieldReference}, {Name: "last_name", Type: dparser.FieldReference}, {Name: "full_name", Type: dparser.FieldReference}, {Name: "first_name", Type: dparser.FieldReference}}}},
			table:  &dparser.Table{Name: "users"},
		},
		{
			name:   "in query",
			query:  "SELECT user_id in (1, 2, 3, parent_id) FROM users",
			fields: []*dparser.Field{{Columns: []*dparser.FieldColumn{{Name: "user_id", Type: dparser.FieldReference}, {Name: "parent_id", Type: dparser.FieldReference}}}},
			table:  &dparser.Table{Name: "users"},
		},
		{
			name:   "is null query",
			query:  "SELECT user_id IS NULL FROM users",
			fields: []*dparser.Field{{Columns: []*dparser.FieldColumn{{Name: "user_id", Type: dparser.FieldReference}}}},
			table:  &dparser.Table{Name: "users"},
		},
		{
			name:   "is bool query",
			query:  "SELECT is_valid IS TRUE FROM users",
			fields: []*dparser.Field{{Columns: []*dparser.FieldColumn{{Name: "is_valid", Type: dparser.FieldReference}}}},
			table:  &dparser.Table{Name: "users"},
		},
		{
			name:   "like query",
			query:  "SELECT name LIKE valid_prefix FROM users",
			fields: []*dparser.Field{{Columns: []*dparser.FieldColumn{{Name: "name", Type: dparser.FieldReference}, {Name: "valid_prefix", Type: dparser.FieldReference}}}},
			table:  &dparser.Table{Name: "users"},
		},
		{
			name:   "parenthesis query",
			query:  "SELECT count1 + (count2 * count3) FROM users",
			fields: []*dparser.Field{{Columns: []*dparser.FieldColumn{{Name: "count1", Type: dparser.FieldReference}, {Name: "count2", Type: dparser.FieldReference}, {Name: "count3", Type: dparser.FieldReference}}}},
			table:  &dparser.Table{Name: "users"},
		},
		{
			name:   "regex query",
			query:  "SELECT name REGEXP valid_pattern FROM users",
			fields: []*dparser.Field{{Columns: []*dparser.FieldColumn{{Name: "name", Type: dparser.FieldReference}, {Name: "valid_pattern", Type: dparser.FieldReference}}}},
			table:  &dparser.Table{Name: "users"},
		},
		{
			name:   "unary operator query",
			query:  "SELECT -user_id FROM users",
			fields: []*dparser.Field{{Columns: []*dparser.FieldColumn{{Name: "user_id", Type: dparser.FieldReference}}}},
			table:  &dparser.Table{Name: "users"},
		},
		{
			name:   "variable query",
			query:  "SELECT @hello",
			fields: nil,
			table:  nil,
		},
		{
			name:  "match query",
			query: "SELECT MATCH(first_name, last_name) AGAINST(valid_pattern) FROM users",
			fields: []*dparser.Field{{Columns: []*dparser.FieldColumn{
				{Name: "valid_pattern", Type: dparser.FieldReference},
				{Name: "first_name", Type: dparser.FieldReference},
				{Name: "last_name", Type: dparser.FieldReference},
			}}},
			table: &dparser.Table{Name: "users"},
		},
		{
			name:   "trim query",
			query:  "SELECT TRIM(name) FROM users",
			fields: []*dparser.Field{{Columns: []*dparser.FieldColumn{{Name: "name", Type: dparser.FieldReference}}}},
			table:  &dparser.Table{Name: "users"},
		},
		{
			name:   "trim query with direction",
			query:  "SELECT TRIM(LEADING ' ' FROM ' hello world ')",
			fields: nil,
			table:  nil,
		},
		{
			name:   "trim query with direction and column reference",
			query:  "SELECT TRIM(LEADING ' ' FROM name) FROM users",
			fields: []*dparser.Field{{Columns: []*dparser.FieldColumn{{Name: "name", Type: dparser.FieldReference}}}},
			table:  &dparser.Table{Name: "users"},
		},
		{
			name:   "cast query",
			query:  "SELECT CAST(user_id AS CHAR) FROM users",
			fields: []*dparser.Field{{Columns: []*dparser.FieldColumn{{Name: "user_id", Type: dparser.FieldReference}}}},
			table:  &dparser.Table{Name: "users"},
		},
		{
			name:   "aggregation query",
			query:  "SELECT SUM(id) FROM users",
			fields: []*dparser.Field{{Columns: []*dparser.FieldColumn{{Name: "id", Type: dparser.FieldReference}}}},
			table:  &dparser.Table{Name: "users"},
		},
		{
			name:   "aggregation query with star",
			query:  "SELECT COUNT(*) FROM users",
			fields: nil,
			table:  &dparser.Table{Name: "users"},
		},
		{
			name:   "window function query",
			query:  "SELECT AVG(task_id) OVER(PARTITION BY user_id ORDER BY created_at ROWS BETWEEN 1 PRECEDING AND 2 FOLLOWING) FROM tasks",
			fields: []*dparser.Field{{Columns: []*dparser.FieldColumn{{Name: "task_id", Type: dparser.FieldReference}, {Name: "user_id", Type: dparser.FieldReference}, {Name: "created_at", Type: dparser.FieldReference}}}},
			table:  &dparser.Table{Name: "tasks"},
		},
		{
			name:   "GET_FORMAT function query",
			query:  "SELECT GET_FORMAT(DATE, 'USA')",
			fields: nil,
			table:  nil,
		},
		{
			name:   "GET_FORMAT function query",
			query:  "SELECT GET_FORMAT(DATE, name) FROM regions",
			fields: []*dparser.Field{{Columns: []*dparser.FieldColumn{{Name: "name", Type: dparser.FieldReference}}}},
			table:  &dparser.Table{Name: "regions"},
		},
		{
			name:   "prepared statement query",
			query:  "SELECT ?",
			fields: nil,
			table:  nil,
		},
	}
	p := tParser.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmtNodes, _, err := p.Parse(tt.query, "", "")
			require.NoError(t, err)
			actualScopes, errs := parser.CollectStmtScopes(stmtNodes[0])
			require.Empty(t, errs)

			expectedScopes := []*dparser.StmtScope{{
				Name:   "<root>",
				Fields: tt.fields,
			}}
			if tt.table != nil {
				expectedScopes[0].Tables = []*dparser.Table{tt.table}
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
		fields []*dparser.Field
		tables []*dparser.Table
	}{
		{
			name:   "no join query",
			query:  "SELECT u.name FROM users AS u",
			fields: []*dparser.Field{{Columns: []*dparser.FieldColumn{{Table: "u", Name: "name", Type: dparser.FieldReference}}}},
			tables: []*dparser.Table{{AsName: "u", Name: "users"}},
		},
		{
			name: "single join query",
			query: `
				SELECT tasks.id
				FROM users AS u
				INNER JOIN tasks ON u.id = tasks.user_id
			`,
			fields: []*dparser.Field{
				{Columns: []*dparser.FieldColumn{{Table: "tasks", Name: "id", Type: dparser.FieldReference}}},
				{Columns: []*dparser.FieldColumn{{Table: "u", Name: "id", Type: dparser.FieldCondition}, {Table: "tasks", Name: "user_id", Type: dparser.FieldCondition}}},
			},
			tables: []*dparser.Table{{AsName: "u", Name: "users"}, {Name: "tasks"}},
		},
		{
			name: "multiple join query",
			query: `
				SELECT tasks.id
				FROM users AS u
				INNER JOIN tasks ON u.id = tasks.user_id
				INNER JOIN tasks AS t2 ON u.id = t2.user_id
				INNER JOIN tasks AS t3 ON u.id = t3.user_id
			`,
			fields: []*dparser.Field{
				{Columns: []*dparser.FieldColumn{{Table: "tasks", Name: "id", Type: dparser.FieldReference}}},
				{Columns: []*dparser.FieldColumn{{Table: "u", Name: "id", Type: dparser.FieldCondition}, {Table: "tasks", Name: "user_id", Type: dparser.FieldCondition}}},
				{Columns: []*dparser.FieldColumn{{Table: "u", Name: "id", Type: dparser.FieldCondition}, {Table: "t2", Name: "user_id", Type: dparser.FieldCondition}}},
				{Columns: []*dparser.FieldColumn{{Table: "u", Name: "id", Type: dparser.FieldCondition}, {Table: "t3", Name: "user_id", Type: dparser.FieldCondition}}},
			},
			tables: []*dparser.Table{{AsName: "u", Name: "users"}, {Name: "tasks"}, {AsName: "t2", Name: "tasks"}, {AsName: "t3", Name: "tasks"}},
		},
		{
			name: "row join query",
			query: `
				SELECT tasks.id
				FROM users AS u
				INNER JOIN tasks ON ROW(u.id, u.id2) = ROW(tasks.user_id, tasks.user_id2)
			`,
			fields: []*dparser.Field{
				{Columns: []*dparser.FieldColumn{{Table: "tasks", Name: "id", Type: dparser.FieldReference}}},
				{Columns: []*dparser.FieldColumn{
					{Table: "u", Name: "id", Type: dparser.FieldCondition},
					{Table: "u", Name: "id2", Type: dparser.FieldCondition},
					{Table: "tasks", Name: "user_id", Type: dparser.FieldCondition},
					{Table: "tasks", Name: "user_id2", Type: dparser.FieldCondition},
				}},
			},
			tables: []*dparser.Table{{AsName: "u", Name: "users"}, {Name: "tasks"}},
		},
		{
			name: "subquery join query",
			query: `
				SELECT tasks.id
				FROM users AS u
				INNER JOIN (SELECT * FROM tasks) AS t ON u.id = t.user_id
			`,
			fields: []*dparser.Field{
				{Columns: []*dparser.FieldColumn{{Table: "tasks", Name: "id", Type: dparser.FieldReference}}},
				{Columns: []*dparser.FieldColumn{
					{Table: "u", Name: "id", Type: dparser.FieldCondition},
					{Table: "t", Name: "user_id", Type: dparser.FieldCondition},
				}},
			},
			tables: []*dparser.Table{{AsName: "u", Name: "users"}, {AsName: "t"}},
		},
	}
	p := tParser.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmtNodes, _, err := p.Parse(tt.query, "", "")
			require.NoError(t, err)
			actualScopes, errs := parser.CollectStmtScopes(stmtNodes[0])
			require.Empty(t, errs)

			expectedScopes := []*dparser.StmtScope{{
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

func TestCollectStmtScopes_NestedScope(t *testing.T) {
	tests := []struct {
		name   string
		query  string
		fields []*dparser.Field
		tables []*dparser.Table
	}{
		{
			name:   "subquery",
			query:  "SELECT id, (SELECT COUNT(status) FROM tasks) AS status_count FROM users",
			fields: []*dparser.Field{{Columns: []*dparser.FieldColumn{{Name: "", Type: dparser.FieldSubquery}}}},
			tables: []*dparser.Table{},
		},
		{
			name:   "subquery with comparison",
			query:  "SELECT (SELECT COUNT(status) FROM tasks) - (SELECT COUNT(description) FROM tasks) AS no_desc_count",
			fields: []*dparser.Field{{Columns: []*dparser.FieldColumn{{Name: "", Type: dparser.FieldSubquery}}}},
			tables: []*dparser.Table{},
		},
		{
			name:   "exists subquery",
			query:  "SELECT id, EXISTS (SELECT * FROM tasks WHERE user_id = users.user_id) FROM users",
			fields: []*dparser.Field{{Columns: []*dparser.FieldColumn{{Name: "", Type: dparser.FieldSubquery}}}},
			tables: []*dparser.Table{},
		},
		{
			name:   "in subquery",
			query:  "SELECT user_id IN (SELECT user_id FROM tasks) FROM users",
			fields: []*dparser.Field{{Columns: []*dparser.FieldColumn{{Name: "", Type: dparser.FieldSubquery}}}},
			tables: []*dparser.Table{},
		},
	}
	p := tParser.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Skipf("Nest not implemented yet")
			stmtNodes, _, err := p.Parse(tt.query, "", "")
			require.NoError(t, err)
			actualScopes, errs := parser.CollectStmtScopes(stmtNodes[0])
			require.Empty(t, errs)

			require.Equal(t, 1, len(actualScopes))
			assert.ElementsMatch(t, tt.fields, actualScopes[0].Fields)
			assert.ElementsMatch(t, tt.tables, actualScopes[0].Tables)
			expectedScopes := []*dparser.StmtScope{{
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

func TestCollectStmtScopes_Where(t *testing.T) {
	tests := []struct {
		name   string
		query  string
		fields []*dparser.Field
		table  *dparser.Table
	}{
		{
			name:  "simple query(column at left)",
			query: "SELECT foo FROM users WHERE name = 'hello'",
			fields: []*dparser.Field{
				{Columns: []*dparser.FieldColumn{{Name: "foo", Type: dparser.FieldReference}}},
				{Columns: []*dparser.FieldColumn{{Name: "name", Type: dparser.FieldCondition}}},
			},
			table: &dparser.Table{Name: "users"},
		},
		{
			name:  "simple query(column at right",
			query: "SELECT foo FROM users WHERE 'hello' = name",
			fields: []*dparser.Field{
				{Columns: []*dparser.FieldColumn{{Name: "foo", Type: dparser.FieldReference}}},
				{Columns: []*dparser.FieldColumn{{Name: "name", Type: dparser.FieldCondition}}},
			},
			table: &dparser.Table{Name: "users"},
		},
		{
			name:  "same name with selected column",
			query: "SELECT name FROM users WHERE name = 'hello'",
			fields: []*dparser.Field{
				{Columns: []*dparser.FieldColumn{{Name: "name", Type: dparser.FieldReference}}},
				{Columns: []*dparser.FieldColumn{{Name: "name", Type: dparser.FieldCondition}}},
			},
			table: &dparser.Table{Name: "users"},
		},
		{
			name:  "compare with column",
			query: "SELECT foo FROM users WHERE name = first_name",
			fields: []*dparser.Field{
				{Columns: []*dparser.FieldColumn{{Name: "foo", Type: dparser.FieldReference}}},
				{Columns: []*dparser.FieldColumn{
					{Name: "name", Type: dparser.FieldCondition},
					{Name: "first_name", Type: dparser.FieldCondition},
				}},
			},
			table: &dparser.Table{Name: "users"},
		},
		{
			name:  "multiple comparison",
			query: "SELECT foo FROM users WHERE first_name = 'hello' AND middle_name = 'nice' AND last_name = 'world'",
			fields: []*dparser.Field{
				{Columns: []*dparser.FieldColumn{{Name: "foo", Type: dparser.FieldReference}}},
				{Columns: []*dparser.FieldColumn{
					{Name: "first_name", Type: dparser.FieldCondition},
					{Name: "middle_name", Type: dparser.FieldCondition},
					{Name: "last_name", Type: dparser.FieldCondition},
				}},
			},
			table: &dparser.Table{Name: "users"},
		},
	}
	p := tParser.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmtNodes, _, err := p.Parse(tt.query, "", "")
			require.NoError(t, err)
			actualScopes, errs := parser.CollectStmtScopes(stmtNodes[0])
			require.Empty(t, errs)

			expectedScopes := []*dparser.StmtScope{{
				Name:   "<root>",
				Fields: tt.fields,
				Tables: []*dparser.Table{tt.table},
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
		fields []*dparser.Field
		tables []*dparser.Table
	}{
		{
			name:  "simple query",
			query: "SELECT COUNT(*) FROM users GROUP BY name",
			fields: []*dparser.Field{
				{Columns: []*dparser.FieldColumn{{Name: "name", Type: dparser.FieldReference}}},
			},
			tables: []*dparser.Table{{Name: "users"}},
		},
		{
			name:   "group by value",
			query:  "SELECT COUNT(*) FROM users GROUP BY 1",
			fields: nil,
			tables: []*dparser.Table{{Name: "users"}},
		},
		{
			name:   "group by null",
			query:  "SELECT COUNT(*) FROM users GROUP BY null",
			fields: nil,
			tables: []*dparser.Table{{Name: "users"}},
		},
		{
			name:  "multiple columns",
			query: "SELECT COUNT(*) FROM users GROUP BY first_name, middle_name, last_name",
			fields: []*dparser.Field{
				{Columns: []*dparser.FieldColumn{{Name: "first_name", Type: dparser.FieldReference}}},
				{Columns: []*dparser.FieldColumn{{Name: "middle_name", Type: dparser.FieldReference}}},
				{Columns: []*dparser.FieldColumn{{Name: "last_name", Type: dparser.FieldReference}}},
			},
			tables: []*dparser.Table{{Name: "users"}},
		},
		/*
			Not supporting WITH ROLLUP because pingcap/tidb doesn't support it.
			https://github.com/pingcap/tidb/issues/4250
			{
				name:  "with rollup",
				query: "SELECT COUNT(*) FROM users GROUP BY name WITH ROLLUP",
				fields: []*dparser.Field{
					{Columns: []*dparser.FieldColumn{{Name: "name", Type: dparser.FieldReference}}},
				},
				tables: []*dparser.Table{{Name: "users"}},
			},
		*/
	}
	p := tParser.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmtNodes, _, err := p.Parse(tt.query, "", "")
			require.NoError(t, err)
			actualScopes, errs := parser.CollectStmtScopes(stmtNodes[0])
			require.Empty(t, errs)

			expectedScopes := []*dparser.StmtScope{{
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
		fields []*dparser.Field
		tables []*dparser.Table
	}{
		{
			name:  "simple query",
			query: "SELECT count(*) FROM users HAVING COUNT(name) > 1",
			fields: []*dparser.Field{
				{Columns: []*dparser.FieldColumn{{Name: "name", Type: dparser.FieldReference}}},
			},
			tables: []*dparser.Table{{Name: "users"}},
		},
		{
			name:  "comparison",
			query: "SELECT count(*) FROM users HAVING COUNT(first_name) > COUNT(last_name)",
			fields: []*dparser.Field{
				{Columns: []*dparser.FieldColumn{
					{Name: "first_name", Type: dparser.FieldReference},
					{Name: "last_name", Type: dparser.FieldReference},
				}},
			},
			tables: []*dparser.Table{{Name: "users"}},
		},
	}
	p := tParser.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmtNodes, _, err := p.Parse(tt.query, "", "")
			require.NoError(t, err)
			actualScopes, errs := parser.CollectStmtScopes(stmtNodes[0])
			require.Empty(t, errs)

			expectedScopes := []*dparser.StmtScope{{
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
		fields []*dparser.Field
		tables []*dparser.Table
	}{
		{
			name:  "simple query",
			query: "SELECT 1 FROM users ORDER BY name",
			fields: []*dparser.Field{
				{Columns: []*dparser.FieldColumn{{Name: "name", Type: dparser.FieldReference}}},
			},
			tables: []*dparser.Table{{Name: "users"}},
		},
		{
			name:   "order by value",
			query:  "SELECT COUNT(*) FROM users ORDER BY 1",
			fields: nil,
			tables: []*dparser.Table{{Name: "users"}},
		},
		{
			name:   "group by null",
			query:  "SELECT COUNT(*) FROM users ORDER BY null",
			fields: nil,
			tables: []*dparser.Table{{Name: "users"}},
		},
		{
			name:  "multiple columns",
			query: "SELECT COUNT(*) FROM users ORDER BY first_name, middle_name, last_name",
			fields: []*dparser.Field{
				{Columns: []*dparser.FieldColumn{{Name: "first_name", Type: dparser.FieldReference}}},
				{Columns: []*dparser.FieldColumn{{Name: "middle_name", Type: dparser.FieldReference}}},
				{Columns: []*dparser.FieldColumn{{Name: "last_name", Type: dparser.FieldReference}}},
			},
			tables: []*dparser.Table{{Name: "users"}},
		},
	}
	p := tParser.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmtNodes, _, err := p.Parse(tt.query, "", "")
			require.NoError(t, err)
			actualScopes, errs := parser.CollectStmtScopes(stmtNodes[0])
			require.Empty(t, errs)

			expectedScopes := []*dparser.StmtScope{{
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
		fields []*dparser.Field
		tables []*dparser.Table
	}{
		{
			name:   "simple query",
			query:  "SELECT 1 FROM users LIMIT 1",
			fields: nil,
			tables: []*dparser.Table{{Name: "users"}},
		},
		{
			name:   "offset exists",
			query:  "SELECT 1 FROM users LIMIT 1 OFFSET 1",
			fields: nil,
			tables: []*dparser.Table{{Name: "users"}},
		},
		{
			name:   "offset without OFFSET keyword",
			query:  "SELECT 1 FROM users LIMIT 1, 1",
			fields: nil,
			tables: []*dparser.Table{{Name: "users"}},
		},
	}
	p := tParser.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmtNodes, _, err := p.Parse(tt.query, "", "")
			require.NoError(t, err)
			actualScopes, errs := parser.CollectStmtScopes(stmtNodes[0])
			require.Empty(t, errs)

			expectedScopes := []*dparser.StmtScope{{
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
