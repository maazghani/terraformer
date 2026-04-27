> **Phase 1 of 9** | [Index: 00-PLAN.md](00-PLAN.md) | Next: [Phase 2 →](02-phase2.md)

# Phase 1: Test harness, command runner abstraction, and safety boundaries

Build the command runner abstraction, path safety primitives, and test infrastructure that all later phases depend on. No Terraform execution occurs in this phase—only test doubles and path validation logic.

## Phase 1: Test harness, command runner abstraction, and safety boundaries

### Goal

Build the foundation that all command execution and safety checks depend on.

### TDD loop

- [x] Write failing unit tests for command runner result shape.
- [x] Write failing tests for fake runner expectations.
- [x] Write failing safety tests for repo root validation.
- [x] Write failing safety tests for path traversal.
- [x] Implement the smallest runner and safety code needed.
- [x] Run targeted tests for `internal/runner` and `internal/safety`.
- [x] Refactor after tests pass.
- [x] Update this phase file's checklist and the status tracker in [00-PLAN.md](00-PLAN.md).
- [x] Leave code committed or ready to commit.

### Tasks

- [x] Define `runner.Command`.
- [x] Define `runner.Result`.
- [x] Define `runner.Runner`.
- [x] Implement fake runner for tests.
- [x] Implement real process runner.
- [x] Ensure real runner does not invoke a shell.
- [x] Ensure real runner captures stdout.
- [x] Ensure real runner captures stderr.
- [x] Ensure real runner captures exit code.
- [x] Ensure real runner captures duration.
- [x] Ensure real runner respects working directory.
- [x] Add repo root validation helper.
- [x] Add safe path resolution helper.
- [x] Add path traversal tests.
- [x] Add symlink escape tests.
- [x] Add forbidden path tests for `.git`.
- [x] Add default exclusion tests for `.terraform`.
- [x] Run `go test ./internal/runner`.

### Completion criteria

- [x] Runner abstraction exists and is testable.
- [x] Fake runner can assert exact command name, args, working directory, and environment behavior.
- [x] Path safety rejects traversal and symlink escapes.
- [x] Redaction helpers deferred to v0.1.
- [x] No Terraform service exists yet except possibly types.
