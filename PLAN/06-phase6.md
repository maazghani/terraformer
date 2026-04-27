> **Phase 6 of 9** | [← Phase 5](05-phase5.md) | [Index: 00-PLAN.md](00-PLAN.md) | Next: [Phase 7 →](07-phase7.md)

# Phase 6: Patch application workflow

Implement the `apply_patch` tool using a structured JSON format. Patches can create, write, and delete repo-local files. All operations validate paths and prevent traversal. No Terraform execution occurs.

## Phase 6: Patch application workflow

### Goal

Allow the agent to modify repo-local files through a safe patch tool.

### TDD loop

- [x] Write failing tests for valid patch application.
- [x] Write failing tests for invalid patch handling.
- [x] Write failing tests for patch path traversal.
- [x] Write failing tests for symlink escape.
- [x] Write failing tests proving patch application does not run Terraform.
- [x] Implement smallest patch service.
- [x] Run targeted tests.
- [x] Refactor after tests pass.
- [x] Update this phase file's checklist and the status tracker in [00-PLAN.md](00-PLAN.md).
- [x] Leave code committed or ready to commit.

### Tasks

- [x] Define `patch.Service`.
- [x] Define structured patch request and response.
- [x] Support structured JSON patch format: array of `{path, operation, content}` objects.
- [x] Support "write" operations to create/overwrite files.
- [x] Support "delete" operations for file removal (if tested).
- [x] Validate all target paths before writing.
- [x] Apply patches to existing files.
- [x] Add support for creating new files if safe and tested.
- [x] Return changed files.
- [x] Return rejected files.
- [x] Fail atomically where practical.
- [x] Avoid command execution during patching.
- [ ] Add golden response tests.
- [x] Run `go test ./internal/patch ./internal/repo`.

### Completion criteria

- [x] `apply_patch` is implemented internally.
- [x] Patches cannot escape the repo.
- [x] Patch failures are structured.
- [x] Patch application does not run Terraform or any shell command.