package collectors_test

import (
	"github.com/mrasu/GravityR/database/mysql/models"
	"github.com/mrasu/GravityR/database/mysql/models/collectors"
	"github.com/pingcap/tidb/parser"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCollectStmtScopes_SingleField(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		asName  string
		columns []*models.FieldColumn
		table   *models.Table
	}{
		{
			name:    "value query",
			query:   "SELECT 1",
			asName:  "",
			columns: []*models.FieldColumn{},
			table:   nil,
		},
		{
			name:    "column reference query",
			query:   "SELECT name FROM users AS u",
			asName:  "",
			columns: []*models.FieldColumn{{Name: "name", Type: models.FieldReference}},
			table:   &models.Table{AsName: "u", Name: "users"},
		},
		{
			name:    "column reference query with table name",
			query:   "SELECT u.name FROM users AS u",
			columns: []*models.FieldColumn{{Table: "u", Name: "name", Type: models.FieldReference}},
			table:   &models.Table{AsName: "u", Name: "users"},
		},
		{
			name:    "star query",
			query:   "SELECT * FROM users",
			asName:  "",
			columns: []*models.FieldColumn{{Type: models.FieldStar}},
			table:   &models.Table{Name: "users"},
		},
		{
			name:    "star query with table name",
			query:   "SELECT users.* FROM users",
			asName:  "",
			columns: []*models.FieldColumn{{Table: "users", Type: models.FieldStar}},
			table:   &models.Table{Name: "users"},
		},
		{
			name:    "aliased query",
			query:   "SELECT name AS name1 FROM users",
			asName:  "name1",
			columns: []*models.FieldColumn{{Name: "name", Type: models.FieldReference}},
			table:   &models.Table{Name: "users"},
		},
		{
			name:    "between query",
			query:   "SELECT lunch_at BETWEEN start_at AND end_at FROM users",
			columns: []*models.FieldColumn{{Name: "lunch_at", Type: models.FieldReference}, {Name: "start_at", Type: models.FieldReference}, {Name: "end_at", Type: models.FieldReference}},
			table:   &models.Table{Name: "users"},
		},
		{
			name:    "operator query",
			query:   "SELECT end_at - start_at FROM users",
			columns: []*models.FieldColumn{{Name: "start_at", Type: models.FieldReference}, {Name: "end_at", Type: models.FieldReference}},
			table:   &models.Table{Name: "users"},
		},
		{
			name:    "case query",
			query:   "SELECT CASE WHEN name = last_name THEN full_name ELSE first_name END FROM users",
			columns: []*models.FieldColumn{{Name: "name", Type: models.FieldReference}, {Name: "last_name", Type: models.FieldReference}, {Name: "full_name", Type: models.FieldReference}, {Name: "first_name", Type: models.FieldReference}},
			table:   &models.Table{Name: "users"},
		},
		{
			name:    "in query",
			query:   "SELECT user_id in (1, 2, 3, parent_id) FROM users",
			columns: []*models.FieldColumn{{Name: "user_id", Type: models.FieldReference}, {Name: "parent_id", Type: models.FieldReference}},
			table:   &models.Table{Name: "users"},
		},
		{
			name:    "is null query",
			query:   "SELECT user_id IS NULL FROM users",
			columns: []*models.FieldColumn{{Name: "user_id", Type: models.FieldReference}},
			table:   &models.Table{Name: "users"},
		},
		{
			name:    "is bool query",
			query:   "SELECT is_valid IS TRUE FROM users",
			columns: []*models.FieldColumn{{Name: "is_valid", Type: models.FieldReference}},
			table:   &models.Table{Name: "users"},
		},
		{
			name:    "like query",
			query:   "SELECT name LIKE valid_prefix FROM users",
			columns: []*models.FieldColumn{{Name: "name", Type: models.FieldReference}, {Name: "valid_prefix", Type: models.FieldReference}},
			table:   &models.Table{Name: "users"},
		},
		{
			name:    "parenthesis query",
			query:   "SELECT count1 + (count2 * count3) FROM users",
			columns: []*models.FieldColumn{{Name: "count1", Type: models.FieldReference}, {Name: "count2", Type: models.FieldReference}, {Name: "count3", Type: models.FieldReference}},
			table:   &models.Table{Name: "users"},
		},
		{
			name:    "regex query",
			query:   "SELECT name REGEXP valid_pattern FROM users",
			columns: []*models.FieldColumn{{Name: "name", Type: models.FieldReference}, {Name: "valid_pattern", Type: models.FieldReference}},
			table:   &models.Table{Name: "users"},
		},
		{
			name:    "unary operator query",
			query:   "SELECT -user_id FROM users",
			columns: []*models.FieldColumn{{Name: "user_id", Type: models.FieldReference}},
			table:   &models.Table{Name: "users"},
		},
		{
			name:    "variable query",
			query:   "SELECT @hello",
			columns: []*models.FieldColumn{},
			table:   nil,
		},
		{
			name:    "match query",
			query:   "SELECT MATCH(first_name, last_name) AGAINST(valid_pattern) FROM users",
			columns: []*models.FieldColumn{{Name: "first_name", Type: models.FieldReference}, {Name: "last_name", Type: models.FieldReference}, {Name: "valid_pattern", Type: models.FieldReference}},
			table:   &models.Table{Name: "users"},
		},
		{
			name:    "trim query",
			query:   "SELECT TRIM(name) FROM users",
			columns: []*models.FieldColumn{{Name: "name", Type: models.FieldReference}},
			table:   &models.Table{Name: "users"},
		},
		{
			name:    "trim query with direction",
			query:   "SELECT TRIM(LEADING ' ' FROM ' hello world ')",
			columns: []*models.FieldColumn{},
			table:   nil,
		},
		{
			name:    "trim query with direction and column reference",
			query:   "SELECT TRIM(LEADING ' ' FROM name) FROM users",
			columns: []*models.FieldColumn{{Name: "name", Type: models.FieldReference}},
			table:   &models.Table{Name: "users"},
		},
		{
			name:    "cast query",
			query:   "SELECT CAST(user_id AS CHAR) FROM users",
			columns: []*models.FieldColumn{{Name: "user_id", Type: models.FieldReference}},
			table:   &models.Table{Name: "users"},
		},
		{
			name:    "aggregation query",
			query:   "SELECT SUM(id) FROM users",
			columns: []*models.FieldColumn{{Name: "id", Type: models.FieldReference}},
			table:   &models.Table{Name: "users"},
		},
		{
			name:    "aggregation query with star",
			query:   "SELECT COUNT(*) FROM users",
			columns: []*models.FieldColumn{},
			table:   &models.Table{Name: "users"},
		},
		{
			name:    "window function query",
			query:   "SELECT AVG(todo_id) OVER(PARTITION BY user_id ORDER BY created_at ROWS BETWEEN 1 PRECEDING AND 2 FOLLOWING) FROM todos",
			columns: []*models.FieldColumn{{Name: "todo_id", Type: models.FieldReference}, {Name: "user_id", Type: models.FieldReference}, {Name: "created_at", Type: models.FieldReference}},
			table:   &models.Table{Name: "todos"},
		},
		{
			name:    "GET_FORMAT function query",
			query:   "SELECT GET_FORMAT(DATE, 'USA')",
			columns: []*models.FieldColumn{},
			table:   nil,
		},
		{
			name:    "GET_FORMAT function query",
			query:   "SELECT GET_FORMAT(DATE, name) FROM regions",
			columns: []*models.FieldColumn{{Name: "name", Type: models.FieldReference}},
			table:   &models.Table{Name: "regions"},
		},
		{
			name:    "prepared statement query",
			query:   "SELECT ?",
			columns: []*models.FieldColumn{},
			table:   nil,
		},
	}
	p := parser.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmtNodes, _, err := p.Parse(tt.query, "", "")
			assert.NoError(t, err)
			actualScopes, errs := collectors.CollectStmtScopes(stmtNodes[0])
			assert.Empty(t, errs)

			assert.Equal(t, 1, len(actualScopes))
			actualScope := actualScopes[0]

			assert.Equal(t, "<root>", actualScope.Name)
			assert.Equal(t, 1, len(actualScope.Fields))
			assert.Equal(t, tt.asName, actualScope.Fields[0].AsName)
			assert.ElementsMatch(t, tt.columns, actualScope.Fields[0].Columns)

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
		fields []*models.Field
		tables []*models.Table
	}{
		{
			name:   "no join query",
			query:  "SELECT u.name FROM users AS u",
			fields: []*models.Field{{Columns: []*models.FieldColumn{{Table: "u", Name: "name", Type: models.FieldReference}}}},
			tables: []*models.Table{{AsName: "u", Name: "users"}},
		},
		{
			name: "single join query",
			query: `
				SELECT todos.id
				FROM users AS u
				INNER JOIN todos ON u.id = todos.user_id
			`,
			fields: []*models.Field{
				{Columns: []*models.FieldColumn{{Table: "todos", Name: "id", Type: models.FieldReference}}},
				{Columns: []*models.FieldColumn{{Table: "u", Name: "id", Type: models.FieldCondition}, {Table: "todos", Name: "user_id", Type: models.FieldCondition}}},
			},
			tables: []*models.Table{{AsName: "u", Name: "users"}, {Name: "todos"}},
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
			fields: []*models.Field{
				{Columns: []*models.FieldColumn{{Table: "todos", Name: "id", Type: models.FieldReference}}},
				{Columns: []*models.FieldColumn{{Table: "u", Name: "id", Type: models.FieldCondition}, {Table: "todos", Name: "user_id", Type: models.FieldCondition}}},
				{Columns: []*models.FieldColumn{{Table: "u", Name: "id", Type: models.FieldCondition}, {Table: "t2", Name: "user_id", Type: models.FieldCondition}}},
				{Columns: []*models.FieldColumn{{Table: "u", Name: "id", Type: models.FieldCondition}, {Table: "t3", Name: "user_id", Type: models.FieldCondition}}},
			},
			tables: []*models.Table{{AsName: "u", Name: "users"}, {Name: "todos"}, {AsName: "t2", Name: "todos"}, {AsName: "t3", Name: "todos"}},
		},
		{
			name: "row join query",
			query: `
				SELECT todos.id
				FROM users AS u
				INNER JOIN todos ON ROW(u.id, u.id2) = ROW(todos.user_id, todos.user_id2)
			`,
			fields: []*models.Field{
				{Columns: []*models.FieldColumn{{Table: "todos", Name: "id", Type: models.FieldReference}}},
				{Columns: []*models.FieldColumn{
					{Table: "u", Name: "id", Type: models.FieldCondition},
					{Table: "u", Name: "id2", Type: models.FieldCondition},
					{Table: "todos", Name: "user_id", Type: models.FieldCondition},
					{Table: "todos", Name: "user_id2", Type: models.FieldCondition},
				}},
			},
			tables: []*models.Table{{AsName: "u", Name: "users"}, {Name: "todos"}},
		},
		{
			name: "subquery join query",
			query: `
				SELECT todos.id
				FROM users AS u
				INNER JOIN (SELECT * FROM todos) AS t ON u.id = t.user_id
			`,
			fields: []*models.Field{
				{Columns: []*models.FieldColumn{{Table: "todos", Name: "id", Type: models.FieldReference}}},
				{Columns: []*models.FieldColumn{
					{Table: "u", Name: "id", Type: models.FieldCondition},
					{Table: "t", Name: "user_id", Type: models.FieldCondition},
				}},
			},
			tables: []*models.Table{{AsName: "u", Name: "users"}, {AsName: "t"}},
		},
	}
	p := parser.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmtNodes, _, err := p.Parse(tt.query, "", "")
			assert.NoError(t, err)
			actualScopes, errs := collectors.CollectStmtScopes(stmtNodes[0])
			assert.Empty(t, errs)

			assert.Equal(t, 1, len(actualScopes))
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
		fields []*models.Field
		tables []string
	}{
		{
			name:   "subquery",
			query:  "SELECT id, (SELECT COUNT(status) FROM todos) AS status_count FROM users",
			fields: []*models.Field{{Columns: []*models.FieldColumn{{Name: "", Type: models.FieldSubquery}}}},
			tables: []string{},
		},
		{
			name:   "subquery with comparison",
			query:  "SELECT (SELECT COUNT(status) FROM todos) - (SELECT COUNT(description) FROM todos) AS no_desc_count",
			fields: []*models.Field{{Columns: []*models.FieldColumn{{Name: "", Type: models.FieldSubquery}}}},
			tables: []string{},
		},
		{
			name:   "exists subquery",
			query:  "SELECT id, EXISTS (SELECT * FROM todos WHERE user_id = users.user_id) FROM users",
			fields: []*models.Field{{Columns: []*models.FieldColumn{{Name: "", Type: models.FieldSubquery}}}},
			tables: []string{},
		},
		{
			name:   "in subquery",
			query:  "SELECT user_id IN (SELECT user_id FROM todos) FROM users",
			fields: []*models.Field{{Columns: []*models.FieldColumn{{Name: "", Type: models.FieldSubquery}}}},
			tables: []string{},
		},
	}
	p := parser.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Skipf("Nest not implemented yet")
			stmtNodes, _, err := p.Parse(tt.query, "", "")
			assert.NoError(t, err)
			actualScopes, errs := collectors.CollectStmtScopes(stmtNodes[0])
			assert.Empty(t, errs)

			assert.Equal(t, 1, len(actualScopes))
			assert.ElementsMatch(t, tt.fields, actualScopes[0].Fields)

			// TODO
			//assert.ElementsMatch(t, tt.tables, actualScopes[0].Tables)
		})
	}
}