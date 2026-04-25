> **Phase 2 of 9** | [← Phase 1](01-phase1.md) | [Index: 00-PLAN.md](00-PLAN.md) | Next: [Phase 3 →](03-phase3.md)

# Phase 2: Repo file access tools with path safety

Implement safe repo-local file listing and reading through the `list_repo_files` and `read_repo_file` internal handlers. Path safety from Phase 1 is now actively used to prevent traversal and symlink escape.

## Phase 2: Repo file access tools with path safety

### Goal

Implement safe repo-local file listing and reading.

### TDD loop

- [ ] Write failing tests for repo file listing.
- [ ] Write failing tests for reading files.
- [ ] Write failing safety tests for traversal and symlink escape.
- [ ] Implement smallest repo service code.
- [ ] Run targeted repo tests.
- [ ] Refactor after tests pass.
- [ ] Update this phase file's checklist and the status tracker in [00-PLAN.md](00-PLAN.md).
- [ ] Leave code committed or ready to commit.

### Tasks

- [ ] Implement `repo.Service`.
- [ ] Implement `list_repo_files` internal handler.
- [ ] Implement `read_repo_file` internal handler.
- [ ] Return normalized relative paths.
- [ ] Exclude `.git` by default.
- [ ] Exclude `.terraform` from listing by default.
- [ ] Enforce `max_files`.
- [ ] Enforce `max_bytes`.
- [ ] Return `truncated=true` when content is truncated.
- [ ] Add golden tests for file response shape.
- [ ] Run `go test ./internal/repo ./internal/tools`.

### Completion criteria

- [ ] `list_repo_files` behavior is implemented internally.
- [ ] `read_repo_file` behavior is implemented internally.
- [ ] All repo file access goes through safe path resolution.
- [ ] Tests cover traversal, absolute paths, symlink escapes, `.git`, and `.terraform`.
- [ ] Secret file handling and redaction deferred to v0.1.
- [ ] No HTTP wiring required yet.