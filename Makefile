.PHONY: sample

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

define sample_query
SELECT
	name,
	t.description
FROM
	users
	INNER JOIN todos AS t ON users.id = t.user_id
WHERE
	users.email = 'test31776@example.com'
endef
export sample_query

sample:
	./dist/gr db suggest -o "sample/output.html" -q "$${sample_query}"
	./dist/gr db suggest --with-examine -o "sample/output_examine.html" -q "$${sample_query}"

