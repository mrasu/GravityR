import analyzeNodes from "./analyzeNodes.json";

const query = `
with count_tbl as
( select
  status,
  count(status) as count
  from users
  inner join todos on users.user_id = todos.user_id
  where users.created_at > now() - interval 2  month
  group by status
)
select
  case when status = 1 then 'Todo'
       when status = 2 then 'Doing'
       when status = 3 then 'Done' end as status,
  sum(count)
from  count_tbl;
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

export const dummySuggestData = {
  query: query,
  analyzeNodes: analyzeNodes,
  indexTargets: indexTargets,
  examinationResult: examinationResult,
  examinationCommandOptions: [
    { isShort: true, name: "o", value: "output_examine.html" },
    { isShort: true, name: "q", value: query },
  ],
};
