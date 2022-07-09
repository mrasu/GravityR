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
