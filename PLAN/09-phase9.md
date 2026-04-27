> **Phase 9 of 9** | [← Phase 8](08-phase8.md) | [Index: 00-PLAN.md](00-PLAN.md)

# Phase 9: Documentation, hardening, and demo workflow

Phases 1–8 produced a working v0 server: all nine tools are wired through the HTTP/JSON server on port 9001 (configurable), structured JSON logs are emitted, the repo root is validated at startup and immutable thereafter, forbidden commands are blocked, and end-to-end tests pass with the fake runner. This phase is **about turning that into something a human (or another agent) can actually pick up and use**: write the user-facing docs, fill the small remaining hardening gaps, and ship a demo that matches real behavior.

## Phase 9: Documentation, hardening, and demo workflow

### Goal

Document real v0 usage, close the remaining hardening gaps (CLI surface, response size limits, config validation), and provide a demo workflow plus per-client integration instructions that match the actual implementation.

### Already implemented (do not redo)

These were completed in earlier phases. Verify them while writing docs; do not re-checkbox them here.

- HTTP server on port 9001, configurable via `--port`.
- `--repo-root` flag, required, validated at startup, immutable after startup.
- Structured JSON logs to stdout/stderr (see [internal/httpserver/logger.go](../internal/httpserver/logger.go)).
- All v0 tools registered: `terraform_init`, `terraform_fmt`, `terraform_validate`, `terraform_plan`, `terraform_show_json`, `list_repo_files`, `read_repo_file`, `apply_patch`, `check_desired_state`.
- Forbidden Terraform subcommands rejected (`apply`, `destroy`, `import`, `workspace`, `state`, …).
- Path traversal and symlink escape prevention.
- `make test`, `make test-unit`, `make test-integration`, `make check`, `make build`.
- Integration test opt-in via `TERRAFORMER_RUN_INTEGRATION=1`.
- End-to-end fake-runner test in [internal/httpserver/e2e_test.go](../internal/httpserver/e2e_test.go).

### TDD loop

- [ ] Write failing tests for each new CLI flag's parsing and validation.
- [ ] Write failing tests for response size limit truncation.
- [ ] Write failing tests for startup config validation (repo root must exist, must be a directory, must be absolute).
- [ ] Implement smallest hardening change per red commit (Two-Commit Rule applies).
- [ ] Run targeted package tests after each green commit.
- [ ] Run `make check` at the end of the phase.
- [ ] Refactor only after tests pass.
- [ ] Update this phase file's checklist and the status tracker in [00-PLAN.md](00-PLAN.md).
- [ ] Leave code committed.

### Hardening tasks (code)

These are the small remaining gaps from earlier phases. Each requires a red+green commit pair.

- [ ] Add `--log-level` flag (default `info`; values `debug|info|warn|error`). Wire it through `httpserver.Config` and the JSON logger.
- [ ] Add `--max-response-bytes` flag (default `1048576` = 1 MiB). Truncate oversized tool response bodies and set a `truncated: true` marker in the JSON response.
- [ ] Add `--terraform-bin` flag (default `terraform`). Thread it through `runner.LocalRunner` / `terraform.Service` so a custom binary path can be used.
- [ ] Strengthen startup config validation in `cmd/terraformer-mcp`:
  - [ ] Repo root must be an absolute path.
  - [ ] Repo root must exist on disk.
  - [ ] Repo root must be a directory.
  - [ ] Port must be in `1..65535`.
  - [ ] On any failure, exit non-zero with a single-line JSON error to stderr.
- [ ] Confirm `--help` output (provided by `flag`) is readable and lists every flag with its default. Add a test that exercises `flag.CommandLine` if practical.

### Documentation tasks

> **Documentation must match shipped behavior.** Do not document flags, tools, or guarantees that are not in the code. If a doc claim cannot be backed by a test or by reading the source, delete it.

