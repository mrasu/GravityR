.PHONY: example

export DB_USERNAME := root
export DB_DATABASE := gravityr
export HASURA_URL := http://localhost:8081
export HASURA_ADMIN_SECRET := myadminsecretkey

devSetup:
	go install github.com/spf13/cobra-cli@latest

lint:
	golangci-lint run --fix

test:
	go test ./...

build:
	cd client && yarn build
	go build -o dist/gr main.go

clean:
	rm -rf dist/
	rm -rf client/dist/

define example_query
SELECT
	name,
	t.description
FROM
	users
	INNER JOIN todos AS t ON users.id = t.user_id
WHERE
	users.email = 'test31776@example.com'
endef
export example_query

define example_gql
query MyQuery($$email: String) {
  users(where: {email: {_eq: $$email}}) {
    email
    name
    todos(where: {}) {
      title
      description
    }
  }
}
endef
export example_gql

define example_gql_variables
{"email": "test33333@example.com"}
endef
export example_gql_variables


example:
	mkdir -p example
	./dist/gr db suggest hasura -o "example/hasura.html" -q "$${example_gql}" --json-variables "$${example_gql_variables}"
	./dist/gr db suggest hasura --with-examine -o "example/hasura_examine.html" -q "$${example_gql}" --json-variables "$${example_gql_variables}"
	./dist/gr db suggest mysql -o "example/mysql.html" -q "$${example_query}"
	./dist/gr db suggest mysql --with-examine -o "example/mysql_examine.html" -q "$${example_query}"
	./dist/gr db suggest postgres -o "example/postgres.html" -q "$${example_query}"
	./dist/gr db suggest postgres --with-examine -o "example/postgres_examine.html" -q "$${example_query}"
	./dist/gr db dig performance-insights -o "example/performance-insights.html"  --use-mock --start-from 2022-08-04T14:00:00Z

