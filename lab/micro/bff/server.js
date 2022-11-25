require("./tel.js");
const express = require("express");
const { graphqlHTTP } = require("express-graphql");
const { buildSchema } = require("graphql");
const opentelemetry = require("@opentelemetry/api");

const schema = buildSchema(`
type Query {
  users(ids: [Int!]!): [User!]!
}

type User {
  id: Int!
  name: String!
  email: String!
  histories: [PaymentHistory!]
}

type PaymentHistory {
  amount: Int!
  last_numbers: String!
}
`);

const getTraceId = () => {
  const span = opentelemetry.trace.getActiveSpan();
  const ctx = span.spanContext();

  return `00-${ctx.traceId}-${ctx.spanId}-00`;
};

const handleUsers = async (args) => {
  const ids = args.ids;
  const traceId = getTraceId();
  const promises = ids.map((id) =>
    (async function (id) {
      return fetch(`http://user:9002/users/${id}`, {
        headers: {
          // Couldn't find a way to inject traceparent header automatically
          traceparent: traceId,
        },
      })
        .then((r) => r.json())
        .then((v) => ({
          id: v.id,
          name: v.name,
          email: v.email,
          histories: v.histories.map((h) => ({
            amount: h.amount,
            last_numbers: h.last_numbers,
          })),
        }));
    })(id)
  );
  return await Promise.all(promises);
};

const root = {
  users: handleUsers,
};

const app = express();
app.use(
  "/graphql",
  graphqlHTTP({
    schema: schema,
    rootValue: root,
  })
);
app.listen(9001);
console.log("Running a BFF service at http://localhost:9001/graphql");
