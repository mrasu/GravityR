import { PostgresIndexTarget } from "@/models/PostgresIndexTarget";

export class HasuraIndexTarget extends PostgresIndexTarget {
  toExecutableText(): string {
    return `curl --request POST \
  --url "\${HASURA_URL}/v2/query" \
  --header 'Content-Type: application/json' \
  --header "x-hasura-admin-secret: \${HASURA_ADMIN_SECRET}" \
  --data '{
	"type": "run_sql",
	"args": {
		"sql": "${this.toAlterAddSQL()}"
	}
}'`;
  }
}
