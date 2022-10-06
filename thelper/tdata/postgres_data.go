package tdata

import "github.com/mrasu/GravityR/database"

type postgresData struct {
	SQL    string
	Scopes []*database.StmtScope
}

var PostgresSubqueryInSelectData = &postgresData{
	SQL: "SELECT (SELECT user_id FROM tasks limit 1)",
	Scopes: []*database.StmtScope{
		{
			Name: "<root>",
			Fields: []*database.Field{
				{Columns: []*database.FieldColumn{{ReferenceName: "<field0>", Type: database.FieldSubquery}}},
			},
			FieldScopes: []*database.StmtScope{
				{
					Name: "<field0>",
					Fields: []*database.Field{
						{Columns: []*database.FieldColumn{{Name: "user_id", Type: database.FieldReference}}},
					},
					Tables: []*database.Table{
						{Name: "tasks"},
					},
				},
			},
		},
	},
}

var PostgresSubqueryInSelectAliasedFieldData = &postgresData{
	SQL: "SELECT name, is_admin, (SELECT COUNT(status) FROM tasks) AS status_count FROM users",
	Scopes: []*database.StmtScope{{
		Name: "<root>",
		Fields: []*database.Field{
			{Columns: []*database.FieldColumn{{Name: "name", Type: database.FieldReference}}},
			{Columns: []*database.FieldColumn{{Name: "is_admin", Type: database.FieldReference}}},
			{AsName: "status_count", Columns: []*database.FieldColumn{{ReferenceName: "<field0>", Type: database.FieldSubquery}}},
		},
		FieldScopes: []*database.StmtScope{
			{
				Name: "<field0>",
				Fields: []*database.Field{
					{Columns: []*database.FieldColumn{{Name: "status", Type: database.FieldReference}}},
				},
				Tables: []*database.Table{
					{Name: "tasks"},
				},
			},
		},
		Tables: []*database.Table{
			{Name: "users"},
		},
	}},
}

var PostgresSubqueryWithComparisonData = &postgresData{
	SQL: "SELECT (SELECT COUNT(status) FROM tasks) - (SELECT COUNT(description) FROM tasks) AS no_desc_count",
	Scopes: []*database.StmtScope{{
		Name: "<root>",
		Fields: []*database.Field{
			{AsName: "no_desc_count", Columns: []*database.FieldColumn{
				{ReferenceName: "<field0>", Type: database.FieldSubquery},
				{ReferenceName: "<field1>", Type: database.FieldSubquery},
			}},
		},
		FieldScopes: []*database.StmtScope{
			{
				Name: "<field0>",
				Fields: []*database.Field{
					{Columns: []*database.FieldColumn{{Name: "status", Type: database.FieldReference}}},
				},
				Tables: []*database.Table{
					{Name: "tasks"},
				},
			},
			{
				Name: "<field1>",
				Fields: []*database.Field{
					{Columns: []*database.FieldColumn{{Name: "description", Type: database.FieldReference}}},
				},
				Tables: []*database.Table{
					{Name: "tasks"},
				},
			},
		},
	}},
}

var PostgresSubqueryWithExistData = &postgresData{
	SQL: "SELECT name, EXISTS (SELECT * FROM tasks WHERE user_id = users.id) FROM users",
	Scopes: []*database.StmtScope{{
		Name: "<root>",
		Fields: []*database.Field{
			{Columns: []*database.FieldColumn{{Name: "name", Type: database.FieldReference}}},
			{Columns: []*database.FieldColumn{{ReferenceName: "<field0>", Type: database.FieldSubquery}}},
		},
		FieldScopes: []*database.StmtScope{
			{
				Name: "<field0>",
				Fields: []*database.Field{
					{Columns: []*database.FieldColumn{{Type: database.FieldStar}}},
					{Columns: []*database.FieldColumn{
						{Name: "user_id", Type: database.FieldCondition},
						{Table: "users", Name: "id", Type: database.FieldCondition},
					}},
				},
				Tables: []*database.Table{
					{Name: "tasks"},
				},
			},
		},
		Tables: []*database.Table{
			{Name: "users"},
		},
	}},
}

