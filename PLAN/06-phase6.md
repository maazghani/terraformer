> **Phase 6 of 9** | [← Phase 5](05-phase5.md) | [Index: 00-PLAN.md](00-PLAN.md) | Next: [Phase 7 →](07-phase7.md)

# Phase 6: Patch application workflow

Implement the `apply_patch` tool using a structured JSON format. Patches can create, write, and delete repo-local files. All operations validate paths and prevent traversal. No Terraform execution occurs.

## Phase 6: Patch application workflow

### Goal

Allow the agent to modify repo-local files through a safe patch tool.

### TDD loop

- [ ] Write failing tests for valid patch application.
- [ ] Write failing tests for invalid patch handling.
- [ ] Write failing tests for patch path traversal.
- [ ] Write failing tests for symlink escape.
- [ ] Write failing tests proving patch application does not run Terraform.
- [ ] Implement smallest patch service.
- [ ] Run targeted tests.
- [ ] Refactor after tests pass.
- [ ] Update this phase file's checklist and the status tracker in [00-PLAN.md](00-PLAN.md).
- [ ] Leave code committed or ready to commit.

### Tasks

- [ ] Define `patch.Service`.
- [ ] Define structured patch request and response.
- [ ] Support structured JSON patch format: array of `{path, operation, content}` objects.
- [ ] Support "write" operations to create/overwrite files.
- [ ] Support "delete" operations for file removal (if tested).
- [ ] Validate all target paths before writing.
- [ ] Apply patches to existing files.
- [ ] Add support for creating new files if safe and tested.
- [ ] Return changed files.
- [ ] Return rejected files.
- [ ] Fail atomically where practical.
- [ ] Avoid command execution during patching.
- [ ] Add golden response tests.
- [ ] Run `go test ./internal/patch ./internal/repo`.

### Completion criteria

- [ ] `apply_patch` is implemented internally.
- [ ] Patches cannot escape the repo.
- [ ] Patch failures are structured.
- [ ] Patch application does not run Terraform or any shell command.