

.PHONY: generate
generate:
	go generate ./...

.PHONY: test
test: generate
	go test -v -cover -race ./...