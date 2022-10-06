package tdata

import "github.com/mrasu/GravityR/database/dmodel"

type postgresData struct {
	SQL    string
	Scopes []*dmodel.StmtScope
}

var PostgresSubqueryInSelectData = &postgresData{
	SQL: "SELECT (SELECT user_id FROM tasks limit 1)",
	Scopes: []*dmodel.StmtScope{
		{
			Name: "<root>",
			Fields: []*dmodel.Field{
				{Columns: []*dmodel.FieldColumn{{ReferenceName: "<field0>", Type: dmodel.FieldSubquery}}},
			},
			FieldScopes: []*dmodel.StmtScope{
				{
					Name: "<field0>",
					Fields: []*dmodel.Field{
						{Columns: []*dmodel.FieldColumn{{Name: "user_id", Type: dmodel.FieldReference}}},
					},
					Tables: []*dmodel.Table{
						{Name: "tasks"},
					},
				},
			},
		},
	},
}

var PostgresSubqueryInSelectAliasedFieldData = &postgresData{
	SQL: "SELECT name, is_admin, (SELECT COUNT(status) FROM tasks) AS status_count FROM users",
	Scopes: []*dmodel.StmtScope{{
		Name: "<root>",
		Fields: []*dmodel.Field{
			{Columns: []*dmodel.FieldColumn{{Name: "name", Type: dmodel.FieldReference}}},
			{Columns: []*dmodel.FieldColumn{{Name: "is_admin", Type: dmodel.FieldReference}}},
			{AsName: "status_count", Columns: []*dmodel.FieldColumn{{ReferenceName: "<field0>", Type: dmodel.FieldSubquery}}},
		},
		FieldScopes: []*dmodel.StmtScope{
			{
				Name: "<field0>",
				Fields: []*dmodel.Field{
					{Columns: []*dmodel.FieldColumn{{Name: "status", Type: dmodel.FieldReference}}},
				},
				Tables: []*dmodel.Table{
					{Name: "tasks"},
				},
			},
		},
		Tables: []*dmodel.Table{
			{Name: "users"},
		},
	}},
}

var PostgresSubqueryWithComparisonData = &postgresData{
	SQL: "SELECT (SELECT COUNT(status) FROM tasks) - (SELECT COUNT(description) FROM tasks) AS no_desc_count",
	Scopes: []*dmodel.StmtScope{{
		Name: "<root>",
		Fields: []*dmodel.Field{
			{AsName: "no_desc_count", Columns: []*dmodel.FieldColumn{
				{ReferenceName: "<field0>", Type: dmodel.FieldSubquery},
				{ReferenceName: "<field1>", Type: dmodel.FieldSubquery},
			}},
		},
		FieldScopes: []*dmodel.StmtScope{
			{
				Name: "<field0>",
				Fields: []*dmodel.Field{
					{Columns: []*dmodel.FieldColumn{{Name: "status", Type: dmodel.FieldReference}}},
				},
				Tables: []*dmodel.Table{
					{Name: "tasks"},
				},
			},
			{
				Name: "<field1>",
				Fields: []*dmodel.Field{
					{Columns: []*dmodel.FieldColumn{{Name: "description", Type: dmodel.FieldReference}}},
				},
				Tables: []*dmodel.Table{
					{Name: "tasks"},
				},
			},
		},
	}},
}

var PostgresSubqueryWithExistData = &postgresData{
	SQL: "SELECT name, EXISTS (SELECT * FROM tasks WHERE user_id = users.id) FROM users",
	Scopes: []*dmodel.StmtScope{{
		Name: "<root>",
		Fields: []*dmodel.Field{
			{Columns: []*dmodel.FieldColumn{{Name: "name", Type: dmodel.FieldReference}}},
			{Columns: []*dmodel.FieldColumn{{ReferenceName: "<field0>", Type: dmodel.FieldSubquery}}},
		},
		FieldScopes: []*dmodel.StmtScope{
			{
				Name: "<field0>",
				Fields: []*dmodel.Field{
					{Columns: []*dmodel.FieldColumn{{Type: dmodel.FieldStar}}},
					{Columns: []*dmodel.FieldColumn{
						{Name: "user_id", Type: dmodel.FieldCondition},
						{Table: "users", Name: "id", Type: dmodel.FieldCondition},
					}},
				},
				Tables: []*dmodel.Table{
					{Name: "tasks"},
				},
			},
		},
		Tables: []*dmodel.Table{
			{Name: "users"},
		},
	}},
}

