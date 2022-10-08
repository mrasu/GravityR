package tdata

import (
	"github.com/mrasu/GravityR/database/dservice/dparser"
)

type postgresData struct {
	SQL    string
	Scopes []*dparser.StmtScope
}

var PostgresSubqueryInSelectData = &postgresData{
	SQL: "SELECT (SELECT user_id FROM tasks limit 1)",
	Scopes: []*dparser.StmtScope{
		{
			Name: "<root>",
			Fields: []*dparser.Field{
				{Columns: []*dparser.FieldColumn{{ReferenceName: "<field0>", Type: dparser.FieldSubquery}}},
			},
			FieldScopes: []*dparser.StmtScope{
				{
					Name: "<field0>",
					Fields: []*dparser.Field{
						{Columns: []*dparser.FieldColumn{{Name: "user_id", Type: dparser.FieldReference}}},
					},
					Tables: []*dparser.Table{
						{Name: "tasks"},
					},
				},
			},
		},
	},
}

var PostgresSubqueryInSelectAliasedFieldData = &postgresData{
	SQL: "SELECT name, is_admin, (SELECT COUNT(status) FROM tasks) AS status_count FROM users",
	Scopes: []*dparser.StmtScope{{
		Name: "<root>",
		Fields: []*dparser.Field{
			{Columns: []*dparser.FieldColumn{{Name: "name", Type: dparser.FieldReference}}},
			{Columns: []*dparser.FieldColumn{{Name: "is_admin", Type: dparser.FieldReference}}},
			{AsName: "status_count", Columns: []*dparser.FieldColumn{{ReferenceName: "<field0>", Type: dparser.FieldSubquery}}},
		},
		FieldScopes: []*dparser.StmtScope{
			{
				Name: "<field0>",
				Fields: []*dparser.Field{
					{Columns: []*dparser.FieldColumn{{Name: "status", Type: dparser.FieldReference}}},
				},
				Tables: []*dparser.Table{
					{Name: "tasks"},
				},
			},
		},
		Tables: []*dparser.Table{
			{Name: "users"},
		},
	}},
}

var PostgresSubqueryWithComparisonData = &postgresData{
	SQL: "SELECT (SELECT COUNT(status) FROM tasks) - (SELECT COUNT(description) FROM tasks) AS no_desc_count",
	Scopes: []*dparser.StmtScope{{
		Name: "<root>",
		Fields: []*dparser.Field{
			{AsName: "no_desc_count", Columns: []*dparser.FieldColumn{
				{ReferenceName: "<field0>", Type: dparser.FieldSubquery},
				{ReferenceName: "<field1>", Type: dparser.FieldSubquery},
			}},
		},
		FieldScopes: []*dparser.StmtScope{
			{
				Name: "<field0>",
				Fields: []*dparser.Field{
					{Columns: []*dparser.FieldColumn{{Name: "status", Type: dparser.FieldReference}}},
				},
				Tables: []*dparser.Table{
					{Name: "tasks"},
				},
			},
			{
				Name: "<field1>",
				Fields: []*dparser.Field{
					{Columns: []*dparser.FieldColumn{{Name: "description", Type: dparser.FieldReference}}},
				},
				Tables: []*dparser.Table{
					{Name: "tasks"},
				},
			},
		},
	}},
}

var PostgresSubqueryWithExistData = &postgresData{
	SQL: "SELECT name, EXISTS (SELECT * FROM tasks WHERE user_id = users.id) FROM users",
	Scopes: []*dparser.StmtScope{{
		Name: "<root>",
		Fields: []*dparser.Field{
			{Columns: []*dparser.FieldColumn{{Name: "name", Type: dparser.FieldReference}}},
			{Columns: []*dparser.FieldColumn{{ReferenceName: "<field0>", Type: dparser.FieldSubquery}}},
		},
		FieldScopes: []*dparser.StmtScope{
			{
				Name: "<field0>",
				Fields: []*dparser.Field{
					{Columns: []*dparser.FieldColumn{{Type: dparser.FieldStar}}},
					{Columns: []*dparser.FieldColumn{
						{Name: "user_id", Type: dparser.FieldCondition},
						{Table: "users", Name: "id", Type: dparser.FieldCondition},
					}},
				},
				Tables: []*dparser.Table{
					{Name: "tasks"},
				},
			},
		},
		Tables: []*dparser.Table{
			{Name: "users"},
		},
	}},
}

