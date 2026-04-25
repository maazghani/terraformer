> **Phase 1 of 9** | [Index: 00-PLAN.md](00-PLAN.md) | Next: [Phase 2 →](02-phase2.md)

# Phase 1: Test harness, command runner abstraction, and safety boundaries

Build the command runner abstraction, path safety primitives, and test infrastructure that all later phases depend on. No Terraform execution occurs in this phase—only test doubles and path validation logic.

## Phase 1: Test harness, command runner abstraction, and safety boundaries

### Goal

Build the foundation that all command execution and safety checks depend on.

### TDD loop

- [ ] Write failing unit tests for command runner result shape.
- [ ] Write failing tests for fake runner expectations.
- [ ] Write failing safety tests for repo root validation.
- [ ] Write failing safety tests for path traversal.
- [ ] Implement the smallest runner and safety code needed.
- [ ] Run targeted tests for `internal/runner` and `internal/safety`.
- [ ] Refactor after tests pass.
- [ ] Update this phase file's checklist and the status tracker in [00-PLAN.md](00-PLAN.md).
- [ ] Leave code committed or ready to commit.

### Tasks

- [ ] Define `runner.Command`.
- [ ] Define `runner.Result`.
- [ ] Define `runner.Runner`.
- [ ] Implement fake runner for tests.
- [ ] Implement real process runner.
- [ ] Ensure real runner does not invoke a shell.
- [ ] Ensure real runner captures stdout.
- [ ] Ensure real runner captures stderr.
- [ ] Ensure real runner captures exit code.
- [ ] Ensure real runner captures duration.
- [ ] Ensure real runner respects working directory.
- [ ] Add repo root validation helper.
- [ ] Add safe path resolution helper.
- [ ] Add path traversal tests.
- [ ] Add symlink escape tests.
- [ ] Add forbidden path tests for `.git`.
- [ ] Add default exclusion tests for `.terraform`.
- [ ] Run `go test ./internal/runner`.

### Completion criteria

- [ ] Runner abstraction exists and is testable.
- [ ] Fake runner can assert exact command name, args, working directory, and environment behavior.
- [ ] Path safety rejects traversal and symlink escapes.
- [ ] Redaction helpers deferred to v0.1.
- [ ] No Terraform service exists yet except possibly types.
