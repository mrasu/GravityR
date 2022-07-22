.PHONY: sample

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

sample:
	./dist/gr db suggest -o "sample/output.html" -q "SELECT name, t.description FROM users INNER JOIN todos AS t ON users.id = t.user_id WHERE users.name = 'foo'"
	./dist/gr db suggest --with-examine -o "sample/output_examine.html" -q "SELECT name, t.description FROM users INNER JOIN todos AS t ON users.id = t.user_id WHERE users.name = 'foo'"
