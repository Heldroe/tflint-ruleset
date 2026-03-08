BINARY=tflint-ruleset

build:
	go build -o $(BINARY)

install: build
	mkdir -p ~/.tflint.d/plugins
	cp $(BINARY) ~/.tflint.d/plugins/

clean:
	rm -f $(BINARY)

fmt:
	go fmt ./...

lint:
	go vet ./...

dev: fmt lint build
