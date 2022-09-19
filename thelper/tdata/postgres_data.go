package tdata

import "github.com/mrasu/GravityR/database/common_model"

type postgresData struct {
	SQL    string
	Scopes []*common_model.StmtScope
}

var PostgresSubqueryInSelectData = &postgresData{
	SQL: "SELECT (SELECT user_id FROM todos limit 1)",
	Scopes: []*common_model.StmtScope{
		{
			Name: "<root>",
			Fields: []*common_model.Field{
				{Columns: []*common_model.FieldColumn{{ReferenceName: "<field0>", Type: common_model.FieldSubquery}}},
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
		},
	},
}

var PostgresSubqueryInSelectAliasedFieldData = &postgresData{
	SQL: "SELECT name, is_admin, (SELECT COUNT(status) FROM todos) AS status_count FROM users",
	Scopes: []*common_model.StmtScope{{
		Name: "<root>",
		Fields: []*common_model.Field{
			{Columns: []*common_model.FieldColumn{{Name: "name", Type: common_model.FieldReference}}},
			{Columns: []*common_model.FieldColumn{{Name: "is_admin", Type: common_model.FieldReference}}},
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
}

var PostgresSubqueryWithComparisonData = &postgresData{
	SQL: "SELECT (SELECT COUNT(status) FROM todos) - (SELECT COUNT(description) FROM todos) AS no_desc_count",
	Scopes: []*common_model.StmtScope{{
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
}

var PostgresSubqueryWithExistData = &postgresData{
	SQL: "SELECT name, EXISTS (SELECT * FROM todos WHERE user_id = users.id) FROM users",
	Scopes: []*common_model.StmtScope{{
		Name: "<root>",
		Fields: []*common_model.Field{
			{Columns: []*common_model.FieldColumn{{Name: "name", Type: common_model.FieldReference}}},
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
}

var PostgresSubqueryWithInData = &postgresData{
	SQL: "SELECT user_id IN (SELECT user_id FROM todos) FROM users",
	Scopes: []*common_model.StmtScope{{
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
}

var PostgresSubqueryInSelectFunctionData = &postgresData{
	SQL: `
SELECT row_to_json(
  (
    SELECT t
    FROM (SELECT id, description FROM todos LIMIT 1) AS t
  )
)
`,
	Scopes: []*common_model.StmtScope{
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
					SubScopes: []*common_model.StmtScope{
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
}

var PostgresSubqueryInFromData = &postgresData{
	SQL: `
SELECT *
FROM (SELECT id, user_id FROM todos)
`,
	Scopes: []*common_model.StmtScope{
		{
			Name: "<root>",
			Fields: []*common_model.Field{
				{Columns: []*common_model.FieldColumn{{Type: common_model.FieldStar}}},
			},
			Tables: []*common_model.Table{
				{Name: "<select0>"},
			},
			SubScopes: []*common_model.StmtScope{
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
}

var PostgresSubqueryInJoinData = &postgresData{
	SQL: `
SELECT name, t.id, user_id
FROM users
INNER JOIN (SELECT id, user_id FROM todos) AS t ON users.id = t.user_id
`,
	Scopes: []*common_model.StmtScope{
		{
			Name: "<root>",
			Fields: []*common_model.Field{
				{Columns: []*common_model.FieldColumn{{Name: "name", Type: common_model.FieldReference}}},
				{Columns: []*common_model.FieldColumn{{Table: "t", Name: "id", Type: common_model.FieldReference}}},
				{Columns: []*common_model.FieldColumn{{Name: "user_id", Type: common_model.FieldReference}}},
				{Columns: []*common_model.FieldColumn{
					{Table: "users", Name: "id", Type: common_model.FieldCondition},
					{Table: "t", Name: "user_id", Type: common_model.FieldCondition},
				}},
			},
			Tables: []*common_model.Table{
				{Name: "users"},
				{Name: "<select0>", AsName: "t"},
			},
			SubScopes: []*common_model.StmtScope{
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
}

var PostgresSubqueryUsingStarData = &postgresData{
	SQL: `
SELECT description
FROM
	(SELECT * FROM todos) AS t
`,
	Scopes: []*common_model.StmtScope{
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
			SubScopes: []*common_model.StmtScope{
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
}

var PostgresSubqueryWithLateralJoinData = &postgresData{
	SQL: `
SELECT description
FROM
	(SELECT id FROM users) AS u
	LEFT OUTER JOIN LATERAL (SELECT user_id, description FROM todos WHERE user_id = u.id) AS t ON u.id = t.user_id
`,
	Scopes: []*common_model.StmtScope{
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
			SubScopes: []*common_model.StmtScope{
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
						{Columns: []*common_model.FieldColumn{{Name: "user_id", Type: common_model.FieldReference}}},
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
}

var PostgresSubqueryWithLateralObjectReferenceData = &postgresData{
	SQL: `
SELECT b
FROM (SELECT name FROM users) AS a
	LEFT JOIN LATERAL (SELECT a) AS b ON (TRUE)
`,
	Scopes: []*common_model.StmtScope{
		{
			Name: "<root>",
			Fields: []*common_model.Field{
				{Columns: []*common_model.FieldColumn{
					{Name: "b", Type: common_model.FieldReference},
				}},
			},
			Tables: []*common_model.Table{
				{Name: "<select0>", AsName: "a"},
				{Name: "<select1>", AsName: "b", IsLateral: true},
			},
			SubScopes: []*common_model.StmtScope{
				{
					Name: "<select0>",
					Fields: []*common_model.Field{
						{Columns: []*common_model.FieldColumn{{Name: "name", Type: common_model.FieldReference}}},
					},
					Tables: []*common_model.Table{
						{Name: "users"},
					},
				},
				{
					Name: "<select1>",
					Fields: []*common_model.Field{
						{Columns: []*common_model.FieldColumn{{Name: "a", Type: common_model.FieldReference}}},
					},
				},
			},
		},
	},
}

var PostgresSubqueryWithAggregateFunctionData = &postgresData{
	SQL: `
SELECT json_agg(t)
FROM (SELECT description FROM todos LIMIT 3) AS t;
`,
	Scopes: []*common_model.StmtScope{
		{
			Name: "<root>",
			Fields: []*common_model.Field{
				{Columns: []*common_model.FieldColumn{{Name: "t", Type: common_model.FieldReference}}},
			},
			Tables: []*common_model.Table{
				{AsName: "t", Name: "<select0>"},
			},
			SubScopes: []*common_model.StmtScope{
				{
					Name: "<select0>",
					Fields: []*common_model.Field{
						{Columns: []*common_model.FieldColumn{{Name: "description", Type: common_model.FieldReference}}},
					},
					Tables: []*common_model.Table{
						{Name: "todos"},
					},
				},
			},
		},
	},
}

var PostgresCorrelatedSubqueryData = &postgresData{
	SQL: `
SELECT description
FROM todos
WHERE EXISTS(SELECT * FROM users WHERE users.id = user_id)
`,
	Scopes: []*common_model.StmtScope{
		{
			Name: "<root>",
			Fields: []*common_model.Field{
				{Columns: []*common_model.FieldColumn{{Name: "description", Type: common_model.FieldReference}}},
				{Columns: []*common_model.FieldColumn{{ReferenceName: "<field0>", Type: common_model.FieldSubquery}}},
			},
			FieldScopes: []*common_model.StmtScope{
				{
					Name: "<field0>",
					Fields: []*common_model.Field{
						{Columns: []*common_model.FieldColumn{{Type: common_model.FieldStar}}},
						{Columns: []*common_model.FieldColumn{
							{Table: "users", Name: "id", Type: common_model.FieldCondition},
							{Name: "user_id", Type: common_model.FieldCondition},
						}},
					},
					Tables: []*common_model.Table{
						{Name: "users"},
					},
				},
			},
			Tables: []*common_model.Table{
				{Name: "todos"},
			},
		},
	},
}
