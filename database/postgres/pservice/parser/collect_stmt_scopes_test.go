package parser_test

import (
	pParser "github.com/auxten/postgresql-parser/pkg/sql/parser"
	"github.com/google/go-cmp/cmp"
	"github.com/mrasu/GravityR/database/dservice/dparser"
	"github.com/mrasu/GravityR/database/postgres/pservice"
	"github.com/mrasu/GravityR/database/postgres/pservice/parser"
	"github.com/mrasu/GravityR/thelper/tdata"
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
			query:  "SELECT user_id IN (1, 2, 3, parent_id) FROM users",
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
			name:   "similar to query",
			query:  "SELECT name SIMILAR TO valid_pattern FROM users",
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
			name:   "trim query",
			query:  "SELECT TRIM(name) FROM users",
			fields: []*dparser.Field{{Columns: []*dparser.FieldColumn{{Name: "name", Type: dparser.FieldReference}}}},
			table:  &dparser.Table{Name: "users"},
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
			name:   "prepared statement query",
			query:  "SELECT $1",
			fields: nil,
			table:  nil,
		},
		{
			name:   "array query",
			query:  "SELECT ARRAY[1, id] FROM users",
			fields: []*dparser.Field{{Columns: []*dparser.FieldColumn{{Name: "id", Type: dparser.FieldReference}}}},
			table:  &dparser.Table{Name: "users"},
		},
		{
			name:   "nullif query",
			query:  "SELECT NULLIF(id, 1) FROM users",
			fields: []*dparser.Field{{Columns: []*dparser.FieldColumn{{Name: "id", Type: dparser.FieldReference}}}},
			table:  &dparser.Table{Name: "users"},
		},
		{
			name:   "coalesce query",
			query:  "SELECT COALESCE(id, 1) FROM users",
			fields: []*dparser.Field{{Columns: []*dparser.FieldColumn{{Name: "id", Type: dparser.FieldReference}}}},
			table:  &dparser.Table{Name: "users"},
		},
		{
			name:   "collate query",
			query:  "SELECT name COLLATE \"C\" FROM users",
			fields: []*dparser.Field{{Columns: []*dparser.FieldColumn{{Name: "name", Type: dparser.FieldReference}}}},
			table:  &dparser.Table{Name: "users"},
		},
		{
			name:   "tuple star query",
			query:  "SELECT (users).* FROM users",
			fields: []*dparser.Field{{Columns: []*dparser.FieldColumn{{Table: "users", Type: dparser.FieldStar}}}},
			table:  &dparser.Table{Name: "users"},
		},
		{
			name:   "tuple star query",
			query:  "SELECT (users).id FROM users",
			fields: []*dparser.Field{{Columns: []*dparser.FieldColumn{{Table: "users", Name: "id", Type: dparser.FieldReference}}}},
			table:  &dparser.Table{Name: "users"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := toStatement(t, tt.query)
			actualScopes, errs := parser.CollectStmtScopes(stmt, "public")
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := toStatement(t, tt.query)
			actualScopes, errs := parser.CollectStmtScopes(stmt, "public")
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

func TestCollectStmtScopes_PostgresSubquery(t *testing.T) {
	tests := []struct {
		name           string
		sql            string
		expectedScopes []*dparser.StmtScope
	}{
		{
			name:           "subquery in SELECT",
			sql:            tdata.PostgresSubqueryInSelectData.SQL,
			expectedScopes: tdata.PostgresSubqueryInSelectData.Scopes,
		},
		{
			name:           "subquery in SELECT's aliased field",
			sql:            tdata.PostgresSubqueryInSelectAliasedFieldData.SQL,
			expectedScopes: tdata.PostgresSubqueryInSelectAliasedFieldData.Scopes,
		},
		{
			name:           "subquery with comparison",
			sql:            tdata.PostgresSubqueryWithComparisonData.SQL,
			expectedScopes: tdata.PostgresSubqueryWithComparisonData.Scopes,
		},
		{
			name:           "subquery with EXIST",
			sql:            tdata.PostgresSubqueryWithExistData.SQL,
			expectedScopes: tdata.PostgresSubqueryWithExistData.Scopes,
		},
		{
			name:           "subquery with IN",
			sql:            tdata.PostgresSubqueryWithInData.SQL,
			expectedScopes: tdata.PostgresSubqueryWithInData.Scopes,
		},
		{
			name:           "subquery in SELECT's function",
			sql:            tdata.PostgresSubqueryInSelectFunctionData.SQL,
			expectedScopes: tdata.PostgresSubqueryInSelectFunctionData.Scopes,
		},
		{
			name:           "subquery in FROM",
			sql:            tdata.PostgresSubqueryInFromData.SQL,
			expectedScopes: tdata.PostgresSubqueryInFromData.Scopes,
		},
		{
			name:           "subquery in JOIN",
			sql:            tdata.PostgresSubqueryInJoinData.SQL,
			expectedScopes: tdata.PostgresSubqueryInJoinData.Scopes,
		},
		{
			name:           "subquery using star",
			sql:            tdata.PostgresSubqueryUsingStarData.SQL,
			expectedScopes: tdata.PostgresSubqueryUsingStarData.Scopes,
		},
		{
			name:           "subquery with LATERAL JOIN",
			sql:            tdata.PostgresSubqueryWithLateralJoinData.SQL,
			expectedScopes: tdata.PostgresSubqueryWithLateralJoinData.Scopes,
		},
		{
			name:           "subquery with LATERAL object reference",
			sql:            tdata.PostgresSubqueryWithLateralObjectReferenceData.SQL,
			expectedScopes: tdata.PostgresSubqueryWithLateralObjectReferenceData.Scopes,
		},
		{
			name:           "subquery with aggregate function",
			sql:            tdata.PostgresSubqueryWithAggregateFunctionData.SQL,
			expectedScopes: tdata.PostgresSubqueryWithAggregateFunctionData.Scopes,
		},
		{
			name:           "correlated subquery",
			sql:            tdata.PostgresCorrelatedSubqueryData.SQL,
			expectedScopes: tdata.PostgresCorrelatedSubqueryData.Scopes,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := toStatement(t, tt.sql)
			actualScopes, errs := parser.CollectStmtScopes(stmt, "public")
			require.Empty(t, errs)

			if diff := cmp.Diff(tt.expectedScopes, actualScopes); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestCollectStmtScopes_HasuraSubquery(t *testing.T) {
	tests := []struct {
		name           string
		sql            string
		expectedScopes []*dparser.StmtScope
	}{
		{
			name:           "subquery with where",
			sql:            tdata.HasuraSubqueryWithWhereData.SQL,
			expectedScopes: tdata.HasuraSubqueryWithWhereData.Scopes,
		},
		{
			name:           "subquery with variables",
			sql:            tdata.HasuraSubqueryWithVariablesData.SQL,
			expectedScopes: tdata.HasuraSubqueryWithVariablesData.Scopes,
		},
		{
			name:           "subquery with aggregation",
			sql:            tdata.HasuraSubqueryWithAggregationData.SQL,
			expectedScopes: tdata.HasuraSubqueryWithAggregationData.Scopes,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := toStatement(t, tt.sql)
			actualScopes, errs := parser.CollectStmtScopes(stmt, "public")
			require.Empty(t, errs)

			if diff := cmp.Diff(tt.expectedScopes, actualScopes); diff != "" {
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := toStatement(t, tt.query)
			actualScopes, errs := parser.CollectStmtScopes(stmt, "public")
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
			query:  "SELECT COUNT(*) FROM users GROUP BY NULL",
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
		{
			name:  "rollup",
			query: "SELECT COUNT(*) FROM users GROUP BY ROLLUP(name)",
			fields: []*dparser.Field{
				{Columns: []*dparser.FieldColumn{{Name: "name", Type: dparser.FieldReference}}},
			},
			tables: []*dparser.Table{{Name: "users"}},
		},
		{
			name:  "cube",
			query: "SELECT COUNT(*) FROM users GROUP BY CUBE(name)",
			fields: []*dparser.Field{
				{Columns: []*dparser.FieldColumn{{Name: "name", Type: dparser.FieldReference}}},
			},
			tables: []*dparser.Table{{Name: "users"}},
		},
		/*
			Not supporting GROUPING SETS because cockroachdb doesn't support it
			https://github.com/cockroachdb/cockroach/issues/46280
			{
				name:  "grouping sets",
				query: "SELECT COUNT(*) FROM users GROUP BY GROUPING SETS(name)",
				fields: []*dparser.Field{
					{Columns: []*dparser.FieldColumn{{Name: "name", Type: dparser.FieldReference}}},
				},
				tables: []*dparser.Table{{Name: "users"}},
			},
		*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := toStatement(t, tt.query)
			actualScopes, errs := parser.CollectStmtScopes(stmt, "public")
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
			query: "SELECT COUNT(*) FROM users HAVING COUNT(name) > 1",
			fields: []*dparser.Field{
				{Columns: []*dparser.FieldColumn{{Name: "name", Type: dparser.FieldReference}}},
			},
			tables: []*dparser.Table{{Name: "users"}},
		},
		{
			name:  "comparison",
			query: "SELECT COUNT(*) FROM users HAVING COUNT(first_name) > COUNT(last_name)",
			fields: []*dparser.Field{
				{Columns: []*dparser.FieldColumn{
					{Name: "first_name", Type: dparser.FieldReference},
					{Name: "last_name", Type: dparser.FieldReference},
				}},
			},
			tables: []*dparser.Table{{Name: "users"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := toStatement(t, tt.query)
			actualScopes, errs := parser.CollectStmtScopes(stmt, "public")
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
			query:  "SELECT COUNT(*) FROM users ORDER BY NULL",
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := toStatement(t, tt.query)
			actualScopes, errs := parser.CollectStmtScopes(stmt, "public")
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := toStatement(t, tt.query)
			actualScopes, errs := parser.CollectStmtScopes(stmt, "public")
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

func TestCollectStmtScopes_SchemaSpecified(t *testing.T) {
	tests := []struct {
		name   string
		query  string
		fields []*dparser.Field
		tables []*dparser.Table
		errMsg string
	}{
		{
			name:  "specifying table name in expression",
			query: "SELECT users.id FROM users",
			fields: []*dparser.Field{
				{Columns: []*dparser.FieldColumn{{Table: "users", Name: "id", Type: dparser.FieldReference}}},
			},
			tables: []*dparser.Table{{Name: "users"}},
		},
		{
			name:  "specifying schema name, public, in expression",
			query: "SELECT public.users.id FROM users",
			fields: []*dparser.Field{
				{Columns: []*dparser.FieldColumn{{Table: "users", Name: "id", Type: dparser.FieldReference}}},
			},
			tables: []*dparser.Table{{Name: "users"}},
		},
		{
			name:   "specifying invalid schema name in expression",
			query:  "SELECT dummy.users.id FROM users",
			errMsg: "schema or catalog/db specification is not supported",
		},
		{
			name:   "specifying catalog name in table name",
			query:  "SELECT foo.public.users.id FROM users",
			errMsg: "schema or catalog/db specification is not supported",
		},
		{
			name:  "specifying schema name, public, in table name",
			query: "SELECT id FROM public.users",
			fields: []*dparser.Field{
				{Columns: []*dparser.FieldColumn{{Name: "id", Type: dparser.FieldReference}}},
			},
			tables: []*dparser.Table{{Name: "users"}},
		},
		{
			name:   "specifying invalid schema name in table name",
			query:  "SELECT id FROM dummy.users",
			errMsg: "schema or catalog/db specification is not supported",
		},
		{
			name:   "specifying catalog name in table name",
			query:  "SELECT id FROM foo.public.users",
			errMsg: "schema or catalog/db specification is not supported",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := toStatement(t, tt.query)
			actualScopes, errs := parser.CollectStmtScopes(stmt, "public")
			if tt.errMsg != "" {
				require.Equal(t, 1, len(errs))
				assert.Contains(t, errs[0].Error(), tt.errMsg)
			} else {
				require.Empty(t, errs)
				expectedScopes := []*dparser.StmtScope{{
					Name:   "<root>",
					Fields: tt.fields,
					Tables: tt.tables,
				}}
				if diff := cmp.Diff(expectedScopes, actualScopes); diff != "" {
					t.Errorf(diff)
				}
			}
		})
	}
}

func toStatement(t *testing.T, query string) *pParser.Statement {
	t.Helper()
	stmt, err := pservice.Parse(query)
	require.NoError(t, err)

	return stmt
}
