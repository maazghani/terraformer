> **Phase 2 of 9** | [← Phase 1](01-phase1.md) | [Index: 00-PLAN.md](00-PLAN.md) | Next: [Phase 3 →](03-phase3.md)

# Phase 2: Repo file access tools with path safety

Implement safe repo-local file listing and reading through the `list_repo_files` and `read_repo_file` internal handlers. Path safety from Phase 1 is now actively used to prevent traversal and symlink escape.

## Phase 2: Repo file access tools with path safety

### Goal

Implement safe repo-local file listing and reading.

### TDD loop

- [x] Write failing tests for repo file listing.
- [x] Write failing tests for reading files.
- [x] Write failing safety tests for traversal and symlink escape.
- [x] Implement smallest repo service code.
- [x] Run targeted repo tests.
- [x] Refactor after tests pass.
- [x] Update this phase file's checklist and the status tracker in [00-PLAN.md](00-PLAN.md).
- [x] Leave code committed or ready to commit.

### Tasks

- [x] Implement `repo.Service`.
- [x] Implement `list_repo_files` internal handler.
- [x] Implement `read_repo_file` internal handler.
- [x] Return normalized relative paths.
- [x] Exclude `.git` by default.
- [x] Exclude `.terraform` from listing by default.
- [x] Enforce `max_files`.
- [x] Enforce `max_bytes`.
- [x] Return `truncated=true` when content is truncated.
- [x] Add golden tests for file response shape.
- [x] Run `go test ./internal/repo ./internal/tools`.

### Completion criteria

- [x] `list_repo_files` behavior is implemented internally.
- [x] `read_repo_file` behavior is implemented internally.
- [x] All repo file access goes through safe path resolution.
- [x] Tests cover traversal, absolute paths, symlink escapes, `.git`, and `.terraform`.
- [x] Secret file handling and redaction deferred to v0.1.
- [x] No HTTP wiring required yet.