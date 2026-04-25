# Specification: Testing Philosophy and Requirements

This document defines the TDD workflow, test categories, and testing requirements for all implementation phases. Every meaningful capability must follow the prescribed TDD loop and test pattern.

## Testing philosophy
TDD is mandatory.
The default loop for every meaningful change is:
Write a failing test.
Run the targeted test and confirm it fails for the expected reason.
Implement the smallest code needed to pass.
Run the targeted test again.
Run the relevant package tests.
Refactor only after tests pass.
Run tests again.
Update the current phase file's checklist and the status tracker in PLAN/00-PLAN.md.
Commit or leave the code ready to commit.
Do not implement a capability first and add tests later. That is backwards. It creates soft code.
Tests should favor small, explicit fixtures and deterministic behavior. Do not depend on cloud access. Do not depend on provider credentials. Real Terraform integration tests must use local fixtures only.

## Test categories
### Unit tests
Purpose:
 Test individual packages without invoking Terraform.
 Test request and response validation.
 Test path safety.
 Test command construction.
 Test diagnostics parsing.
 Test desired-state comparison behavior.
Expected locations:
internal/**/**/*_test.go
Command runner tests
Purpose:
 Verify real runner captures stdout.
 Verify real runner captures stderr.
 Verify exit code capture.
 Verify duration is recorded.
 Verify commands run in the configured working directory.
 Verify no shell interpolation occurs.
 Verify fake runner can assert command arguments.
Important:
 Unit tests for Terraform tools must use fake runners.
 Integration tests may use the real runner.
### Safety tests
Purpose:
 Reject ../.
 Reject absolute paths outside repo root.
 Reject symlink escapes.
 Reject forbidden Terraform commands.
 Reject arbitrary shell execution.
 Ensure commands cannot set working directory outside repo root.
Note: Redaction tests deferred to v0.1.

### Integration tests
Purpose:
 Exercise package boundaries together.
 Run real Terraform only when explicitly enabled.
 Validate end-to-end tool behavior using local fixtures.
Integration test requirements:
 Must not require network access unless explicitly marked and skipped by default.
 Must not require provider credentials.
 Must be skipped cleanly when Terraform is unavailable.
 Must use temporary copied fixtures, not mutate committed fixtures.

### Golden file tests
Purpose:
 Validate normalized diagnostics output.
 Validate plan JSON parsing output.
 Validate MCP tool response shape.
 Prevent accidental response contract drift.
Golden files live under:
testdata/golden/
Golden file updates must be intentional and documented in PLAN.md.

### Terraform fixture tests
Purpose:
 Provide realistic repo-local Terraform examples.
 Test valid HCL.
 Test invalid HCL.
 Test missing provider behavior.
 Test plan JSON parsing with stable fixture data.
Fixtures live under:
testdata/fixtures/
Fixtures should be copied into a temp directory before modification.

## See also

- [00-spec.md](00-spec.md) — Safety invariants, repo structure, and package responsibilities
- [02-mcp-tool-contracts.md](02-mcp-tool-contracts.md) — Tool contracts that each tool's tests must satisfy
- [03-development.md](03-development.md) — Makefile commands for running tests (make test, make test-unit, make test-integration)
