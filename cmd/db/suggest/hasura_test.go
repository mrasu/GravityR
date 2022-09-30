package suggest

import (
	"github.com/jarcoal/httpmock"
	"github.com/mrasu/GravityR/infra/hasura"
	"github.com/mrasu/GravityR/thelper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/url"
	"os"
	"testing"
)

func Test_runHasuraSuggest(t *testing.T) {
	hr := hasuraRunner{
		query: `
query MyQuery($$email: String) {
  users(where: {email: {_eq: $$email}}) {
    email
    name
    todos(where: {}) {
      title
      description
    }
  }
}`,
		jsonVariables:   `{"email": "test33333@example.com"}`,
		runsExamination: true,
		indexTargets:    []string{"users:name"},
	}
	u, err := url.Parse("https://example.com")
	require.NoError(t, err)
	cli := hasura.NewClient(&hasura.Config{Url: u})

	thelper.InjectClientDist()
	m := hasuraMock{}
	m.mockAll(func() {
		thelper.CreateTemp(t, "tmp.html", func(tmpfile *os.File) {
			err := hr.suggest(tmpfile.Name(), cli)
			require.NoError(t, err)

			html, err := os.ReadFile(tmpfile.Name())
			require.NoError(t, err)
			assert.Contains(t, string(html), "Nested Loop Left Join  (cost=149485.77..151050.00")
			assert.Contains(t, string(html), "MyQuery($$email: String) {")
			assert.Contains(t, string(html), `originalTimeMillis":0,"indexResults":[{"`)
		})
	})
}

type hasuraMock struct{}

func (m *hasuraMock) mockAll(fn func()) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	m.mockExplain()
	m.mockTableSchemaQuery()
	m.mockGqlQuery()

	fn()
}

func (m *hasuraMock) mockExplain() {
	httpmock.RegisterResponder("POST", "https://example.com/v1/graphql/explain",
		httpmock.NewStringResponder(200, `
[
	{
		"field": "todos",
		"plan": [
			"ggregate  (cost=151050.01..151050.02 rows=1 width=32)",
			"  ->  Nested Loop Left Join  (cost=149485.77..151050.00 rows=1 width=62)",
			"        ->  Seq Scan on users  (cost=0.00..1564.20 rows=1 width=34)",
			"              Filter: ((email)::text = 'test33333@example.com'::text)",
			"        ->  Aggregate  (cost=149485.77..149485.78 rows=1 width=32)",
			"              ->  Seq Scan on todos  (cost=0.00..149483.14 rows=175 width=28)",
			"                    Filter: (users.id = user_id)",
			"              SubPlan 2",
			"                ->  Result  (cost=0.00..0.01 rows=1 width=32)",
			"  SubPlan 1",
			"    ->  Result  (cost=0.00..0.01 rows=1 width=32)"
		],
		"sql": "SELECT coalesce(json_agg(\"root\"), '[]') AS \"root\"FROM ( SELECT row_to_json( ( SELECT \"_e\" FROM ( SELECT \"_root.base\".\"email\" AS \"email\", \"_root.base\".\"name\" AS \"name\", \"_root.ar.root.todos\".\"todos\" AS \"todos\" ) AS \"_e\" ) ) AS \"root\" FROM ( SELECT * FROM \"public\".\"users\" WHERE ( (\"public\".\"users\".\"email\") = (('test33333@example.com')::varchar) ) ) AS \"_root.base\" LEFT OUTER JOIN LATERAL ( SELECT coalesce(json_agg(\"todos\"), '[]') AS \"todos\" FROM ( SELECT row_to_json( ( SELECT \"_e\" FROM ( SELECT \"_root.ar.root.todos.base\".\"title\" AS \"title\", \"_root.ar.root.todos.base\".\"description\" AS \"description\" ) AS \"_e\" ) ) AS \"todos\" FROM ( SELECT * FROM \"public\".\"todos\" WHERE ((\"_root.base\".\"id\") = (\"user_id\")) ) AS \"_root.ar.root.todos.base\" ) AS \"_root.ar.root.todos\" ) AS \"_root.ar.root.todos\" ON ('true') ) AS \"_root\""
	}
]`),
	)
}

func (m *hasuraMock) mockTableSchemaQuery() {
	httpmock.RegisterResponder("POST", "https://example.com/v2/query",
		httpmock.NewStringResponder(200, `
{
	"result_type": "TuplesOk",
	"result": [
		[
			"table_name",
			"column_name",
			"is_pk"
		],
		[
			"todos",
			"user_id",
			"f"
		],
		[
			"todos",
			"title",
			"f"
		],
		[
			"todos",
			"description",
			"f"
		],
		[
			"users",
			"id",
			"t"
		],
		[
			"users",
			"name",
			"f"
		],
		[
			"users",
			"email",
			"f"
		]
	]
}
`),
	)
}

func (m *hasuraMock) mockGqlQuery() {
	httpmock.RegisterResponder("POST", "https://example.com/v1/graphql",
		httpmock.NewStringResponder(200, `
{
	"data": {
		"todos": [
			{
				"user": {
					"email": "test33333@example.com",
					"id": 33333,
					"name": "test33333"
				},
				"description": "test description",
				"id": 111201,
				"status": 1
			}
		]
	}
}`),
	)
}
