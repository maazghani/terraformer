# terraformer

**terraformer v0** is a local-only HTTP/JSON server that exposes safe, read-only Terraform operations as MCP (Model Context Protocol) tools for AI coding assistants. It enables LLM-powered agents to inspect Terraform configurations, validate HCL syntax, generate plans, and compare desired state—without ever applying changes or destroying infrastructure.

## Non-Goals

terraformer v0 **does not**:

- Apply Terraform changes (`terraform apply` is forbidden)
- Destroy infrastructure (`terraform destroy` is forbidden)
- Execute arbitrary shell commands
- Mutate remote backends
- Operate on remote repositories (local-only)

This design ensures agents can safely explore and analyze Terraform configurations without risk of unintended infrastructure modifications.

## Build Instructions

```bash
make build
```

This produces a `./terraformer-mcp` binary in the repository root.

## Quickstart

1. **Start the server** against a local Terraform repository:

```bash
./terraformer-mcp --repo-root=/absolute/path/to/your/terraform/repo
```

The server listens on `http://localhost:9001` by default (configurable via `--port`).

2. **Test with curl**:

```bash
# List files in the repo
curl -X POST http://localhost:9001/tools/list_repo_files \
  -H "Content-Type: application/json" \
  -d '{"path": ".", "include_globs": ["*.tf"], "exclude_globs": [".terraform/**"], "max_files": 100}'

# Read a file
curl -X POST http://localhost:9001/tools/read_repo_file \
  -H "Content-Type: application/json" \
  -d '{"path": "main.tf", "max_bytes": 65536}'

# Validate Terraform configuration
curl -X POST http://localhost:9001/tools/terraform_validate \
  -H "Content-Type: application/json" \
  -d '{"json": true}'
```

## Client Integration

For detailed integration instructions with specific AI coding assistants (GitHub Copilot, Cursor, Claude Code, Codex, Cline), see [USAGE.md](USAGE.md).

## Tool Contracts

All tool request/response schemas and detailed behavior are documented in [PLAN/spec/02-mcp-tool-contracts.md](PLAN/spec/02-mcp-tool-contracts.md).

Available tools:
- `terraform_init` — Initialize Terraform working directory
- `terraform_fmt` — Format Terraform files
- `terraform_validate` — Validate configuration syntax
- `terraform_plan` — Generate execution plan
- `terraform_show_json` — Parse plan file to JSON
- `list_repo_files` — List files in repository
- `read_repo_file` — Read file contents
- `apply_patch` — Write structured file changes
- `check_desired_state` — Compare plan against desired state (minimal in v0)

## Testing

```bash
# Run all tests (unit + integration if enabled)
make test

# Run only unit tests (no Terraform required)
make test-unit

# Run integration tests (requires terraform on PATH)
make test-integration

# Run tests + vet
make check
```

### Integration Tests

Integration tests require Terraform to be installed and available on your `PATH`. They are **opt-in** via an environment variable:

```bash
TERRAFORMER_RUN_INTEGRATION=1 make test-integration
```

Without this variable, integration tests are skipped cleanly.

## Command-Line Flags

```bash
./terraformer-mcp \
  --repo-root=/path/to/repo \      # Required: absolute path to Terraform repo
  --port=9001 \                     # Optional: TCP port (default 9001)
  --log-level=info \                # Optional: debug|info|warn|error (default info)
  --max-response-bytes=1048576 \   # Optional: max response size in bytes (default 1 MiB)
  --terraform-bin=terraform         # Optional: path to terraform binary (default "terraform")
```

Use `--help` to see all flags and their defaults.

## Example Workflow

A typical agent workflow for proposing Terraform changes:

1. **`list_repo_files`** — Discover `.tf` files in the repository
2. **`read_repo_file`** — Read specific configuration files
3. **`apply_patch`** — Propose changes to configuration files
4. **`terraform_fmt`** — Format the modified files
5. **`terraform_validate`** — Validate syntax and references
6. **`terraform_plan`** — Generate an execution plan
7. **`check_desired_state`** — Compare plan against desired state

The agent can iterate on changes based on validation errors or plan output, but **never applies changes** to actual infrastructure.

## Safety Model

terraformer enforces strict safety boundaries:

- **Path Safety**: All operations are scoped to the configured `--repo-root`. Path traversal attempts (`../`, absolute paths outside repo, symlink escapes) are rejected. See [internal/safety/safety.go](internal/safety/safety.go) for implementation details.

- **Command Allowlist**: Only safe, read-only Terraform commands are permitted:
  - `terraform init`
  - `terraform fmt`
  - `terraform validate`
  - `terraform plan`
  - `terraform show -json`

  Dangerous commands (`apply`, `destroy`, `import`, `workspace`, `state`, etc.) are **forbidden** and cannot be invoked.

- **No Shell Execution**: Commands are executed directly via Go's `os/exec`, never through a shell. No shell interpolation or metacharacter expansion occurs.

- **Local-Only**: The server operates only on local filesystem paths. It makes no network calls except those made by Terraform providers themselves (e.g., when validating configurations that reference remote data sources).

## Known Limitations (v0)

The following limitations are **intentional** for v0 and deferred to future releases:

- **Local-only**: No remote repository access or remote backend mutation
- **No apply/destroy**: Cannot modify actual infrastructure
- **No arbitrary shell commands**: Only allowlisted Terraform commands
- **Minimal desired-state checking**: `check_desired_state` is intentionally basic in v0
- **No output redaction**: Secrets in Terraform output may be visible (v0.1 will add redaction)
- **No concurrency management**: Concurrent requests share the repo root (v0.1 will address)
- **No secret-file exclusion**: `list_repo_files` and `read_repo_file` don't filter files like `.env` or `credentials.json` (v0.1 will add filtering)

## Demo

See [scripts/demo.sh](scripts/demo.sh) for a complete end-to-end demonstration that builds the server, starts it against a fixture repository, and exercises all tools via `curl`.

## License

See [LICENSE](LICENSE) for details.
