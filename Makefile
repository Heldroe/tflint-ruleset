BINARY=tflint-ruleset-terraform_style

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

install:
	mkdir -p ~/.tflint.d/plugins
	cp $(BINARY) ~/.tflint.d/plugins

dev: fmt lint test build
