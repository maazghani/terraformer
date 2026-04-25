# Specification: terraformer v0 Requirements

This document defines the non-negotiable safety rules, repo structure, package responsibilities, and implementation boundary for terraformer v0. All architectural and implementation decisions must comply with these specifications.

## Non-negotiable safety rules
These rules are project invariants. Tests must cover them.
 The server must only operate inside the configured repository root.
 The server must reject path traversal attempts such as ../, absolute paths outside the repo, symlink escapes, and malformed paths.
 The server must never execute commands outside the configured repo root.
 The server must never execute arbitrary shell commands.
 The only Terraform commands allowed in v0 are:
- terraform init
- terraform fmt
- terraform validate
- terraform plan
- terraform show -json
 The following commands are forbidden in v0:
- terraform apply
- terraform destroy
- arbitrary shell execution
- any Terraform command not explicitly allowlisted
 Command execution must go through a testable command runner abstraction.
 Unit tests must be able to test Terraform command behavior without executing Terraform.
 Redaction and secret filtering deferred to v0.1.
 Every MCP tool must use structured request and response objects.
 Terraform command responses must include:
- stdout
- stderr
- exit code
- duration
- command metadata
- normalized diagnostics where possible

 terraform plan success must be treated as necessary but not sufficient.
 Desired-state comparison must be modeled as a separate check, initially minimal or stubbed, and expanded behind tests.
 Small, boring, auditable Go code is preferred over clever abstractions.
 Every meaningful capability must be developed test-first.
 A task may only be marked complete when tests exist, tests pass, and the code is committed or ready to commit.

## Intended repo structure
terraformer/
├── PLAN.md
├── README.md
├── go.mod
├── go.sum
├── Makefile
├── cmd/
│   └── terraformer-mcp/
│       └── main.go
├── internal/
│   ├── httpserver/
│   ├── tools/
│   ├── terraform/
│   ├── runner/
│   ├── repo/
│   ├── patch/
│   ├── desiredstate/
│   └── diagnostics/
├── testdata/
│   ├── fixtures/
│   │   ├── valid-basic/
│   │   ├── invalid-hcl/
│   │   ├── missing-provider/
│   │   └── plan-basic/
│   └── golden/
└── scripts/

This structure may be adjusted only when tests or implementation reveal a concrete reason. Record the reason in this plan before changing the layout.

## Suggested Go package layout
### cmd/terraformer-mcp
CLI entrypoint for the MCP server.

Responsibilities:
 Parse configuration.
 Resolve and validate repo root.
 Construct dependencies.
 Start the MCP server.
 Avoid business logic.

### internal/httpserver
HTTP/JSON-RPC server wiring.
Responsibilities:
 Start HTTP server on configured port (default 9001).
 Register tools.
 Convert HTTP requests into internal typed requests.
 Convert internal typed responses into HTTP/JSON responses.
 Keep HTTP server and protocol details isolated.
 Emit structured JSON logs to stdout/stderr.
 Avoid Terraform, filesystem, and patch logic.

### internal/tools
Tool orchestration layer.
Responsibilities:
 Define tool request and response structs.
 Call repo, Terraform, patch, desired-state, and diagnostics services.
 Enforce structured outputs.
 Avoid raw shell execution.

### internal/terraform
Terraform-specific command service.
Responsibilities:
 Build allowed Terraform command invocations.
 Reject forbidden commands.
 Parse Terraform outputs where appropriate.
 Normalize Terraform command results.
 Avoid direct filesystem traversal except through safe repo abstractions.

### internal/runner
Command runner abstraction.
Responsibilities:
 Define a Runner interface.
 Provide a real local process runner.
 Provide fake/test runners.
 Capture stdout, stderr, exit code, duration, and errors.
 Execute commands with a fixed working directory.
 Avoid shell interpolation.

### internal/repo
Repo-local filesystem access.
Responsibilities:
 Resolve repo-relative paths safely.
 List files.
 Read files.
 Write files only through approved workflows.
 Reject path traversal and symlink escapes.
 Normalize file metadata.

### internal/patch
Patch application.
Responsibilities:
 Accept structured patch requests.
 Validate target paths through repo.
 Apply patches atomically where practical.
 Return structured patch results.
 Avoid touching files outside the repo.

### internal/desiredstate
Desired-state model and comparison logic.
Responsibilities:
 Define desired-state input schema.
 Define comparison result schema.
 Initially implement minimal or stubbed behavior behind tests.
 Later compare Terraform plan JSON against desired-state expectations.
 Never declare success based on terraform plan alone.

### internal/diagnostics
Diagnostic normalization.
Responsibilities:
 Normalize Terraform validation and plan errors.
 Parse JSON diagnostics where available.
 Extract file, range, severity, summary, and detail where possible.
 Redaction of diagnostics deferred to v0.1.

### internal/safety
Safety policies and reusable guards.
Responsibilities:
 Path safety helpers.
 Command allowlist helpers.
 Secret redaction helpers.
 Environment filtering helpers.
 Safety-specific test fixtures.

## See also

- [01-testing.md](01-testing.md) — Testing philosophy, test categories, runner tests, golden files, and integration test requirements
- [02-mcp-tool-contracts.md](02-mcp-tool-contracts.md) — Tool request/response shapes and per-tool TDD requirements
- [03-development.md](03-development.md) — Makefile targets and local development commands

## Safety model for v0
The safety model focuses on path safety and command allowlisting. Redaction deferred to v0.1.
### Layer 1: Configuration safety
 Resolve repo root at startup.
 Require repo root to exist.
 Require repo root to be a directory.
 Convert repo root to an absolute cleaned path.
 Avoid accepting repo root from an untrusted HTTP request after startup.
 Treat repo root as immutable for the lifetime of the server.
### Layer 2: Path safety
 All user-supplied paths must be repo-relative unless explicitly documented otherwise.
 Absolute user paths are rejected by default.
 Clean paths before use.
 Resolve symlinks when needed.
 Confirm the final resolved path remains inside repo root.
 Reject paths into .git by default.
 Exclude .terraform from file listing by default.
 Secret-looking file exclusion deferred to v0.1.
### Layer 3: Command safety
 No command may be executed except through internal/runner.
 Terraform commands must be built from typed options, not raw strings.
 Never invoke through a shell.
 Never concatenate user input into command strings.
 Use argument arrays.
 Enforce the Terraform command allowlist.
 Run commands only with working directory set to repo root.
 Redact secrets from stdout.
 Redact secrets from stderr.
 Redact secrets from diagnostics.
 Redact secrets from plan summaries.
 Avoid returning provider credentials.
 Avoid returning full environment variables.
 Apply response size limits.
 Return warnings when content is truncated or redacted.
