> **Phase 8 of 9** | [← Phase 7](07-phase7.md) | [Index: 00-PLAN.md](00-PLAN.md) | Next: [Phase 9 →](09-phase9.md)

# Phase 8: HTTP server wiring and end-to-end loop

Wire all tools through an HTTP/JSON server listening on port 9001 (configurable). Emit structured JSON logs to stdout/stderr. Add end-to-end tests. Docker-ready: can be containerized and deployed.

## Phase 8: HTTP server wiring and end-to-end loop

### Goal

Expose the implemented tools through an HTTP/JSON server on port 9001.

### TDD loop

- [ ] Write failing tests for tool registration.
- [ ] Write failing tests for request decoding.
- [ ] Write failing tests for response encoding.
- [ ] Write failing tests for invalid request handling.
- [ ] Write failing end-to-end test using fake runner and temp repo.
- [ ] Implement smallest HTTP server wiring.
- [ ] Run targeted HTTP server tests.
- [ ] Refactor after tests pass.
- [ ] Update this phase file's checklist and the status tracker in [00-PLAN.md](00-PLAN.md).
- [ ] Leave code committed or ready to commit.

### Tasks

- [ ] Choose HTTP/JSON-RPC library approach.
- [ ] Keep HTTP server details inside internal/httpserver.
- [ ] Start HTTP server on port 9001 (configurable).
- [ ] Register terraform_init.
- [ ] Register terraform_fmt.
- [ ] Register terraform_validate.
- [ ] Register terraform_plan.
- [ ] Register terraform_show_json.
- [ ] Register list_repo_files.
- [ ] Register read_repo_file.
- [ ] Register apply_patch.
- [ ] Register check_desired_state.
- [ ] Add typed tool schemas.
- [ ] Add request validation errors.
- [ ] Add structured JSON response serialization.
- [ ] Add structured JSON logging to stdout and stderr.
- [ ] Add end-to-end loop test:
  - [ ] List files.
  - [ ] Read file.
  - [ ] Patch file.
  - [ ] Run fmt.
  - [ ] Run validate.
  - [ ] Run plan with fake runner.
  - [ ] Check desired state.
- [ ] Ensure startup requires a repo root.
- [ ] Ensure repo root cannot be changed by HTTP request.
- [ ] Note: Concurrency management deferred to v0.1.
- [ ] Run `go test ./internal/httpserver ./internal/tools ./...`.

### Completion criteria

- [ ] HTTP server starts on port 9001.
- [ ] All v0 tools are registered.
- [ ] Tool calls use structured requests and responses.
- [ ] Structured JSON logs are emitted to stdout/stderr.
- [ ] End-to-end fake-runner test passes.
- [ ] No unsafe commands are reachable through HTTP.
- [ ] Docker-ready: can be containerized and accessed via port 9001.