- [ ] Write `README.md` at the repo root covering, in this order:
  - [ ] One-paragraph summary of what terraformer v0 is.
  - [ ] Explicit non-goals: no `apply`, no `destroy`, no arbitrary shell, local-only, no remote backends mutated.
  - [ ] Build instructions (`make build`).
  - [ ] Quickstart: start the server against a fixture repo, hit one tool with `curl`.
  - [ ] Pointer to [USAGE.md](../USAGE.md) for client integration.
  - [ ] Pointer to [PLAN/spec/02-mcp-tool-contracts.md](spec/02-mcp-tool-contracts.md) for tool contracts.
  - [ ] Test commands: `make test`, `make test-unit`, `make test-integration`, `make check`.
  - [ ] Integration test opt-in note (`TERRAFORMER_RUN_INTEGRATION=1`, requires `terraform` on `PATH`).
  - [ ] Known limitations section (see below).
- [ ] Write `USAGE.md` at the repo root with **tabbed per-client integration instructions** covering at minimum:
  - [ ] GitHub Copilot (VS Code agent mode, `mcp.json` / `settings.json` snippet).
  - [ ] Cursor (`~/.cursor/mcp.json` snippet).
  - [ ] Claude Code / Claude Desktop (`claude_desktop_config.json` snippet).
  - [ ] Codex CLI (`~/.codex/config.toml` snippet).
  - [ ] Cline (VS Code extension settings snippet).
  - [ ] Use GitHub-flavored markdown collapsible sections (`<details><summary>`) to emulate tabs and keep the page navigable.
  - [ ] Each tab shows: how to launch terraformer-mcp, how to point the client at `http://localhost:9001`, and a smoke-test prompt.
  - [ ] Each tab links back to the canonical tool contracts in [PLAN/spec/02-mcp-tool-contracts.md](spec/02-mcp-tool-contracts.md).
- [ ] Document each tool's request/response shape with one minimal example per tool. Prefer linking into the spec file rather than duplicating it.
- [ ] Document the example iterative workflow: `list_repo_files` → `read_repo_file` → `apply_patch` → `terraform_fmt` → `terraform_validate` → `terraform_plan` → `check_desired_state`.
- [ ] Document forbidden operations and link to [internal/safety/safety.go](../internal/safety/safety.go).
- [ ] Document local-only execution (no network calls except whatever the user's Terraform providers themselves make).

### Demo workflow

- [ ] Add `scripts/demo.sh` that:
  - [ ] Builds `terraformer-mcp`.
  - [ ] Starts it against `testdata/fixtures/valid-basic` on a free port.
  - [ ] Issues a representative sequence of `curl` calls against the running server.
  - [ ] Tears the server down cleanly.
- [ ] Add a short `scripts/README.md` (or section in `README.md`) describing the demo and its prerequisites.
- [ ] The demo must not depend on any unimplemented tool, flag, or behavior.

### Known limitations to record verbatim

These belong in `README.md` under a clearly-labeled section. They are not tasks to fix in v0.

- [ ] v0 is local-only.
- [ ] v0 does not apply or destroy infrastructure.
- [ ] v0 does not execute arbitrary shell commands.
- [ ] `check_desired_state` reports its implementation level honestly and is intentionally minimal in v0.
- [ ] No output redaction in v0; secrets in Terraform output may reach the model. (Deferred to v0.1.)
- [ ] No concurrency management; concurrent tool calls share the repo root. (Deferred to v0.1.)
- [ ] No secret-looking file exclusion in `list_repo_files` / `read_repo_file`. (Deferred to v0.1.)

### Completion criteria

- [ ] `README.md` exists, describes only implemented behavior, and links to `USAGE.md` and the spec.
- [ ] `USAGE.md` exists with working integration snippets for GitHub Copilot, Cursor, Claude, Codex, and Cline.
- [ ] `--log-level`, `--max-response-bytes`, and `--terraform-bin` flags exist, are tested, and are documented.
- [ ] Startup config validation rejects bad repo roots and bad ports with a clear JSON error.
- [ ] `scripts/demo.sh` runs against the bundled fixtures end-to-end without errors.
- [ ] `make check` passes.
- [ ] Integration tests pass when `TERRAFORMER_RUN_INTEGRATION=1` and `terraform` is installed; otherwise they skip cleanly.
- [ ] Status tracker in [00-PLAN.md](00-PLAN.md) updated to mark Phase 9 complete.
- [ ] Redaction, concurrency management, and secret-file filtering are explicitly recorded in `README.md` as v0.1 work.