var PostgresSubqueryWithInData = &postgresData{
	SQL: "SELECT user_id IN (SELECT user_id FROM tasks) FROM users",
	Scopes: []*database.StmtScope{{
		Name: "<root>",
		Fields: []*database.Field{
			{Columns: []*database.FieldColumn{
				{Name: "user_id", Type: database.FieldReference},
				{ReferenceName: "<field0>", Type: database.FieldSubquery},
			}},
		},
		FieldScopes: []*database.StmtScope{
			{
				Name: "<field0>",
				Fields: []*database.Field{
					{Columns: []*database.FieldColumn{{Name: "user_id", Type: database.FieldReference}}},
				},
				Tables: []*database.Table{
					{Name: "tasks"},
				},
			},
		},
		Tables: []*database.Table{
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
	Scopes: []*database.StmtScope{
		{
			Name: "<root>",
			Fields: []*database.Field{
				{Columns: []*database.FieldColumn{{ReferenceName: "<field0>", Type: database.FieldSubquery}}},
			},
			FieldScopes: []*database.StmtScope{
				{
					Name: "<field0>",
					Fields: []*database.Field{
						{Columns: []*database.FieldColumn{{Name: "t", Type: database.FieldReference}}},
					},
					Tables: []*database.Table{
						{Name: "<select0>", AsName: "t"},
					},
					SubScopes: []*database.StmtScope{
						{
							Name: "<select0>",
							Fields: []*database.Field{
								{Columns: []*database.FieldColumn{{Name: "id", Type: database.FieldReference}}},
								{Columns: []*database.FieldColumn{{Name: "description", Type: database.FieldReference}}},
							},
							Tables: []*database.Table{
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
	Scopes: []*database.StmtScope{
		{
			Name: "<root>",
			Fields: []*database.Field{
				{Columns: []*database.FieldColumn{{Type: database.FieldStar}}},
			},
			Tables: []*database.Table{
				{Name: "<select0>"},
			},
			SubScopes: []*database.StmtScope{
				{
					Name: "<select0>",
					Fields: []*database.Field{
						{Columns: []*database.FieldColumn{{Name: "id", Type: database.FieldReference}}},
						{Columns: []*database.FieldColumn{{Name: "user_id", Type: database.FieldReference}}},
					},
					Tables: []*database.Table{
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
	Scopes: []*database.StmtScope{
		{
			Name: "<root>",
			Fields: []*database.Field{
				{Columns: []*database.FieldColumn{{Name: "name", Type: database.FieldReference}}},
				{Columns: []*database.FieldColumn{{Table: "t", Name: "id", Type: database.FieldReference}}},
				{Columns: []*database.FieldColumn{{Name: "user_id", Type: database.FieldReference}}},
				{Columns: []*database.FieldColumn{
					{Table: "users", Name: "id", Type: database.FieldCondition},
					{Table: "t", Name: "user_id", Type: database.FieldCondition},
				}},
			},
			Tables: []*database.Table{
				{Name: "users"},
				{Name: "<select0>", AsName: "t"},
			},
			SubScopes: []*database.StmtScope{
				{
					Name: "<select0>",
					Fields: []*database.Field{
						{Columns: []*database.FieldColumn{{Name: "id", Type: database.FieldReference}}},
						{Columns: []*database.FieldColumn{{Name: "user_id", Type: database.FieldReference}}},
					},
					Tables: []*database.Table{
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
	Scopes: []*database.StmtScope{
		{
			Name: "<root>",
			Fields: []*database.Field{
				{Columns: []*database.FieldColumn{
					{Name: "description", Type: database.FieldReference},
				}},
			},
			Tables: []*database.Table{
				{Name: "<select0>", AsName: "t"},
			},
			SubScopes: []*database.StmtScope{
				{
					Name: "<select0>",
					Fields: []*database.Field{
						{Columns: []*database.FieldColumn{{Type: database.FieldStar}}},
					},
					Tables: []*database.Table{
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
	Scopes: []*database.StmtScope{
		{
			Name: "<root>",
			Fields: []*database.Field{
				{Columns: []*database.FieldColumn{
					{Name: "description", Type: database.FieldReference},
				}},
				{Columns: []*database.FieldColumn{
					{Table: "u", Name: "id", Type: database.FieldCondition},
					{Table: "t", Name: "user_id", Type: database.FieldCondition},
				}},
			},
			Tables: []*database.Table{
				{Name: "<select0>", AsName: "u"},
				{Name: "<select1>", AsName: "t", IsLateral: true},
			},
			SubScopes: []*database.StmtScope{
				{
					Name: "<select0>",
					Fields: []*database.Field{
						{Columns: []*database.FieldColumn{{Name: "id", Type: database.FieldReference}}},
					},
					Tables: []*database.Table{
						{Name: "users"},
					},
				},
				{
					Name: "<select1>",
					Fields: []*database.Field{
						{Columns: []*database.FieldColumn{{Name: "user_id", Type: database.FieldReference}}},
						{Columns: []*database.FieldColumn{{Name: "description", Type: database.FieldReference}}},
						{Columns: []*database.FieldColumn{
							{Name: "user_id", Type: database.FieldCondition},
							{Table: "u", Name: "id", Type: database.FieldCondition},
						}},
					},
					Tables: []*database.Table{
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
	Scopes: []*database.StmtScope{
		{
			Name: "<root>",
			Fields: []*database.Field{
				{Columns: []*database.FieldColumn{
					{Name: "b", Type: database.FieldReference},
				}},
			},
			Tables: []*database.Table{
				{Name: "<select0>", AsName: "a"},
				{Name: "<select1>", AsName: "b", IsLateral: true},
			},
			SubScopes: []*database.StmtScope{
				{
					Name: "<select0>",
					Fields: []*database.Field{
						{Columns: []*database.FieldColumn{{Name: "name", Type: database.FieldReference}}},
					},
					Tables: []*database.Table{
						{Name: "users"},
					},
				},
				{
					Name: "<select1>",
					Fields: []*database.Field{
						{Columns: []*database.FieldColumn{{Name: "a", Type: database.FieldReference}}},
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
	Scopes: []*database.StmtScope{
		{
			Name: "<root>",
			Fields: []*database.Field{
				{Columns: []*database.FieldColumn{{Name: "t", Type: database.FieldReference}}},
			},
			Tables: []*database.Table{
				{AsName: "t", Name: "<select0>"},
			},
			SubScopes: []*database.StmtScope{
				{
					Name: "<select0>",
					Fields: []*database.Field{
						{Columns: []*database.FieldColumn{{Name: "description", Type: database.FieldReference}}},
					},
					Tables: []*database.Table{
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
	Scopes: []*database.StmtScope{
		{
			Name: "<root>",
			Fields: []*database.Field{
				{Columns: []*database.FieldColumn{{Name: "description", Type: database.FieldReference}}},
				{Columns: []*database.FieldColumn{{ReferenceName: "<field0>", Type: database.FieldSubquery}}},
			},
			FieldScopes: []*database.StmtScope{
				{
					Name: "<field0>",
					Fields: []*database.Field{
						{Columns: []*database.FieldColumn{{Type: database.FieldStar}}},
						{Columns: []*database.FieldColumn{
							{Table: "users", Name: "id", Type: database.FieldCondition},
							{Name: "user_id", Type: database.FieldCondition},
						}},
					},
					Tables: []*database.Table{
						{Name: "users"},
					},
				},
			},
			Tables: []*database.Table{
				{Name: "tasks"},
			},
		},
	},
}
