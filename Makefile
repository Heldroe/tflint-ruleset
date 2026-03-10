BINARY=tflint-ruleset

build:
	go build -o $(BINARY)

clean:
	rm -f $(BINARY)

fmt:
	go fmt ./...

lint:
	go vet ./...

test:
	go test -v ./rules/... ./tests/...

dev: fmt lint test build
