.PHONY: example

export DB_USERNAME := root
export DB_DATABASE :=gravityr

devSetup:
	go install github.com/spf13/cobra-cli@latest

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

example:
	mkdir -p example
	./dist/gr db suggest -o "example/output.html" -q "$${example_query}"
	./dist/gr db suggest --with-examine -o "example/output_examine.html" -q "$${example_query}"
	./dist/gr db dig performance-insights -o "example/performance-insights.html"  --use-mock --start-from 2022-08-04T14:00:00Z

