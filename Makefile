.PHONY: test test-unit test-integration fmt vet check build

test:
	go test -count=10 ./...

test-unit:
	go test -count=10 ./internal/...

test-integration:
	TERRAFORMER_RUN_INTEGRATION=1 go test -count=10 ./... -run Integration

fmt:
	gofmt -w ./cmd ./internal

vet:
	go vet ./...

check:
	go test -count=10 ./...
	go vet ./...

build:
	go build ./cmd/terraformer-mcp