var PostgresSubqueryWithInData = &postgresData{
	SQL: "SELECT user_id IN (SELECT user_id FROM tasks) FROM users",
	Scopes: []*dparser.StmtScope{{
		Name: "<root>",
		Fields: []*dparser.Field{
			{Columns: []*dparser.FieldColumn{
				{Name: "user_id", Type: dparser.FieldReference},
				{ReferenceName: "<field0>", Type: dparser.FieldSubquery},
			}},
		},
		FieldScopes: []*dparser.StmtScope{
			{
				Name: "<field0>",
				Fields: []*dparser.Field{
					{Columns: []*dparser.FieldColumn{{Name: "user_id", Type: dparser.FieldReference}}},
				},
				Tables: []*dparser.Table{
					{Name: "tasks"},
				},
			},
		},
		Tables: []*dparser.Table{
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
	Scopes: []*dparser.StmtScope{
		{
			Name: "<root>",
			Fields: []*dparser.Field{
				{Columns: []*dparser.FieldColumn{{ReferenceName: "<field0>", Type: dparser.FieldSubquery}}},
			},
			FieldScopes: []*dparser.StmtScope{
				{
					Name: "<field0>",
					Fields: []*dparser.Field{
						{Columns: []*dparser.FieldColumn{{Name: "t", Type: dparser.FieldReference}}},
					},
					Tables: []*dparser.Table{
						{Name: "<select0>", AsName: "t"},
					},
					SubScopes: []*dparser.StmtScope{
						{
							Name: "<select0>",
							Fields: []*dparser.Field{
								{Columns: []*dparser.FieldColumn{{Name: "id", Type: dparser.FieldReference}}},
								{Columns: []*dparser.FieldColumn{{Name: "description", Type: dparser.FieldReference}}},
							},
							Tables: []*dparser.Table{
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
	Scopes: []*dparser.StmtScope{
		{
			Name: "<root>",
			Fields: []*dparser.Field{
				{Columns: []*dparser.FieldColumn{{Type: dparser.FieldStar}}},
			},
			Tables: []*dparser.Table{
				{Name: "<select0>"},
			},
			SubScopes: []*dparser.StmtScope{
				{
					Name: "<select0>",
					Fields: []*dparser.Field{
						{Columns: []*dparser.FieldColumn{{Name: "id", Type: dparser.FieldReference}}},
						{Columns: []*dparser.FieldColumn{{Name: "user_id", Type: dparser.FieldReference}}},
					},
					Tables: []*dparser.Table{
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
	Scopes: []*dparser.StmtScope{
		{
			Name: "<root>",
			Fields: []*dparser.Field{
				{Columns: []*dparser.FieldColumn{{Name: "name", Type: dparser.FieldReference}}},
				{Columns: []*dparser.FieldColumn{{Table: "t", Name: "id", Type: dparser.FieldReference}}},
				{Columns: []*dparser.FieldColumn{{Name: "user_id", Type: dparser.FieldReference}}},
				{Columns: []*dparser.FieldColumn{
					{Table: "users", Name: "id", Type: dparser.FieldCondition},
					{Table: "t", Name: "user_id", Type: dparser.FieldCondition},
				}},
			},
			Tables: []*dparser.Table{
				{Name: "users"},
				{Name: "<select0>", AsName: "t"},
			},
			SubScopes: []*dparser.StmtScope{
				{
					Name: "<select0>",
					Fields: []*dparser.Field{
						{Columns: []*dparser.FieldColumn{{Name: "id", Type: dparser.FieldReference}}},
						{Columns: []*dparser.FieldColumn{{Name: "user_id", Type: dparser.FieldReference}}},
					},
					Tables: []*dparser.Table{
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
	Scopes: []*dparser.StmtScope{
		{
			Name: "<root>",
			Fields: []*dparser.Field{
				{Columns: []*dparser.FieldColumn{
					{Name: "description", Type: dparser.FieldReference},
				}},
			},
			Tables: []*dparser.Table{
				{Name: "<select0>", AsName: "t"},
			},
			SubScopes: []*dparser.StmtScope{
				{
					Name: "<select0>",
					Fields: []*dparser.Field{
						{Columns: []*dparser.FieldColumn{{Type: dparser.FieldStar}}},
					},
					Tables: []*dparser.Table{
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
	Scopes: []*dparser.StmtScope{
		{
			Name: "<root>",
			Fields: []*dparser.Field{
				{Columns: []*dparser.FieldColumn{
					{Name: "description", Type: dparser.FieldReference},
				}},
				{Columns: []*dparser.FieldColumn{
					{Table: "u", Name: "id", Type: dparser.FieldCondition},
					{Table: "t", Name: "user_id", Type: dparser.FieldCondition},
				}},
			},
			Tables: []*dparser.Table{
				{Name: "<select0>", AsName: "u"},
				{Name: "<select1>", AsName: "t", IsLateral: true},
			},
			SubScopes: []*dparser.StmtScope{
				{
					Name: "<select0>",
					Fields: []*dparser.Field{
						{Columns: []*dparser.FieldColumn{{Name: "id", Type: dparser.FieldReference}}},
					},
					Tables: []*dparser.Table{
						{Name: "users"},
					},
				},
				{
					Name: "<select1>",
					Fields: []*dparser.Field{
						{Columns: []*dparser.FieldColumn{{Name: "user_id", Type: dparser.FieldReference}}},
						{Columns: []*dparser.FieldColumn{{Name: "description", Type: dparser.FieldReference}}},
						{Columns: []*dparser.FieldColumn{
							{Name: "user_id", Type: dparser.FieldCondition},
							{Table: "u", Name: "id", Type: dparser.FieldCondition},
						}},
					},
					Tables: []*dparser.Table{
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
	Scopes: []*dparser.StmtScope{
		{
			Name: "<root>",
			Fields: []*dparser.Field{
				{Columns: []*dparser.FieldColumn{
					{Name: "b", Type: dparser.FieldReference},
				}},
			},
			Tables: []*dparser.Table{
				{Name: "<select0>", AsName: "a"},
				{Name: "<select1>", AsName: "b", IsLateral: true},
			},
			SubScopes: []*dparser.StmtScope{
				{
					Name: "<select0>",
					Fields: []*dparser.Field{
						{Columns: []*dparser.FieldColumn{{Name: "name", Type: dparser.FieldReference}}},
					},
					Tables: []*dparser.Table{
						{Name: "users"},
					},
				},
				{
					Name: "<select1>",
					Fields: []*dparser.Field{
						{Columns: []*dparser.FieldColumn{{Name: "a", Type: dparser.FieldReference}}},
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
	Scopes: []*dparser.StmtScope{
		{
			Name: "<root>",
			Fields: []*dparser.Field{
				{Columns: []*dparser.FieldColumn{{Name: "t", Type: dparser.FieldReference}}},
			},
			Tables: []*dparser.Table{
				{AsName: "t", Name: "<select0>"},
			},
			SubScopes: []*dparser.StmtScope{
				{
					Name: "<select0>",
					Fields: []*dparser.Field{
						{Columns: []*dparser.FieldColumn{{Name: "description", Type: dparser.FieldReference}}},
					},
					Tables: []*dparser.Table{
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
	Scopes: []*dparser.StmtScope{
		{
			Name: "<root>",
			Fields: []*dparser.Field{
				{Columns: []*dparser.FieldColumn{{Name: "description", Type: dparser.FieldReference}}},
				{Columns: []*dparser.FieldColumn{{ReferenceName: "<field0>", Type: dparser.FieldSubquery}}},
			},
			FieldScopes: []*dparser.StmtScope{
				{
					Name: "<field0>",
					Fields: []*dparser.Field{
						{Columns: []*dparser.FieldColumn{{Type: dparser.FieldStar}}},
						{Columns: []*dparser.FieldColumn{
							{Table: "users", Name: "id", Type: dparser.FieldCondition},
							{Name: "user_id", Type: dparser.FieldCondition},
						}},
					},
					Tables: []*dparser.Table{
						{Name: "users"},
					},
				},
			},
			Tables: []*dparser.Table{
				{Name: "tasks"},
			},
		},
	},
}