var PostgresSubqueryWithInData = &postgresData{
	SQL: "SELECT user_id IN (SELECT user_id FROM tasks) FROM users",
	Scopes: []*dmodel.StmtScope{{
		Name: "<root>",
		Fields: []*dmodel.Field{
			{Columns: []*dmodel.FieldColumn{
				{Name: "user_id", Type: dmodel.FieldReference},
				{ReferenceName: "<field0>", Type: dmodel.FieldSubquery},
			}},
		},
		FieldScopes: []*dmodel.StmtScope{
			{
				Name: "<field0>",
				Fields: []*dmodel.Field{
					{Columns: []*dmodel.FieldColumn{{Name: "user_id", Type: dmodel.FieldReference}}},
				},
				Tables: []*dmodel.Table{
					{Name: "tasks"},
				},
			},
		},
		Tables: []*dmodel.Table{
			{Name: "users"},
		},
	}},
}

var PostgresSubqueryInSelectFunctionData = &postgresData{
	SQL: `
SELECT row_to_json(
  (
    SELECT t
    FROM (SELECT id, description FROM tasks LIMIT 1) AS t
  )
)
`,
	Scopes: []*dmodel.StmtScope{
		{
			Name: "<root>",
			Fields: []*dmodel.Field{
				{Columns: []*dmodel.FieldColumn{{ReferenceName: "<field0>", Type: dmodel.FieldSubquery}}},
			},
			FieldScopes: []*dmodel.StmtScope{
				{
					Name: "<field0>",
					Fields: []*dmodel.Field{
						{Columns: []*dmodel.FieldColumn{{Name: "t", Type: dmodel.FieldReference}}},
					},
					Tables: []*dmodel.Table{
						{Name: "<select0>", AsName: "t"},
					},
					SubScopes: []*dmodel.StmtScope{
						{
							Name: "<select0>",
							Fields: []*dmodel.Field{
								{Columns: []*dmodel.FieldColumn{{Name: "id", Type: dmodel.FieldReference}}},
								{Columns: []*dmodel.FieldColumn{{Name: "description", Type: dmodel.FieldReference}}},
							},
							Tables: []*dmodel.Table{
								{Name: "tasks"},
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
FROM (SELECT id, user_id FROM tasks)
`,
	Scopes: []*dmodel.StmtScope{
		{
			Name: "<root>",
			Fields: []*dmodel.Field{
				{Columns: []*dmodel.FieldColumn{{Type: dmodel.FieldStar}}},
			},
			Tables: []*dmodel.Table{
				{Name: "<select0>"},
			},
			SubScopes: []*dmodel.StmtScope{
				{
					Name: "<select0>",
					Fields: []*dmodel.Field{
						{Columns: []*dmodel.FieldColumn{{Name: "id", Type: dmodel.FieldReference}}},
						{Columns: []*dmodel.FieldColumn{{Name: "user_id", Type: dmodel.FieldReference}}},
					},
					Tables: []*dmodel.Table{
						{Name: "tasks"},
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
INNER JOIN (SELECT id, user_id FROM tasks) AS t ON users.id = t.user_id
`,
	Scopes: []*dmodel.StmtScope{
		{
			Name: "<root>",
			Fields: []*dmodel.Field{
				{Columns: []*dmodel.FieldColumn{{Name: "name", Type: dmodel.FieldReference}}},
				{Columns: []*dmodel.FieldColumn{{Table: "t", Name: "id", Type: dmodel.FieldReference}}},
				{Columns: []*dmodel.FieldColumn{{Name: "user_id", Type: dmodel.FieldReference}}},
				{Columns: []*dmodel.FieldColumn{
					{Table: "users", Name: "id", Type: dmodel.FieldCondition},
					{Table: "t", Name: "user_id", Type: dmodel.FieldCondition},
				}},
			},
			Tables: []*dmodel.Table{
				{Name: "users"},
				{Name: "<select0>", AsName: "t"},
			},
			SubScopes: []*dmodel.StmtScope{
				{
					Name: "<select0>",
					Fields: []*dmodel.Field{
						{Columns: []*dmodel.FieldColumn{{Name: "id", Type: dmodel.FieldReference}}},
						{Columns: []*dmodel.FieldColumn{{Name: "user_id", Type: dmodel.FieldReference}}},
					},
					Tables: []*dmodel.Table{
						{Name: "tasks"},
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
	(SELECT * FROM tasks) AS t
`,
	Scopes: []*dmodel.StmtScope{
		{
			Name: "<root>",
			Fields: []*dmodel.Field{
				{Columns: []*dmodel.FieldColumn{
					{Name: "description", Type: dmodel.FieldReference},
				}},
			},
			Tables: []*dmodel.Table{
				{Name: "<select0>", AsName: "t"},
			},
			SubScopes: []*dmodel.StmtScope{
				{
					Name: "<select0>",
					Fields: []*dmodel.Field{
						{Columns: []*dmodel.FieldColumn{{Type: dmodel.FieldStar}}},
					},
					Tables: []*dmodel.Table{
						{Name: "tasks"},
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
	LEFT OUTER JOIN LATERAL (SELECT user_id, description FROM tasks WHERE user_id = u.id) AS t ON u.id = t.user_id
`,
	Scopes: []*dmodel.StmtScope{
		{
			Name: "<root>",
			Fields: []*dmodel.Field{
				{Columns: []*dmodel.FieldColumn{
					{Name: "description", Type: dmodel.FieldReference},
				}},
				{Columns: []*dmodel.FieldColumn{
					{Table: "u", Name: "id", Type: dmodel.FieldCondition},
					{Table: "t", Name: "user_id", Type: dmodel.FieldCondition},
				}},
			},
			Tables: []*dmodel.Table{
				{Name: "<select0>", AsName: "u"},
				{Name: "<select1>", AsName: "t", IsLateral: true},
			},
			SubScopes: []*dmodel.StmtScope{
				{
					Name: "<select0>",
					Fields: []*dmodel.Field{
						{Columns: []*dmodel.FieldColumn{{Name: "id", Type: dmodel.FieldReference}}},
					},
					Tables: []*dmodel.Table{
						{Name: "users"},
					},
				},
				{
					Name: "<select1>",
					Fields: []*dmodel.Field{
						{Columns: []*dmodel.FieldColumn{{Name: "user_id", Type: dmodel.FieldReference}}},
						{Columns: []*dmodel.FieldColumn{{Name: "description", Type: dmodel.FieldReference}}},
						{Columns: []*dmodel.FieldColumn{
							{Name: "user_id", Type: dmodel.FieldCondition},
							{Table: "u", Name: "id", Type: dmodel.FieldCondition},
						}},
					},
					Tables: []*dmodel.Table{
						{Name: "tasks"},
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
	Scopes: []*dmodel.StmtScope{
		{
			Name: "<root>",
			Fields: []*dmodel.Field{
				{Columns: []*dmodel.FieldColumn{
					{Name: "b", Type: dmodel.FieldReference},
				}},
			},
			Tables: []*dmodel.Table{
				{Name: "<select0>", AsName: "a"},
				{Name: "<select1>", AsName: "b", IsLateral: true},
			},
			SubScopes: []*dmodel.StmtScope{
				{
					Name: "<select0>",
					Fields: []*dmodel.Field{
						{Columns: []*dmodel.FieldColumn{{Name: "name", Type: dmodel.FieldReference}}},
					},
					Tables: []*dmodel.Table{
						{Name: "users"},
					},
				},
				{
					Name: "<select1>",
					Fields: []*dmodel.Field{
						{Columns: []*dmodel.FieldColumn{{Name: "a", Type: dmodel.FieldReference}}},
					},
				},
			},
		},
	},
}

var PostgresSubqueryWithAggregateFunctionData = &postgresData{
	SQL: `
SELECT json_agg(t)
FROM (SELECT description FROM tasks LIMIT 3) AS t;
`,
	Scopes: []*dmodel.StmtScope{
		{
			Name: "<root>",
			Fields: []*dmodel.Field{
				{Columns: []*dmodel.FieldColumn{{Name: "t", Type: dmodel.FieldReference}}},
			},
			Tables: []*dmodel.Table{
				{AsName: "t", Name: "<select0>"},
			},
			SubScopes: []*dmodel.StmtScope{
				{
					Name: "<select0>",
					Fields: []*dmodel.Field{
						{Columns: []*dmodel.FieldColumn{{Name: "description", Type: dmodel.FieldReference}}},
					},
					Tables: []*dmodel.Table{
						{Name: "tasks"},
					},
				},
			},
		},
	},
}

var PostgresCorrelatedSubqueryData = &postgresData{
	SQL: `
SELECT description
FROM tasks
WHERE EXISTS(SELECT * FROM users WHERE users.id = user_id)
`,
	Scopes: []*dmodel.StmtScope{
		{
			Name: "<root>",
			Fields: []*dmodel.Field{
				{Columns: []*dmodel.FieldColumn{{Name: "description", Type: dmodel.FieldReference}}},
				{Columns: []*dmodel.FieldColumn{{ReferenceName: "<field0>", Type: dmodel.FieldSubquery}}},
			},
			FieldScopes: []*dmodel.StmtScope{
				{
					Name: "<field0>",
					Fields: []*dmodel.Field{
						{Columns: []*dmodel.FieldColumn{{Type: dmodel.FieldStar}}},
						{Columns: []*dmodel.FieldColumn{
							{Table: "users", Name: "id", Type: dmodel.FieldCondition},
							{Name: "user_id", Type: dmodel.FieldCondition},
						}},
					},
					Tables: []*dmodel.Table{
						{Name: "users"},
					},
				},
			},
			Tables: []*dmodel.Table{
				{Name: "tasks"},
			},
		},
	},
}
