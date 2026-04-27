> **Phase 8 of 9** | [← Phase 7](07-phase7.md) | [Index: 00-PLAN.md](00-PLAN.md) | Next: [Phase 9 →](09-phase9.md)

# Phase 8: HTTP server wiring and end-to-end loop

Wire all tools through an HTTP/JSON server listening on port 9001 (configurable). Emit structured JSON logs to stdout/stderr. Add end-to-end tests. Docker-ready: can be containerized and deployed.

## Phase 8: HTTP server wiring and end-to-end loop

### Goal

Expose the implemented tools through an HTTP/JSON server on port 9001.

### TDD loop

- [x] Write failing tests for tool registration.
- [x] Write failing tests for request decoding.
- [x] Write failing tests for response encoding.
- [x] Write failing tests for invalid request handling.
- [x] Write failing end-to-end test using fake runner and temp repo.
- [x] Implement smallest HTTP server wiring.
- [x] Run targeted HTTP server tests.
- [x] Refactor after tests pass.
- [x] Update this phase file's checklist and the status tracker in [00-PLAN.md](00-PLAN.md).
- [x] Leave code committed or ready to commit.

### Tasks

- [x] Choose HTTP/JSON-RPC library approach.
- [x] Keep HTTP server details inside internal/httpserver.
- [x] Start HTTP server on port 9001 (configurable).
- [x] Register terraform_init.
- [x] Register terraform_fmt.
- [x] Register terraform_validate.
- [x] Register terraform_plan.
- [x] Register terraform_show_json.
- [x] Register list_repo_files.
- [x] Register read_repo_file.
- [x] Register apply_patch.
- [x] Register check_desired_state.
- [x] Add typed tool schemas.
- [x] Add request validation errors.
- [x] Add structured JSON response serialization.
- [x] Add structured JSON logging to stdout and stderr.
- [x] Add end-to-end loop test:
  - [x] List files.
  - [x] Read file.
  - [x] Patch file.
  - [x] Run fmt.
  - [x] Run validate.
  - [x] Run plan with fake runner.
  - [x] Check desired state.
- [x] Ensure startup requires a repo root.
- [x] Ensure repo root cannot be changed by HTTP request.
- [ ] Note: Concurrency management deferred to v0.1.
- [x] Run `go test ./internal/httpserver ./internal/tools ./...`.

### Completion criteria

- [x] HTTP server starts on port 9001.
- [x] All v0 tools are registered.
- [x] Tool calls use structured requests and responses.
- [x] Structured JSON logs are emitted to stdout/stderr.
- [x] End-to-end fake-runner test passes.
- [x] No unsafe commands are reachable through HTTP.
- [ ] Docker-ready: can be containerized and accessed via port 9001.