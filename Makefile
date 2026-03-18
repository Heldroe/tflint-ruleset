BINARY=tflint-ruleset-terraform_style

.PHONY: build clean fmt lint test install dev check

build:
	go build -o $(BINARY)

clean:
	rm -f $(BINARY)

fmt:
	go fmt ./...

lint:
	go vet ./...

test:
	go test -v ./...

install:
	mkdir -p ~/.tflint.d/plugins
	cp $(BINARY) ~/.tflint.d/plugins

check: fmt lint test

dev: check build
