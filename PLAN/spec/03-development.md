# Specification: Local Development Commands

This document specifies the Makefile targets and development workflow commands for building, testing, and running the terraformer server locally.

## Local development commands
These commands should be added to the Makefile as the repo is built.

### make test
Expected behavior:
go test ./...

### make test-unit
Expected behavior:
go test ./internal/...

### make test-integration
Expected behavior:
TERRAFORMER_RUN_INTEGRATION=1 go test ./... -run Integration

### make fmt
Expected behavior:
gofmt -w ./cmd ./internal

### make vet
Expected behavior:
go vet ./...

### make check
Expected behavior:
go test ./...
go vet ./...

### make build
Expected behavior:
go build ./cmd/terraformer-mcp

### make run
Expected behavior (with environment variable):
./terraformer-mcp --repo-root=/path/to/repo --port=9001

*** No task is complete until the relevant targeted tests pass. Before commit, make check should pass unless the plan explicitly records why it cannot yet pass ***

## See also

- [00-spec.md](00-spec.md) — Safety rules and structural constraints
- [01-testing.md](01-testing.md) — Testing requirements and test categories
- [02-mcp-tool-contracts.md](02-mcp-tool-contracts.md) — Tool contracts and expected response shapes