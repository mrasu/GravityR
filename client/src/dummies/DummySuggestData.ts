import analyzeMysqlNodes from "./analyzeMysqlNodes.json";
import analyzePostgresNodes from "./analyzePostgresNodes.json";
import analyzeHasuraPostgresNodes from "./analyzeHasuraPostgresNodes.json";
import hasuraPostgresQuery from "./hasuraPostgresQuery.sql?raw";

const query = `
with count_tbl as
( select
  status,
  count(status) as count
  from users
  inner join tasks on users.id = tasks.user_id
  where users.created_at > now() - interval 2  month
  group by status
)
select
  case when status = 1 then 'Pending'
       when status = 2 then 'Doing'
       when status = 3 then 'Done' end as status,
  sum(count)
from  count_tbl
group by status;
`;
const indexTargets = [
  { tableName: "tasks", columns: [{ name: "user_id" }] },
  {
    tableName: "tasks",
    columns: [{ name: "user_id" }, { name: "description" }],
  },
  { tableName: "users", columns: [{ name: "id" }, { name: "name" }] },
];

const examinationResult = {
  originalTimeMillis: 1559,
  indexResults: [
    {
      indexTarget: { tableName: "tasks", columns: [{ name: "user_id" }] },
      executionTimeMillis: 12,
    },
    {
      indexTarget: {
        tableName: "tasks",
        columns: [{ name: "user_id" }, { name: "description" }],
      },
      executionTimeMillis: 10,
    },
    {
      indexTarget: {
        tableName: "users",
        columns: [{ name: "id" }, { name: "name" }],
      },
      executionTimeMillis: 1593,
    },
  ],
};

const postgresSummaryText = `Planning:
  Buffers: shared hit=6
Planning Time: 0.167 ms
JIT:
  Functions: 57
  Options: Inlining false, Optimization false, Expressions true, Deforming true
  Timing: Generation 3.100 ms, Inlining 0.000 ms, Optimization 1.174 ms, Emission 20.082 ms, Total 24.357 ms
Execution Time: 375.843 ms`;

const gql = `
query MyQuery($email: String) {
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
`;

const hasuraIndexTargets = [
  { tableName: "tasks", columns: [{ name: "user_id" }] },
  {
    tableName: "tasks",
    columns: [{ name: "user_id" }, { name: "status" }, { name: "title" }],
  },
  { tableName: "users", columns: [{ name: "email" }] },
  { tableName: "users", columns: [{ name: "email" }, { name: "id" }] },
  { tableName: "users", columns: [{ name: "id" }, { name: "email" }] },
  {
    tableName: "users",
    columns: [{ name: "email" }, { name: "id" }, { name: "name" }],
  },
  {
    tableName: "users",
    columns: [{ name: "id" }, { name: "email" }, { name: "name" }],
  },
];

const hasuraPostgresSummaryText = `JIT:
  Functions: 30
  Options: Inlining false, Optimization false, Expressions true, Deforming true`;

export const dummySuggestData = {
  mysql: {
    query: query,
    analyzeNodes: analyzeMysqlNodes,
    indexTargets: indexTargets,
    examinationResult: examinationResult,
    examinationCommandOptions: [
      { isShort: true, name: "o", value: "output_examine.html" },
      { isShort: true, name: "q", value: query },
    ],
  },
  postgres: {
    query: query,
    analyzeNodes: analyzePostgresNodes,
    summaryText: postgresSummaryText,
    indexTargets: indexTargets,
    examinationResult: examinationResult,
    examinationCommandOptions: [
      { isShort: true, name: "o", value: "output_examine.html" },
      { isShort: true, name: "q", value: query },
    ],
  },
  hasura: {
    postgres: {
      gql: gql,
      gqlVariables: { email: "test1111@example.com" },
      query: hasuraPostgresQuery,
      analyzeNodes: analyzeHasuraPostgresNodes,
      summaryText: hasuraPostgresSummaryText,
      indexTargets: hasuraIndexTargets,
      examinationResult: examinationResult,
      examinationCommandOptions: [
        { isShort: true, name: "o", value: "output_examine.html" },
        { isShort: true, name: "q", value: gql },
        {
          isShort: false,
          name: "json-variables",
          value: '{"email": "test1112@example.com"}',
        },
      ],
    },
  },
};
