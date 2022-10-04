package tdata

import "github.com/mrasu/GravityR/database/common_model"

type hasuraData struct {
	GQL    string
	SQL    string
	Scopes []*common_model.StmtScope
}

var HasuraSubqueryWithWhereData = &hasuraData{
	GQL: `
query MyQuery {
  users(where: {email: {_eq: "test33333@example.com"}}) {
    tasks {
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
                "_root.ar.root.tasks"."tasks" AS "tasks",
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
          coalesce(json_agg("tasks"), '[]') AS "tasks"
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
                        "_root.ar.root.tasks.base"."title" AS "title",
                        "_root.ar.root.tasks.base"."status" AS "status"
                    ) AS "_e"
                )
              ) AS "tasks"
            FROM
              (
                SELECT
                  *
                FROM
                  "public"."tasks"
                WHERE
                  (("_root.base"."id") = ("user_id"))
              ) AS "_root.ar.root.tasks.base"
          ) AS "_root.ar.root.tasks"
      ) AS "_root.ar.root.tasks" ON ('true')
  ) AS "_root"
`,
	Scopes: []*common_model.StmtScope{
		{
			Name: "<root>",
			Fields: []*common_model.Field{
				{
					AsName:  "root",
					Columns: []*common_model.FieldColumn{{Name: "root", Type: common_model.FieldReference}},
				},
			},
			Tables: []*common_model.Table{
				{AsName: "_root", Name: "<select0>"},
			},
			SubScopes: []*common_model.StmtScope{
				{
					Name: "<select0>",
					Fields: []*common_model.Field{
						{
							AsName:  "root",
							Columns: []*common_model.FieldColumn{{ReferenceName: "<field0>", Type: common_model.FieldSubquery}},
						},
					},
					FieldScopes: []*common_model.StmtScope{
						{
							Name: "<field0>",
							Fields: []*common_model.Field{
								{Columns: []*common_model.FieldColumn{{Name: "_e", Type: common_model.FieldReference}}},
							},
							Tables: []*common_model.Table{{AsName: "_e", Name: "<select3>"}},
							SubScopes: []*common_model.StmtScope{
								{
									Name: "<select3>",
									Fields: []*common_model.Field{
										{
											AsName:  "tasks",
											Columns: []*common_model.FieldColumn{{Table: "_root.ar.root.tasks", Name: "tasks", Type: common_model.FieldReference}},
										},
										{
											AsName:  "email",
											Columns: []*common_model.FieldColumn{{Table: "_root.base", Name: "email", Type: common_model.FieldReference}},
										},
										{
											AsName:  "name",
											Columns: []*common_model.FieldColumn{{Table: "_root.base", Name: "name", Type: common_model.FieldReference}},
										},
									},
								},
							},
						},
					},
					Tables: []*common_model.Table{
						{AsName: "_root.base", Name: "<select1>"},
						{AsName: "_root.ar.root.tasks", Name: "<select2>", IsLateral: true},
					},
					SubScopes: []*common_model.StmtScope{
						{
							Name: "<select1>",
							Fields: []*common_model.Field{
								{Columns: []*common_model.FieldColumn{{Type: common_model.FieldStar}}},
								{Columns: []*common_model.FieldColumn{{Table: "users", Name: "email", Type: common_model.FieldCondition}}},
							},
							Tables: []*common_model.Table{{Name: "users"}},
						},
						{
							Name: "<select2>",
							Fields: []*common_model.Field{
								{
									AsName:  "tasks",
									Columns: []*common_model.FieldColumn{{Name: "tasks", Type: common_model.FieldReference}},
								},
							},
							Tables: []*common_model.Table{{AsName: "_root.ar.root.tasks", Name: "<select4>"}},
							SubScopes: []*common_model.StmtScope{
								{
									Name: "<select4>",
									Fields: []*common_model.Field{
										{AsName: "tasks", Columns: []*common_model.FieldColumn{{ReferenceName: "<field1>", Type: common_model.FieldSubquery}}},
									},
									FieldScopes: []*common_model.StmtScope{
										{
											Name: "<field1>",
											Fields: []*common_model.Field{
												{Columns: []*common_model.FieldColumn{{Name: "_e", Type: common_model.FieldReference}}},
											},
											Tables: []*common_model.Table{{AsName: "_e", Name: "<select6>"}},
											SubScopes: []*common_model.StmtScope{
												{
													Name: "<select6>",
													Fields: []*common_model.Field{
														{
															AsName:  "title",
															Columns: []*common_model.FieldColumn{{Table: "_root.ar.root.tasks.base", Name: "title", Type: common_model.FieldReference}},
														},
														{
															AsName:  "status",
															Columns: []*common_model.FieldColumn{{Table: "_root.ar.root.tasks.base", Name: "status", Type: common_model.FieldReference}},
														},
													},
												},
											},
										},
									},
									Tables: []*common_model.Table{{AsName: "_root.ar.root.tasks.base", Name: "<select5>"}},
									SubScopes: []*common_model.StmtScope{
										{
											Name: "<select5>",
											Fields: []*common_model.Field{
												{
													Columns: []*common_model.FieldColumn{{Type: common_model.FieldStar}},
												},
												{
													Columns: []*common_model.FieldColumn{
														{Table: "_root.base", Name: "id", Type: common_model.FieldCondition},
														{Name: "user_id", Type: common_model.FieldCondition},
													},
												},
											},
											Tables: []*common_model.Table{{Name: "tasks"}},
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

var HasuraSubqueryWithVariablesData = &hasuraData{
	GQL: `query MyQuery($email: String) {
  tasks(where: {user: {email: {_eq: $email}}}) {
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
          "public"."tasks"
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
                        ("__be_0_users"."id") = ("public"."tasks"."user_id")
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
	Scopes: []*common_model.StmtScope{
		{
			Name: "<root>",
			Fields: []*common_model.Field{
				{
					AsName:  "root",
					Columns: []*common_model.FieldColumn{{Name: "root", Type: common_model.FieldReference}},
				},
			},
			Tables: []*common_model.Table{
				{AsName: "_root", Name: "<select0>"},
			},
			SubScopes: []*common_model.StmtScope{
				{
					Name: "<select0>",
					Fields: []*common_model.Field{
						{
							AsName:  "root",
							Columns: []*common_model.FieldColumn{{ReferenceName: "<field0>", Type: common_model.FieldSubquery}},
						},
					},
					FieldScopes: []*common_model.StmtScope{
						{
							Name: "<field0>",
							Fields: []*common_model.Field{
								{Columns: []*common_model.FieldColumn{{Name: "_e", Type: common_model.FieldReference}}},
							},
							Tables: []*common_model.Table{{AsName: "_e", Name: "<select3>"}},
							SubScopes: []*common_model.StmtScope{
								{
									Name: "<select3>",
									Fields: []*common_model.Field{
										{
											AsName:  "user",
											Columns: []*common_model.FieldColumn{{Table: "_root.or.user", Name: "user", Type: common_model.FieldReference}},
										},
										{
											AsName:  "description",
											Columns: []*common_model.FieldColumn{{Table: "_root.base", Name: "description", Type: common_model.FieldReference}},
										},
										{
											AsName:  "id",
											Columns: []*common_model.FieldColumn{{Table: "_root.base", Name: "id", Type: common_model.FieldReference}},
										},
										{
											AsName:  "status",
											Columns: []*common_model.FieldColumn{{Table: "_root.base", Name: "status", Type: common_model.FieldReference}},
										},
									},
								},
							},
						},
					},
					Tables: []*common_model.Table{
						{AsName: "_root.base", Name: "<select1>"},
						{AsName: "_root.or.user", Name: "<select2>", IsLateral: true},
					},
					SubScopes: []*common_model.StmtScope{
						{
							Name: "<select1>",
							Fields: []*common_model.Field{
								{Columns: []*common_model.FieldColumn{{Type: common_model.FieldStar}}},
								{Columns: []*common_model.FieldColumn{{ReferenceName: "<field1>", Type: common_model.FieldSubquery}}},
							},
							FieldScopes: []*common_model.StmtScope{
								{
									Name: "<field1>",
									Fields: []*common_model.Field{
										{Columns: []*common_model.FieldColumn{
											{Table: "__be_0_users", Name: "id", Type: common_model.FieldCondition},
											{Table: "tasks", Name: "user_id", Type: common_model.FieldCondition},
											{Table: "__be_0_users", Name: "email", Type: common_model.FieldCondition},
										}},
									},
									Tables: []*common_model.Table{{AsName: "__be_0_users", Name: "users"}},
								},
							},
							Tables: []*common_model.Table{{Name: "tasks"}},
						},
						{
							Name: "<select2>",
							Fields: []*common_model.Field{
								{
									AsName:  "user",
									Columns: []*common_model.FieldColumn{{ReferenceName: "<field2>", Type: common_model.FieldSubquery}},
								},
							},
							FieldScopes: []*common_model.StmtScope{
								{
									Name: "<field2>",
									Fields: []*common_model.Field{
										{Columns: []*common_model.FieldColumn{{Name: "_e", Type: common_model.FieldReference}}},
									},
									Tables: []*common_model.Table{{AsName: "_e", Name: "<select5>"}},
									SubScopes: []*common_model.StmtScope{
										{
											Name: "<select5>",
											Fields: []*common_model.Field{
												{
													AsName:  "email",
													Columns: []*common_model.FieldColumn{{Table: "_root.or.user.base", Name: "email", Type: common_model.FieldReference}},
												},
												{
													AsName:  "id",
													Columns: []*common_model.FieldColumn{{Table: "_root.or.user.base", Name: "id", Type: common_model.FieldReference}},
												},
												{
													AsName:  "name",
													Columns: []*common_model.FieldColumn{{Table: "_root.or.user.base", Name: "name", Type: common_model.FieldReference}},
												},
											},
										},
									},
								},
							},
							Tables: []*common_model.Table{{AsName: "_root.or.user.base", Name: "<select4>"}},
							SubScopes: []*common_model.StmtScope{
								{
									Name: "<select4>",
									Fields: []*common_model.Field{
										{Columns: []*common_model.FieldColumn{{Type: common_model.FieldStar}}},
										{Columns: []*common_model.FieldColumn{
											{Table: "_root.base", Name: "user_id", Type: common_model.FieldCondition},
											{Name: "id", Type: common_model.FieldCondition},
										}},
									},
									Tables: []*common_model.Table{{Name: "users"}},
								},
							},
						},
					},
				},
			},
		},
	},
}

var HasuraSubqueryWithAggregationData = &hasuraData{
	GQL: `query MyQuery($email: String) {
  users(where: {email: {_eq: $email}}) {
    email
    name
    tasks_aggregate {
      aggregate {
        avg {
          status
        }
      }
    }
    tasks(where: {}) {
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
                "_root.ar.root.tasks_aggregate"."tasks_aggregate" AS "tasks_aggregate",
                "_root.ar.root.tasks"."tasks" AS "tasks"
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
          ) AS "tasks_aggregate"
        FROM
          (
            SELECT
              "_root.ar.root.tasks_aggregate.base"."status" AS "status"
            FROM
              (
                SELECT
                  *
                FROM
                  "public"."tasks"
                WHERE
                  (("_root.base"."id") = ("user_id"))
              ) AS "_root.ar.root.tasks_aggregate.base"
          ) AS "_root.ar.root.tasks_aggregate"
      ) AS "_root.ar.root.tasks_aggregate" ON ('true')
      LEFT OUTER JOIN LATERAL (
        SELECT
          coalesce(json_agg("tasks"), '[]') AS "tasks"
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
                        "_root.ar.root.tasks.base"."title" AS "title",
                        "_root.ar.root.tasks.base"."description" AS "description"
                    ) AS "_e"
                )
              ) AS "tasks"
            FROM
              (
                SELECT
                  *
                FROM
                  "public"."tasks"
                WHERE
                  (("_root.base"."id") = ("user_id"))
              ) AS "_root.ar.root.tasks.base"
          ) AS "_root.ar.root.tasks"
      ) AS "_root.ar.root.tasks" ON ('true')
  ) AS "_root"`,
	Scopes: []*common_model.StmtScope{
		{
			Name: "<root>",
			Fields: []*common_model.Field{
				{
					AsName:  "root",
					Columns: []*common_model.FieldColumn{{Name: "root", Type: common_model.FieldReference}},
				},
			},
			Tables: []*common_model.Table{
				{AsName: "_root", Name: "<select0>"},
			},
			SubScopes: []*common_model.StmtScope{
				{
					Name: "<select0>",
					Fields: []*common_model.Field{
						{
							AsName:  "root",
							Columns: []*common_model.FieldColumn{{ReferenceName: "<field0>", Type: common_model.FieldSubquery}},
						},
					},
					FieldScopes: []*common_model.StmtScope{
						{
							Name: "<field0>",
							Fields: []*common_model.Field{
								{Columns: []*common_model.FieldColumn{{Name: "_e", Type: common_model.FieldReference}}},
							},
							Tables: []*common_model.Table{{AsName: "_e", Name: "<select4>"}},
							SubScopes: []*common_model.StmtScope{
								{
									Name: "<select4>",
									Fields: []*common_model.Field{
										{
											AsName:  "email",
											Columns: []*common_model.FieldColumn{{Table: "_root.base", Name: "email", Type: common_model.FieldReference}},
										},
										{
											AsName:  "name",
											Columns: []*common_model.FieldColumn{{Table: "_root.base", Name: "name", Type: common_model.FieldReference}},
										},
										{
											AsName:  "tasks_aggregate",
											Columns: []*common_model.FieldColumn{{Table: "_root.ar.root.tasks_aggregate", Name: "tasks_aggregate", Type: common_model.FieldReference}},
										},
										{
											AsName:  "tasks",
											Columns: []*common_model.FieldColumn{{Table: "_root.ar.root.tasks", Name: "tasks", Type: common_model.FieldReference}},
										},
									},
								},
							},
						},
					},
					Tables: []*common_model.Table{
						{AsName: "_root.base", Name: "<select1>"},
						{AsName: "_root.ar.root.tasks_aggregate", Name: "<select2>", IsLateral: true},
						{AsName: "_root.ar.root.tasks", Name: "<select3>", IsLateral: true},
					},
					SubScopes: []*common_model.StmtScope{
						{
							Name: "<select1>",
							Fields: []*common_model.Field{
								{Columns: []*common_model.FieldColumn{{Type: common_model.FieldStar}}},
								{Columns: []*common_model.FieldColumn{{Table: "users", Name: "email", Type: common_model.FieldCondition}}},
							},
							Tables: []*common_model.Table{{Name: "users"}},
						},
						{
							Name: "<select2>",
							Fields: []*common_model.Field{
								{AsName: "tasks_aggregate", Columns: []*common_model.FieldColumn{{Name: "status", Type: common_model.FieldReference}}},
							},
							Tables: []*common_model.Table{{AsName: "_root.ar.root.tasks_aggregate", Name: "<select5>"}},
							SubScopes: []*common_model.StmtScope{
								{
									Name: "<select5>",
									Fields: []*common_model.Field{
										{AsName: "status", Columns: []*common_model.FieldColumn{
											{Table: "_root.ar.root.tasks_aggregate.base", Name: "status", Type: common_model.FieldReference},
										}},
									},
									Tables: []*common_model.Table{{AsName: "_root.ar.root.tasks_aggregate.base", Name: "<select6>"}},
									SubScopes: []*common_model.StmtScope{
										{
											Name: "<select6>",
											Fields: []*common_model.Field{
												{Columns: []*common_model.FieldColumn{{Type: common_model.FieldStar}}},
												{Columns: []*common_model.FieldColumn{
													{Table: "_root.base", Name: "id", Type: common_model.FieldCondition},
													{Name: "user_id", Type: common_model.FieldCondition},
												}},
											},
											Tables: []*common_model.Table{{Name: "tasks"}},
										},
									},
								},
							},
						},
						{
							Name: "<select3>",
							Fields: []*common_model.Field{
								{
									AsName:  "tasks",
									Columns: []*common_model.FieldColumn{{Name: "tasks", Type: common_model.FieldReference}},
								},
							},
							Tables: []*common_model.Table{{AsName: "_root.ar.root.tasks", Name: "<select7>"}},
							SubScopes: []*common_model.StmtScope{
								{
									Name: "<select7>",
									Fields: []*common_model.Field{
										{AsName: "tasks", Columns: []*common_model.FieldColumn{{ReferenceName: "<field1>", Type: common_model.FieldSubquery}}},
									},
									FieldScopes: []*common_model.StmtScope{
										{
											Name: "<field1>",
											Fields: []*common_model.Field{
												{Columns: []*common_model.FieldColumn{{Name: "_e", Type: common_model.FieldReference}}},
											},
											Tables: []*common_model.Table{{AsName: "_e", Name: "<select9>"}},
											SubScopes: []*common_model.StmtScope{
												{
													Name: "<select9>",
													Fields: []*common_model.Field{
														{
															AsName:  "title",
															Columns: []*common_model.FieldColumn{{Table: "_root.ar.root.tasks.base", Name: "title", Type: common_model.FieldReference}},
														},
														{
															AsName:  "description",
															Columns: []*common_model.FieldColumn{{Table: "_root.ar.root.tasks.base", Name: "description", Type: common_model.FieldReference}},
														},
													},
												},
											},
										},
									},
									Tables: []*common_model.Table{{AsName: "_root.ar.root.tasks.base", Name: "<select8>"}},
									SubScopes: []*common_model.StmtScope{
										{
											Name: "<select8>",
											Fields: []*common_model.Field{
												{Columns: []*common_model.FieldColumn{{Type: common_model.FieldStar}}},
												{Columns: []*common_model.FieldColumn{
													{Table: "_root.base", Name: "id", Type: common_model.FieldCondition},
													{Name: "user_id", Type: common_model.FieldCondition},
												}},
											},
											Tables: []*common_model.Table{{Name: "tasks"}},
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
