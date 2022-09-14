package tdata

import "github.com/mrasu/GravityR/database/db_models"

type hasuraData struct {
	GQL    string
	SQL    string
	Scopes []*db_models.StmtScope
}

var Hasura1 = hasuraData{
	GQL: `
query MyQuery {
  users(where: {email: {_eq: "test33333@example.com"}}) {
    todos {
      title
      status
    }
    email
    name
  }
}`,
	SQL: `
SELECT
  coalesce(json_agg("root"), '[]') AS "root"
FROM
  (
    SELECT
      row_to_json(
        (
          SELECT
            "_e"
          FROM
            (
              SELECT
                "_root.ar.root.todos"."todos" AS "todos",
                "_root.base"."email" AS "email",
                "_root.base"."name" AS "name"
            ) AS "_e"
        )
      ) AS "root"
    FROM
      (
        SELECT
          *
        FROM
          "public"."users"
        WHERE
          (
            ("public"."users"."email") = (('test33333@example.com') :: varchar)
          )
      ) AS "_root.base"
      LEFT OUTER JOIN LATERAL (
        SELECT
          coalesce(json_agg("todos"), '[]') AS "todos"
        FROM
          (
            SELECT
              row_to_json(
                (
                  SELECT
                    "_e"
                  FROM
                    (
                      SELECT
                        "_root.ar.root.todos.base"."title" AS "title",
                        "_root.ar.root.todos.base"."status" AS "status"
                    ) AS "_e"
                )
              ) AS "todos"
            FROM
              (
                SELECT
                  *
                FROM
                  "public"."todos"
                WHERE
                  (("_root.base"."id") = ("user_id"))
              ) AS "_root.ar.root.todos.base"
          ) AS "_root.ar.root.todos"
      ) AS "_root.ar.root.todos" ON ('true')
  ) AS "_root"
`,
	Scopes: []*db_models.StmtScope{
		{
			Name: "<root>",
			Fields: []*db_models.Field{
				{
					AsName:  "root",
					Columns: []*db_models.FieldColumn{{Name: "root", Type: db_models.FieldReference}},
				},
			},
			Tables: []*db_models.Table{
				{AsName: "_root", Name: "<select0>"},
			},
			Scopes: []*db_models.StmtScope{
				{
					Name: "<select0>",
					Fields: []*db_models.Field{
						{
							AsName:  "root",
							Columns: []*db_models.FieldColumn{{ReferenceName: "<field0>", Type: db_models.FieldSubquery}},
						},
					},
					FieldScopes: []*db_models.StmtScope{
						{
							Name: "<field0>",
							Fields: []*db_models.Field{
								{Columns: []*db_models.FieldColumn{{Name: "_e", Type: db_models.FieldReference}}},
							},
							Tables: []*db_models.Table{{AsName: "_e", Name: "<select3>"}},
							Scopes: []*db_models.StmtScope{
								{
									Name: "<select3>",
									Fields: []*db_models.Field{
										{
											AsName:  "todos",
											Columns: []*db_models.FieldColumn{{Table: "_root.ar.root.todos", Name: "todos", Type: db_models.FieldReference}},
										},
										{
											AsName:  "email",
											Columns: []*db_models.FieldColumn{{Table: "_root.base", Name: "email", Type: db_models.FieldReference}},
										},
										{
											AsName:  "name",
											Columns: []*db_models.FieldColumn{{Table: "_root.base", Name: "name", Type: db_models.FieldReference}},
										},
									},
								},
							},
						},
					},
					Tables: []*db_models.Table{
						{AsName: "_root.base", Name: "<select1>"},
						{AsName: "_root.ar.root.todos", Name: "<select2>", IsLateral: true},
					},
					Scopes: []*db_models.StmtScope{
						{
							Name: "<select1>",
							Fields: []*db_models.Field{
								{Columns: []*db_models.FieldColumn{{Type: db_models.FieldStar}}},
								{Columns: []*db_models.FieldColumn{{Table: "users", Name: "email", Type: db_models.FieldCondition}}},
							},
							Tables: []*db_models.Table{{Name: "users"}},
						},
						{
							Name: "<select2>",
							Fields: []*db_models.Field{
								{
									AsName:  "todos",
									Columns: []*db_models.FieldColumn{{Name: "todos", Type: db_models.FieldReference}},
								},
							},
							Tables: []*db_models.Table{{AsName: "_root.ar.root.todos", Name: "<select4>"}},
							Scopes: []*db_models.StmtScope{
								{
									Name: "<select4>",
									Fields: []*db_models.Field{
										{AsName: "todos", Columns: []*db_models.FieldColumn{{ReferenceName: "<field1>", Type: db_models.FieldSubquery}}},
									},
									FieldScopes: []*db_models.StmtScope{
										{
											Name: "<field1>",
											Fields: []*db_models.Field{
												{Columns: []*db_models.FieldColumn{{Name: "_e", Type: db_models.FieldReference}}},
											},
											Tables: []*db_models.Table{{AsName: "_e", Name: "<select6>"}},
											Scopes: []*db_models.StmtScope{
												{
													Name: "<select6>",
													Fields: []*db_models.Field{
														{
															AsName:  "title",
															Columns: []*db_models.FieldColumn{{Table: "_root.ar.root.todos.base", Name: "title", Type: db_models.FieldReference}},
														},
														{
															AsName:  "status",
															Columns: []*db_models.FieldColumn{{Table: "_root.ar.root.todos.base", Name: "status", Type: db_models.FieldReference}},
														},
													},
												},
											},
										},
									},
									Tables: []*db_models.Table{{AsName: "_root.ar.root.todos.base", Name: "<select5>"}},
									Scopes: []*db_models.StmtScope{
										{
											Name: "<select5>",
											Fields: []*db_models.Field{
												{
													Columns: []*db_models.FieldColumn{{Type: db_models.FieldStar}},
												},
												{
													Columns: []*db_models.FieldColumn{
														{Table: "_root.base", Name: "id", Type: db_models.FieldCondition},
														{Name: "user_id", Type: db_models.FieldCondition},
													},
												},
											},
											Tables: []*db_models.Table{{Name: "todos"}},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	},
}

var Hasura2 = hasuraData{
	GQL: `query MyQuery($email: String) {
  todos(where: {user: {email: {_eq: $email}}}) {
    user {
      email
      id
      name
    }
    description
    id
    status
  }
}
`,
	SQL: `
SELECT
  COALESCE(JSON_AGG("root"), '[]') AS "root"
FROM
  (
    SELECT
      ROW_TO_JSON(
        (
          SELECT
            "_e"
          FROM
            (
              SELECT
                "_root.or.user"."user" AS "user",
                 "_root.base"."description" AS "description",
                 "_root.base"."id" AS "id",
                 "_root.base"."status" AS "status"
            ) AS "_e"
        )
      ) AS "root"
    FROM
      (
        SELECT
          *
        FROM
          "public"."todos"
        WHERE
          (
            EXISTS
              (
                SELECT 1
                FROM
                  "public"."users" AS "__be_0_users"
                WHERE
                  (
                    (
                      (
                        ("__be_0_users"."id") = ("public"."todos"."user_id")
                      ) AND
                      ('true')
                    ) AND
                    (
                      ('true') AND
                      (
                        (
                          (
                            (
                              ("__be_0_users"."email") = (('test1111@example.com')::varchar)
                            ) AND
                            ('true')
                          ) AND
                          ('true')
                        ) AND
                        ('true')
                      )
                    )
                  )
            )
          )
      ) AS "_root.base"
    LEFT OUTER JOIN LATERAL
      (
        SELECT row_to_json(
          (
            SELECT
              "_e"
            FROM
              (
                SELECT
                  "_root.or.user.base"."email" AS "email",
                  "_root.or.user.base"."id" AS "id",
                  "_root.or.user.base"."name" AS "name"
                ) AS "_e"
          )
        ) AS "user"
        FROM
          (
            SELECT
              *
            FROM
              "public"."users"
            WHERE
              (("_root.base"."user_id") = ("id"))
            LIMIT 1
          ) AS "_root.or.user.base"
      ) AS "_root.or.user" ON ('true')
  ) AS "_root"
`,

	Scopes: []*db_models.StmtScope{
		{
			Name: "<root>",
			Fields: []*db_models.Field{
				{
					AsName:  "root",
					Columns: []*db_models.FieldColumn{{Name: "root", Type: db_models.FieldReference}},
				},
			},
			Tables: []*db_models.Table{
				{AsName: "_root", Name: "<select0>"},
			},
			Scopes: []*db_models.StmtScope{
				{
					Name: "<select0>",
					Fields: []*db_models.Field{
						{
							AsName:  "root",
							Columns: []*db_models.FieldColumn{{ReferenceName: "<field0>", Type: db_models.FieldSubquery}},
						},
					},
					FieldScopes: []*db_models.StmtScope{
						{
							Name: "<field0>",
							Fields: []*db_models.Field{
								{Columns: []*db_models.FieldColumn{{Name: "_e", Type: db_models.FieldReference}}},
							},
							Tables: []*db_models.Table{{AsName: "_e", Name: "<select3>"}},
							Scopes: []*db_models.StmtScope{
								{
									Name: "<select3>",
									Fields: []*db_models.Field{
										{
											AsName:  "user",
											Columns: []*db_models.FieldColumn{{Table: "_root.or.user", Name: "user", Type: db_models.FieldReference}},
										},
										{
											AsName:  "description",
											Columns: []*db_models.FieldColumn{{Table: "_root.base", Name: "description", Type: db_models.FieldReference}},
										},
										{
											AsName:  "id",
											Columns: []*db_models.FieldColumn{{Table: "_root.base", Name: "id", Type: db_models.FieldReference}},
										},
										{
											AsName:  "status",
											Columns: []*db_models.FieldColumn{{Table: "_root.base", Name: "status", Type: db_models.FieldReference}},
										},
									},
								},
							},
						},
					},
					Tables: []*db_models.Table{
						{AsName: "_root.base", Name: "<select1>"},
						{AsName: "_root.or.user", Name: "<select2>", IsLateral: true},
					},
					Scopes: []*db_models.StmtScope{
						{
							Name: "<select1>",
							Fields: []*db_models.Field{
								{Columns: []*db_models.FieldColumn{{Type: db_models.FieldStar}}},
								{Columns: []*db_models.FieldColumn{{ReferenceName: "<field1>", Type: db_models.FieldSubquery}}},
							},
							FieldScopes: []*db_models.StmtScope{
								{
									Name: "<field1>",
									Fields: []*db_models.Field{
										{Columns: []*db_models.FieldColumn{
											{Table: "__be_0_users", Name: "id", Type: db_models.FieldCondition},
											{Table: "todos", Name: "user_id", Type: db_models.FieldCondition},
											{Table: "__be_0_users", Name: "email", Type: db_models.FieldCondition},
										}},
									},
									Tables: []*db_models.Table{{AsName: "__be_0_users", Name: "users"}},
								},
							},
							Tables: []*db_models.Table{{Name: "todos"}},
						},
						{
							Name: "<select2>",
							Fields: []*db_models.Field{
								{
									AsName:  "user",
									Columns: []*db_models.FieldColumn{{ReferenceName: "<field2>", Type: db_models.FieldSubquery}},
								},
							},
							FieldScopes: []*db_models.StmtScope{
								{
									Name: "<field2>",
									Fields: []*db_models.Field{
										{Columns: []*db_models.FieldColumn{{Name: "_e", Type: db_models.FieldReference}}},
									},
									Tables: []*db_models.Table{{AsName: "_e", Name: "<select5>"}},
									Scopes: []*db_models.StmtScope{
										{
											Name: "<select5>",
											Fields: []*db_models.Field{
												{
													AsName:  "email",
													Columns: []*db_models.FieldColumn{{Table: "_root.or.user.base", Name: "email", Type: db_models.FieldReference}},
												},
												{
													AsName:  "id",
													Columns: []*db_models.FieldColumn{{Table: "_root.or.user.base", Name: "id", Type: db_models.FieldReference}},
												},
												{
													AsName:  "name",
													Columns: []*db_models.FieldColumn{{Table: "_root.or.user.base", Name: "name", Type: db_models.FieldReference}},
												},
											},
										},
									},
								},
							},
							Tables: []*db_models.Table{{AsName: "_root.or.user.base", Name: "<select4>"}},
							Scopes: []*db_models.StmtScope{
								{
									Name: "<select4>",
									Fields: []*db_models.Field{
										{Columns: []*db_models.FieldColumn{{Type: db_models.FieldStar}}},
										{Columns: []*db_models.FieldColumn{
											{Table: "_root.base", Name: "user_id", Type: db_models.FieldCondition},
											{Name: "id", Type: db_models.FieldCondition},
										}},
									},
									Tables: []*db_models.Table{{Name: "users"}},
								},
							},
						},
					},
				},
			},
		},
	},
}

var Hasura3 = hasuraData{
	GQL: `query MyQuery($email: String) {
  users(where: {email: {_eq: $email}}) {
    email
    name
    todos_aggregate {
      aggregate {
        avg {
          status
        }
      }
    }
    todos(where: {}) {
      title
      description
    }
  }
}`,
	SQL: `SELECT
  coalesce(json_agg("root"), '[]') AS "root"
FROM
  (
    SELECT
      row_to_json(
        (
          SELECT
            "_e"
          FROM
            (
              SELECT
                "_root.base"."email" AS "email",
                "_root.base"."name" AS "name",
                "_root.ar.root.todos_aggregate"."todos_aggregate" AS "todos_aggregate",
                "_root.ar.root.todos"."todos" AS "todos"
            ) AS "_e"
        )
      ) AS "root"
    FROM
      (
        SELECT
          *
        FROM
          "public"."users"
        WHERE
          (
            ("public"."users"."email") = (('test33333@example.com') :: varchar)
          )
      ) AS "_root.base"
      LEFT OUTER JOIN LATERAL (
        SELECT
          json_build_object(
            'aggregate',
            json_build_object(
              'avg',
              json_build_object('status', avg("status"))
            )
          ) AS "todos_aggregate"
        FROM
          (
            SELECT
              "_root.ar.root.todos_aggregate.base"."status" AS "status"
            FROM
              (
                SELECT
                  *
                FROM
                  "public"."todos"
                WHERE
                  (("_root.base"."id") = ("user_id"))
              ) AS "_root.ar.root.todos_aggregate.base"
          ) AS "_root.ar.root.todos_aggregate"
      ) AS "_root.ar.root.todos_aggregate" ON ('true')
      LEFT OUTER JOIN LATERAL (
        SELECT
          coalesce(json_agg("todos"), '[]') AS "todos"
        FROM
          (
            SELECT
              row_to_json(
                (
                  SELECT
                    "_e"
                  FROM
                    (
                      SELECT
                        "_root.ar.root.todos.base"."title" AS "title",
                        "_root.ar.root.todos.base"."description" AS "description"
                    ) AS "_e"
                )
              ) AS "todos"
            FROM
              (
                SELECT
                  *
                FROM
                  "public"."todos"
                WHERE
                  (("_root.base"."id") = ("user_id"))
              ) AS "_root.ar.root.todos.base"
          ) AS "_root.ar.root.todos"
      ) AS "_root.ar.root.todos" ON ('true')
  ) AS "_root"`,
	Scopes: []*db_models.StmtScope{
		{
			Name: "<root>",
			Fields: []*db_models.Field{
				{
					AsName:  "root",
					Columns: []*db_models.FieldColumn{{Name: "root", Type: db_models.FieldReference}},
				},
			},
			Tables: []*db_models.Table{
				{AsName: "_root", Name: "<select0>"},
			},
			Scopes: []*db_models.StmtScope{
				{
					Name: "<select0>",
					Fields: []*db_models.Field{
						{
							AsName:  "root",
							Columns: []*db_models.FieldColumn{{ReferenceName: "<field0>", Type: db_models.FieldSubquery}},
						},
					},
					FieldScopes: []*db_models.StmtScope{
						{
							Name: "<field0>",
							Fields: []*db_models.Field{
								{Columns: []*db_models.FieldColumn{{Name: "_e", Type: db_models.FieldReference}}},
							},
							Tables: []*db_models.Table{{AsName: "_e", Name: "<select4>"}},
							Scopes: []*db_models.StmtScope{
								{
									Name: "<select4>",
									Fields: []*db_models.Field{
										{
											AsName:  "email",
											Columns: []*db_models.FieldColumn{{Table: "_root.base", Name: "email", Type: db_models.FieldReference}},
										},
										{
											AsName:  "name",
											Columns: []*db_models.FieldColumn{{Table: "_root.base", Name: "name", Type: db_models.FieldReference}},
										},
										{
											AsName:  "todos_aggregate",
											Columns: []*db_models.FieldColumn{{Table: "_root.ar.root.todos_aggregate", Name: "todos_aggregate", Type: db_models.FieldReference}},
										},
										{
											AsName:  "todos",
											Columns: []*db_models.FieldColumn{{Table: "_root.ar.root.todos", Name: "todos", Type: db_models.FieldReference}},
										},
									},
								},
							},
						},
					},
					Tables: []*db_models.Table{
						{AsName: "_root.base", Name: "<select1>"},
						{AsName: "_root.ar.root.todos_aggregate", Name: "<select2>", IsLateral: true},
						{AsName: "_root.ar.root.todos", Name: "<select3>", IsLateral: true},
					},
					Scopes: []*db_models.StmtScope{
						{
							Name: "<select1>",
							Fields: []*db_models.Field{
								{Columns: []*db_models.FieldColumn{{Type: db_models.FieldStar}}},
								{Columns: []*db_models.FieldColumn{{Table: "users", Name: "email", Type: db_models.FieldCondition}}},
							},
							Tables: []*db_models.Table{{Name: "users"}},
						},
						{
							Name: "<select2>",
							Fields: []*db_models.Field{
								{AsName: "todos_aggregate", Columns: []*db_models.FieldColumn{{Name: "status", Type: db_models.FieldReference}}},
							},
							Tables: []*db_models.Table{{AsName: "_root.ar.root.todos_aggregate", Name: "<select5>"}},
							Scopes: []*db_models.StmtScope{
								{
									Name: "<select5>",
									Fields: []*db_models.Field{
										{AsName: "status", Columns: []*db_models.FieldColumn{
											{Table: "_root.ar.root.todos_aggregate.base", Name: "status", Type: db_models.FieldReference},
										}},
									},
									Tables: []*db_models.Table{{AsName: "_root.ar.root.todos_aggregate.base", Name: "<select6>"}},
									Scopes: []*db_models.StmtScope{
										{
											Name: "<select6>",
											Fields: []*db_models.Field{
												{Columns: []*db_models.FieldColumn{{Type: db_models.FieldStar}}},
												{Columns: []*db_models.FieldColumn{
													{Table: "_root.base", Name: "id", Type: db_models.FieldCondition},
													{Name: "user_id", Type: db_models.FieldCondition},
												}},
											},
											Tables: []*db_models.Table{{Name: "todos"}},
										},
									},
								},
							},
						},
						{
							Name: "<select3>",
							Fields: []*db_models.Field{
								{
									AsName:  "todos",
									Columns: []*db_models.FieldColumn{{Name: "todos", Type: db_models.FieldReference}},
								},
							},
							Tables: []*db_models.Table{{AsName: "_root.ar.root.todos", Name: "<select7>"}},
							Scopes: []*db_models.StmtScope{
								{
									Name: "<select7>",
									Fields: []*db_models.Field{
										{AsName: "todos", Columns: []*db_models.FieldColumn{{ReferenceName: "<field1>", Type: db_models.FieldSubquery}}},
									},
									FieldScopes: []*db_models.StmtScope{
										{
											Name: "<field1>",
											Fields: []*db_models.Field{
												{Columns: []*db_models.FieldColumn{{Name: "_e", Type: db_models.FieldReference}}},
											},
											Tables: []*db_models.Table{{AsName: "_e", Name: "<select9>"}},
											Scopes: []*db_models.StmtScope{
												{
													Name: "<select9>",
													Fields: []*db_models.Field{
														{
															AsName:  "title",
															Columns: []*db_models.FieldColumn{{Table: "_root.ar.root.todos.base", Name: "title", Type: db_models.FieldReference}},
														},
														{
															AsName:  "description",
															Columns: []*db_models.FieldColumn{{Table: "_root.ar.root.todos.base", Name: "description", Type: db_models.FieldReference}},
														},
													},
												},
											},
										},
									},
									Tables: []*db_models.Table{{AsName: "_root.ar.root.todos.base", Name: "<select8>"}},
									Scopes: []*db_models.StmtScope{
										{
											Name: "<select8>",
											Fields: []*db_models.Field{
												{Columns: []*db_models.FieldColumn{{Type: db_models.FieldStar}}},
												{Columns: []*db_models.FieldColumn{
													{Table: "_root.base", Name: "id", Type: db_models.FieldCondition},
													{Name: "user_id", Type: db_models.FieldCondition},
												}},
											},
											Tables: []*db_models.Table{{Name: "todos"}},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	},
}
