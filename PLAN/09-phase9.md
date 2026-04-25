> **Phase 9 of 9** | [← Phase 8](08-phase8.md) | [Index: 00-PLAN.md](00-PLAN.md)

# Phase 9: Documentation, hardening, and demo workflow

Write documentation matching actual v0 behavior. Define CLI flags, validate config, set response size limits. Run full test suite. Document known limitations and deferred work (redaction, concurrency, secret filtering for v0.1).

## Phase 9: Documentation, hardening, and demo workflow

### Goal

Document real usage, harden rough edges, and provide a demo that matches actual implemented behavior.

### TDD loop

- [ ] Write failing tests for documented CLI flags where practical.
- [ ] Write failing tests for config validation.
- [ ] Write failing tests for response size limits.
- [ ] Implement smallest hardening changes.
- [ ] Run full test suite.
- [ ] Refactor only after tests pass.
- [ ] Update this phase file's checklist and the status tracker in [00-PLAN.md](00-PLAN.md).
- [ ] Leave code committed or ready to commit.

### Tasks

- [ ] Write README.md with accurate v0 behavior.
- [ ] Document forbidden operations.
- [ ] Document local-only execution.
- [ ] Document HTTP server and available tools.
- [ ] Document example iterative workflow.
- [ ] Document test commands (make test, make test-integration).
- [ ] Document integration test opt-in (TERRAFORMER_RUN_INTEGRATION=1).
- [ ] Add CLI flags:
  - [ ] `--repo-root` (required, validated at startup)
  - [ ] `--port` (default 9001)
  - [ ] `--log-level` (default info, options: debug, info, warn, error)
  - [ ] `--max-response-bytes` (default 1MB)
  - [ ] `--terraform-bin` (default terraform, path to terraform binary)
- [ ] Add CLI usage help (--help).
- [ ] Add config validation (repo root exists, is directory, etc.).
- [ ] Add response size limits (truncate large outputs).
- [ ] Add structured JSON logging setup to stdout/stderr.
- [ ] Add demo fixture workflow under scripts/ if useful.
- [ ] Run make check.
- [ ] Run integration tests if Terraform is available.
- [ ] Record known limitations (no redaction v0, no concurrency mgmt, no secret file filtering).
- [ ] Note: Redaction deferred to v0.1.
- [ ] Note: Concurrency management deferred to v0.1.

### Completion criteria

- [ ] Documentation matches actual behavior.
- [ ] Demo does not rely on unimplemented features.
- [ ] Full unit test suite passes.
- [ ] Integration tests are either passing or cleanly skipped.
- [ ] Plan is updated with completed work and known limitations.
- [ ] Redaction and concurrency management documented as v0.1 work.
