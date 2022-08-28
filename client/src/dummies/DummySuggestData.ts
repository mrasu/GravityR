import analyzeMysqlNodes from "./analyzeMysqlNodes.json";
import analyzePostgresNodes from "./analyzePostgresNodes.json";

const query = `
with count_tbl as
( select
  status,
  count(status) as count
  from users
  inner join todos on users.id = todos.user_id
  where users.created_at > now() - interval 2  month
  group by status
)
select
  case when status = 1 then 'Todo'
       when status = 2 then 'Doing'
       when status = 3 then 'Done' end as status,
  sum(count)
from  count_tbl
group by status;
`;
const indexTargets = [
  { tableName: "todos", columns: [{ name: "user_id" }] },
  {
    tableName: "todos",
    columns: [{ name: "user_id" }, { name: "description" }],
  },
  { tableName: "users", columns: [{ name: "id" }, { name: "name" }] },
];

const examinationResult = {
  originalTimeMillis: 1559,
  indexResults: [
    {
      indexTarget: { tableName: "todos", columns: [{ name: "user_id" }] },
      executionTimeMillis: 12,
    },
    {
      indexTarget: {
        tableName: "todos",
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

const postgresPlanningText = `Planning:
  Buffers: shared hit=6
Planning Time: 0.167 ms
JIT:
  Functions: 57
  Options: Inlining false, Optimization false, Expressions true, Deforming true
  Timing: Generation 3.100 ms, Inlining 0.000 ms, Optimization 1.174 ms, Emission 20.082 ms, Total 24.357 ms
Execution Time: 375.843 ms`;

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
    planningText: postgresPlanningText,
    indexTargets: indexTargets,
    examinationResult: examinationResult,
    examinationCommandOptions: [
      { isShort: true, name: "o", value: "output_examine.html" },
      { isShort: true, name: "q", value: query },
    ],
  },
};
