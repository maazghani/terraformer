.PHONY: test test-unit test-integration fmt vet check build

test:
	go test ./...

test-unit:
	go test ./internal/...

test-integration:
	TERRAFORMER_RUN_INTEGRATION=1 go test ./... -run Integration

fmt:
	gofmt -w ./cmd ./internal

vet:
	go vet ./...

check:
	go test ./...
	go vet ./...

build:
	go build ./cmd/terraformer-mcp
