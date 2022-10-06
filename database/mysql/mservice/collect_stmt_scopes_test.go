package mservice_test

import (
	"github.com/google/go-cmp/cmp"
	"github.com/mrasu/GravityR/database/dmodel"
	"github.com/mrasu/GravityR/database/mysql/mservice"
	"github.com/pingcap/tidb/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCollectStmtScopes_SingleField(t *testing.T) {
	tests := []struct {
		name   string
		query  string
		fields []*dmodel.Field
		table  *dmodel.Table
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
			fields: []*dmodel.Field{{Columns: []*dmodel.FieldColumn{{Name: "name", Type: dmodel.FieldReference}}}},
			table:  &dmodel.Table{AsName: "u", Name: "users"},
		},
		{
			name:   "column reference query with table name",
			query:  "SELECT u.name FROM users AS u",
			fields: []*dmodel.Field{{Columns: []*dmodel.FieldColumn{{Table: "u", Name: "name", Type: dmodel.FieldReference}}}},
			table:  &dmodel.Table{AsName: "u", Name: "users"},
		},
		{
			name:   "star query",
			query:  "SELECT * FROM users",
			fields: []*dmodel.Field{{Columns: []*dmodel.FieldColumn{{Type: dmodel.FieldStar}}}},
			table:  &dmodel.Table{Name: "users"},
		},
		{
			name:   "star query with table name",
			query:  "SELECT users.* FROM users",
			fields: []*dmodel.Field{{Columns: []*dmodel.FieldColumn{{Table: "users", Type: dmodel.FieldStar}}}},
			table:  &dmodel.Table{Name: "users"},
		},
		{
			name:   "aliased query",
			query:  "SELECT name AS name1 FROM users",
			fields: []*dmodel.Field{{AsName: "name1", Columns: []*dmodel.FieldColumn{{Name: "name", Type: dmodel.FieldReference}}}},
			table:  &dmodel.Table{Name: "users"},
		},
		{
			name:   "between query",
			query:  "SELECT lunch_at BETWEEN start_at AND end_at FROM users",
			fields: []*dmodel.Field{{Columns: []*dmodel.FieldColumn{{Name: "lunch_at", Type: dmodel.FieldReference}, {Name: "start_at", Type: dmodel.FieldReference}, {Name: "end_at", Type: dmodel.FieldReference}}}},
			table:  &dmodel.Table{Name: "users"},
		},
		{
			name:  "operator query",
			query: "SELECT end_at - start_at FROM users",
			fields: []*dmodel.Field{{Columns: []*dmodel.FieldColumn{
				{Name: "end_at", Type: dmodel.FieldReference},
				{Name: "start_at", Type: dmodel.FieldReference},
			}}},
			table: &dmodel.Table{Name: "users"},
		},
		{
			name:   "case query",
			query:  "SELECT CASE WHEN name = last_name THEN full_name ELSE first_name END FROM users",
			fields: []*dmodel.Field{{Columns: []*dmodel.FieldColumn{{Name: "name", Type: dmodel.FieldReference}, {Name: "last_name", Type: dmodel.FieldReference}, {Name: "full_name", Type: dmodel.FieldReference}, {Name: "first_name", Type: dmodel.FieldReference}}}},
			table:  &dmodel.Table{Name: "users"},
		},
		{
			name:   "in query",
			query:  "SELECT user_id in (1, 2, 3, parent_id) FROM users",
			fields: []*dmodel.Field{{Columns: []*dmodel.FieldColumn{{Name: "user_id", Type: dmodel.FieldReference}, {Name: "parent_id", Type: dmodel.FieldReference}}}},
			table:  &dmodel.Table{Name: "users"},
		},
		{
			name:   "is null query",
			query:  "SELECT user_id IS NULL FROM users",
			fields: []*dmodel.Field{{Columns: []*dmodel.FieldColumn{{Name: "user_id", Type: dmodel.FieldReference}}}},
			table:  &dmodel.Table{Name: "users"},
		},
		{
			name:   "is bool query",
			query:  "SELECT is_valid IS TRUE FROM users",
			fields: []*dmodel.Field{{Columns: []*dmodel.FieldColumn{{Name: "is_valid", Type: dmodel.FieldReference}}}},
			table:  &dmodel.Table{Name: "users"},
		},
		{
			name:   "like query",
			query:  "SELECT name LIKE valid_prefix FROM users",
			fields: []*dmodel.Field{{Columns: []*dmodel.FieldColumn{{Name: "name", Type: dmodel.FieldReference}, {Name: "valid_prefix", Type: dmodel.FieldReference}}}},
			table:  &dmodel.Table{Name: "users"},
		},
		{
			name:   "parenthesis query",
			query:  "SELECT count1 + (count2 * count3) FROM users",
			fields: []*dmodel.Field{{Columns: []*dmodel.FieldColumn{{Name: "count1", Type: dmodel.FieldReference}, {Name: "count2", Type: dmodel.FieldReference}, {Name: "count3", Type: dmodel.FieldReference}}}},
			table:  &dmodel.Table{Name: "users"},
		},
		{
			name:   "regex query",
			query:  "SELECT name REGEXP valid_pattern FROM users",
			fields: []*dmodel.Field{{Columns: []*dmodel.FieldColumn{{Name: "name", Type: dmodel.FieldReference}, {Name: "valid_pattern", Type: dmodel.FieldReference}}}},
			table:  &dmodel.Table{Name: "users"},
		},
		{
			name:   "unary operator query",
			query:  "SELECT -user_id FROM users",
			fields: []*dmodel.Field{{Columns: []*dmodel.FieldColumn{{Name: "user_id", Type: dmodel.FieldReference}}}},
			table:  &dmodel.Table{Name: "users"},
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
			fields: []*dmodel.Field{{Columns: []*dmodel.FieldColumn{
				{Name: "valid_pattern", Type: dmodel.FieldReference},
				{Name: "first_name", Type: dmodel.FieldReference},
				{Name: "last_name", Type: dmodel.FieldReference},
			}}},
			table: &dmodel.Table{Name: "users"},
		},
		{
			name:   "trim query",
			query:  "SELECT TRIM(name) FROM users",
			fields: []*dmodel.Field{{Columns: []*dmodel.FieldColumn{{Name: "name", Type: dmodel.FieldReference}}}},
			table:  &dmodel.Table{Name: "users"},
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
			fields: []*dmodel.Field{{Columns: []*dmodel.FieldColumn{{Name: "name", Type: dmodel.FieldReference}}}},
			table:  &dmodel.Table{Name: "users"},
		},
		{
			name:   "cast query",
			query:  "SELECT CAST(user_id AS CHAR) FROM users",
			fields: []*dmodel.Field{{Columns: []*dmodel.FieldColumn{{Name: "user_id", Type: dmodel.FieldReference}}}},
			table:  &dmodel.Table{Name: "users"},
		},
		{
			name:   "aggregation query",
			query:  "SELECT SUM(id) FROM users",
			fields: []*dmodel.Field{{Columns: []*dmodel.FieldColumn{{Name: "id", Type: dmodel.FieldReference}}}},
			table:  &dmodel.Table{Name: "users"},
		},
		{
			name:   "aggregation query with star",
			query:  "SELECT COUNT(*) FROM users",
			fields: nil,
			table:  &dmodel.Table{Name: "users"},
		},
		{
			name:   "window function query",
			query:  "SELECT AVG(task_id) OVER(PARTITION BY user_id ORDER BY created_at ROWS BETWEEN 1 PRECEDING AND 2 FOLLOWING) FROM tasks",
			fields: []*dmodel.Field{{Columns: []*dmodel.FieldColumn{{Name: "task_id", Type: dmodel.FieldReference}, {Name: "user_id", Type: dmodel.FieldReference}, {Name: "created_at", Type: dmodel.FieldReference}}}},
			table:  &dmodel.Table{Name: "tasks"},
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
			fields: []*dmodel.Field{{Columns: []*dmodel.FieldColumn{{Name: "name", Type: dmodel.FieldReference}}}},
			table:  &dmodel.Table{Name: "regions"},
		},
		{
			name:   "prepared statement query",
			query:  "SELECT ?",
			fields: nil,
			table:  nil,
		},
	}
	p := parser.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmtNodes, _, err := p.Parse(tt.query, "", "")
			require.NoError(t, err)
			actualScopes, errs := mservice.CollectStmtScopes(stmtNodes[0])
			require.Empty(t, errs)

			expectedScopes := []*dmodel.StmtScope{{
				Name:   "<root>",
				Fields: tt.fields,
			}}
			if tt.table != nil {
				expectedScopes[0].Tables = []*dmodel.Table{tt.table}
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
		fields []*dmodel.Field
		tables []*dmodel.Table
	}{
		{
			name:   "no join query",
			query:  "SELECT u.name FROM users AS u",
			fields: []*dmodel.Field{{Columns: []*dmodel.FieldColumn{{Table: "u", Name: "name", Type: dmodel.FieldReference}}}},
			tables: []*dmodel.Table{{AsName: "u", Name: "users"}},
		},
		{
			name: "single join query",
			query: `
				SELECT tasks.id
				FROM users AS u
				INNER JOIN tasks ON u.id = tasks.user_id
			`,
			fields: []*dmodel.Field{
				{Columns: []*dmodel.FieldColumn{{Table: "tasks", Name: "id", Type: dmodel.FieldReference}}},
				{Columns: []*dmodel.FieldColumn{{Table: "u", Name: "id", Type: dmodel.FieldCondition}, {Table: "tasks", Name: "user_id", Type: dmodel.FieldCondition}}},
			},
			tables: []*dmodel.Table{{AsName: "u", Name: "users"}, {Name: "tasks"}},
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
			fields: []*dmodel.Field{
				{Columns: []*dmodel.FieldColumn{{Table: "tasks", Name: "id", Type: dmodel.FieldReference}}},
				{Columns: []*dmodel.FieldColumn{{Table: "u", Name: "id", Type: dmodel.FieldCondition}, {Table: "tasks", Name: "user_id", Type: dmodel.FieldCondition}}},
				{Columns: []*dmodel.FieldColumn{{Table: "u", Name: "id", Type: dmodel.FieldCondition}, {Table: "t2", Name: "user_id", Type: dmodel.FieldCondition}}},
				{Columns: []*dmodel.FieldColumn{{Table: "u", Name: "id", Type: dmodel.FieldCondition}, {Table: "t3", Name: "user_id", Type: dmodel.FieldCondition}}},
			},
			tables: []*dmodel.Table{{AsName: "u", Name: "users"}, {Name: "tasks"}, {AsName: "t2", Name: "tasks"}, {AsName: "t3", Name: "tasks"}},
		},
		{
			name: "row join query",
			query: `
				SELECT tasks.id
				FROM users AS u
				INNER JOIN tasks ON ROW(u.id, u.id2) = ROW(tasks.user_id, tasks.user_id2)
			`,
			fields: []*dmodel.Field{
				{Columns: []*dmodel.FieldColumn{{Table: "tasks", Name: "id", Type: dmodel.FieldReference}}},
				{Columns: []*dmodel.FieldColumn{
					{Table: "u", Name: "id", Type: dmodel.FieldCondition},
					{Table: "u", Name: "id2", Type: dmodel.FieldCondition},
					{Table: "tasks", Name: "user_id", Type: dmodel.FieldCondition},
					{Table: "tasks", Name: "user_id2", Type: dmodel.FieldCondition},
				}},
			},
			tables: []*dmodel.Table{{AsName: "u", Name: "users"}, {Name: "tasks"}},
		},
		{
			name: "subquery join query",
			query: `
				SELECT tasks.id
				FROM users AS u
				INNER JOIN (SELECT * FROM tasks) AS t ON u.id = t.user_id
			`,
			fields: []*dmodel.Field{
				{Columns: []*dmodel.FieldColumn{{Table: "tasks", Name: "id", Type: dmodel.FieldReference}}},
				{Columns: []*dmodel.FieldColumn{
					{Table: "u", Name: "id", Type: dmodel.FieldCondition},
					{Table: "t", Name: "user_id", Type: dmodel.FieldCondition},
				}},
			},
			tables: []*dmodel.Table{{AsName: "u", Name: "users"}, {AsName: "t"}},
		},
	}
	p := parser.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmtNodes, _, err := p.Parse(tt.query, "", "")
			require.NoError(t, err)
			actualScopes, errs := mservice.CollectStmtScopes(stmtNodes[0])
			require.Empty(t, errs)

			expectedScopes := []*dmodel.StmtScope{{
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
		fields []*dmodel.Field
		tables []*dmodel.Table
	}{
		{
			name:   "subquery",
			query:  "SELECT id, (SELECT COUNT(status) FROM tasks) AS status_count FROM users",
			fields: []*dmodel.Field{{Columns: []*dmodel.FieldColumn{{Name: "", Type: dmodel.FieldSubquery}}}},
			tables: []*dmodel.Table{},
		},
		{
			name:   "subquery with comparison",
			query:  "SELECT (SELECT COUNT(status) FROM tasks) - (SELECT COUNT(description) FROM tasks) AS no_desc_count",
			fields: []*dmodel.Field{{Columns: []*dmodel.FieldColumn{{Name: "", Type: dmodel.FieldSubquery}}}},
			tables: []*dmodel.Table{},
		},
		{
			name:   "exists subquery",
			query:  "SELECT id, EXISTS (SELECT * FROM tasks WHERE user_id = users.user_id) FROM users",
			fields: []*dmodel.Field{{Columns: []*dmodel.FieldColumn{{Name: "", Type: dmodel.FieldSubquery}}}},
			tables: []*dmodel.Table{},
		},
		{
			name:   "in subquery",
			query:  "SELECT user_id IN (SELECT user_id FROM tasks) FROM users",
			fields: []*dmodel.Field{{Columns: []*dmodel.FieldColumn{{Name: "", Type: dmodel.FieldSubquery}}}},
			tables: []*dmodel.Table{},
		},
	}
	p := parser.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Skipf("Nest not implemented yet")
			stmtNodes, _, err := p.Parse(tt.query, "", "")
			require.NoError(t, err)
			actualScopes, errs := mservice.CollectStmtScopes(stmtNodes[0])
			require.Empty(t, errs)

			require.Equal(t, 1, len(actualScopes))
			assert.ElementsMatch(t, tt.fields, actualScopes[0].Fields)
			assert.ElementsMatch(t, tt.tables, actualScopes[0].Tables)
			expectedScopes := []*dmodel.StmtScope{{
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
		fields []*dmodel.Field
		table  *dmodel.Table
	}{
		{
			name:  "simple query(column at left)",
			query: "SELECT foo FROM users WHERE name = 'hello'",
			fields: []*dmodel.Field{
				{Columns: []*dmodel.FieldColumn{{Name: "foo", Type: dmodel.FieldReference}}},
				{Columns: []*dmodel.FieldColumn{{Name: "name", Type: dmodel.FieldCondition}}},
			},
			table: &dmodel.Table{Name: "users"},
		},
		{
			name:  "simple query(column at right",
			query: "SELECT foo FROM users WHERE 'hello' = name",
			fields: []*dmodel.Field{
				{Columns: []*dmodel.FieldColumn{{Name: "foo", Type: dmodel.FieldReference}}},
				{Columns: []*dmodel.FieldColumn{{Name: "name", Type: dmodel.FieldCondition}}},
			},
			table: &dmodel.Table{Name: "users"},
		},
		{
			name:  "same name with selected column",
			query: "SELECT name FROM users WHERE name = 'hello'",
			fields: []*dmodel.Field{
				{Columns: []*dmodel.FieldColumn{{Name: "name", Type: dmodel.FieldReference}}},
				{Columns: []*dmodel.FieldColumn{{Name: "name", Type: dmodel.FieldCondition}}},
			},
			table: &dmodel.Table{Name: "users"},
		},
		{
			name:  "compare with column",
			query: "SELECT foo FROM users WHERE name = first_name",
			fields: []*dmodel.Field{
				{Columns: []*dmodel.FieldColumn{{Name: "foo", Type: dmodel.FieldReference}}},
				{Columns: []*dmodel.FieldColumn{
					{Name: "name", Type: dmodel.FieldCondition},
					{Name: "first_name", Type: dmodel.FieldCondition},
				}},
			},
			table: &dmodel.Table{Name: "users"},
		},
		{
			name:  "multiple comparison",
			query: "SELECT foo FROM users WHERE first_name = 'hello' AND middle_name = 'nice' AND last_name = 'world'",
			fields: []*dmodel.Field{
				{Columns: []*dmodel.FieldColumn{{Name: "foo", Type: dmodel.FieldReference}}},
				{Columns: []*dmodel.FieldColumn{
					{Name: "first_name", Type: dmodel.FieldCondition},
					{Name: "middle_name", Type: dmodel.FieldCondition},
					{Name: "last_name", Type: dmodel.FieldCondition},
				}},
			},
			table: &dmodel.Table{Name: "users"},
		},
	}
	p := parser.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmtNodes, _, err := p.Parse(tt.query, "", "")
			require.NoError(t, err)
			actualScopes, errs := mservice.CollectStmtScopes(stmtNodes[0])
			require.Empty(t, errs)

			expectedScopes := []*dmodel.StmtScope{{
				Name:   "<root>",
				Fields: tt.fields,
				Tables: []*dmodel.Table{tt.table},
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
		fields []*dmodel.Field
		tables []*dmodel.Table
	}{
		{
			name:  "simple query",
			query: "SELECT COUNT(*) FROM users GROUP BY name",
			fields: []*dmodel.Field{
				{Columns: []*dmodel.FieldColumn{{Name: "name", Type: dmodel.FieldReference}}},
			},
			tables: []*dmodel.Table{{Name: "users"}},
		},
		{
			name:   "group by value",
			query:  "SELECT COUNT(*) FROM users GROUP BY 1",
			fields: nil,
			tables: []*dmodel.Table{{Name: "users"}},
		},
		{
			name:   "group by null",
			query:  "SELECT COUNT(*) FROM users GROUP BY null",
			fields: nil,
			tables: []*dmodel.Table{{Name: "users"}},
		},
		{
			name:  "multiple columns",
			query: "SELECT COUNT(*) FROM users GROUP BY first_name, middle_name, last_name",
			fields: []*dmodel.Field{
				{Columns: []*dmodel.FieldColumn{{Name: "first_name", Type: dmodel.FieldReference}}},
				{Columns: []*dmodel.FieldColumn{{Name: "middle_name", Type: dmodel.FieldReference}}},
				{Columns: []*dmodel.FieldColumn{{Name: "last_name", Type: dmodel.FieldReference}}},
			},
			tables: []*dmodel.Table{{Name: "users"}},
		},
		/*
			Not supporting WITH ROLLUP because pingcap/tidb doesn't support it.
			https://github.com/pingcap/tidb/issues/4250
			{
				name:  "with rollup",
				query: "SELECT COUNT(*) FROM users GROUP BY name WITH ROLLUP",
				fields: []*dmodel.Field{
					{Columns: []*dmodel.FieldColumn{{Name: "name", Type: dmodel.FieldReference}}},
				},
				tables: []*dmodel.Table{{Name: "users"}},
			},
		*/
	}
	p := parser.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmtNodes, _, err := p.Parse(tt.query, "", "")
			require.NoError(t, err)
			actualScopes, errs := mservice.CollectStmtScopes(stmtNodes[0])
			require.Empty(t, errs)

			expectedScopes := []*dmodel.StmtScope{{
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
		fields []*dmodel.Field
		tables []*dmodel.Table
	}{
		{
			name:  "simple query",
			query: "SELECT count(*) FROM users HAVING COUNT(name) > 1",
			fields: []*dmodel.Field{
				{Columns: []*dmodel.FieldColumn{{Name: "name", Type: dmodel.FieldReference}}},
			},
			tables: []*dmodel.Table{{Name: "users"}},
		},
		{
			name:  "comparison",
			query: "SELECT count(*) FROM users HAVING COUNT(first_name) > COUNT(last_name)",
			fields: []*dmodel.Field{
				{Columns: []*dmodel.FieldColumn{
					{Name: "first_name", Type: dmodel.FieldReference},
					{Name: "last_name", Type: dmodel.FieldReference},
				}},
			},
			tables: []*dmodel.Table{{Name: "users"}},
		},
	}
	p := parser.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmtNodes, _, err := p.Parse(tt.query, "", "")
			require.NoError(t, err)
			actualScopes, errs := mservice.CollectStmtScopes(stmtNodes[0])
			require.Empty(t, errs)

			expectedScopes := []*dmodel.StmtScope{{
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
		fields []*dmodel.Field
		tables []*dmodel.Table
	}{
		{
			name:  "simple query",
			query: "SELECT 1 FROM users ORDER BY name",
			fields: []*dmodel.Field{
				{Columns: []*dmodel.FieldColumn{{Name: "name", Type: dmodel.FieldReference}}},
			},
			tables: []*dmodel.Table{{Name: "users"}},
		},
		{
			name:   "order by value",
			query:  "SELECT COUNT(*) FROM users ORDER BY 1",
			fields: nil,
			tables: []*dmodel.Table{{Name: "users"}},
		},
		{
			name:   "group by null",
			query:  "SELECT COUNT(*) FROM users ORDER BY null",
			fields: nil,
			tables: []*dmodel.Table{{Name: "users"}},
		},
		{
			name:  "multiple columns",
			query: "SELECT COUNT(*) FROM users ORDER BY first_name, middle_name, last_name",
			fields: []*dmodel.Field{
				{Columns: []*dmodel.FieldColumn{{Name: "first_name", Type: dmodel.FieldReference}}},
				{Columns: []*dmodel.FieldColumn{{Name: "middle_name", Type: dmodel.FieldReference}}},
				{Columns: []*dmodel.FieldColumn{{Name: "last_name", Type: dmodel.FieldReference}}},
			},
			tables: []*dmodel.Table{{Name: "users"}},
		},
	}
	p := parser.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmtNodes, _, err := p.Parse(tt.query, "", "")
			require.NoError(t, err)
			actualScopes, errs := mservice.CollectStmtScopes(stmtNodes[0])
			require.Empty(t, errs)

			expectedScopes := []*dmodel.StmtScope{{
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
		fields []*dmodel.Field
		tables []*dmodel.Table
	}{
		{
			name:   "simple query",
			query:  "SELECT 1 FROM users LIMIT 1",
			fields: nil,
			tables: []*dmodel.Table{{Name: "users"}},
		},
		{
			name:   "offset exists",
			query:  "SELECT 1 FROM users LIMIT 1 OFFSET 1",
			fields: nil,
			tables: []*dmodel.Table{{Name: "users"}},
		},
		{
			name:   "offset without OFFSET keyword",
			query:  "SELECT 1 FROM users LIMIT 1, 1",
			fields: nil,
			tables: []*dmodel.Table{{Name: "users"}},
		},
	}
	p := parser.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmtNodes, _, err := p.Parse(tt.query, "", "")
			require.NoError(t, err)
			actualScopes, errs := mservice.CollectStmtScopes(stmtNodes[0])
			require.Empty(t, errs)

			expectedScopes := []*dmodel.StmtScope{{
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