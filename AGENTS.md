# AGENTS.md — Working Contract for terraformer v0

> **You are a coding agent (GitHub Copilot in VS Code, backed by Claude Opus or Sonnet) implementing terraformer v0.**
> This file is your operating contract. Read it before every action. The rules here are mechanical, not aspirational. The most important one is [The Two-Commit Rule](#the-two-commit-rule) below — follow it strictly on the honor system.

---

## How to read this repo

Phases and their goals are defined in [PLAN/00-PLAN.md](PLAN/00-PLAN.md). Specifications live in [PLAN/spec/](PLAN/spec/). Read both before starting any phase.

**Order of authority (highest first):**
1. `PLAN/spec/00-spec.md` — non-negotiable safety rules
2. `PLAN/spec/02-mcp-tool-contracts.md` — tool request/response contracts
3. `PLAN/spec/01-testing.md` — TDD requirements
4. `PLAN/spec/03-development.md` — Makefile targets
5. The current phase file in `PLAN/`
6. This file (`AGENTS.md`) — workflow rules
7. Existing code in `internal/`

If anything in this file conflicts with `PLAN/spec/`, the spec wins. Stop and surface the conflict to the human.

---

## The Two-Commit Rule

**This is the single most important rule in this document. Follow it strictly.**

Every meaningful capability requires **two separate commits**:

1. **Red commit:** Adds a failing test. Contains *only* `*_test.go` files (and possibly `testdata/` fixtures).
2. **Green commit:** Adds the production code that makes the test pass. Contains *no* `*_test.go` changes.

If you genuinely need to refactor both at once, the commit message must contain `[refactor]`.

There is no commit hook enforcing this. You are on the honor system. The human will audit your commit history and reject work that violates the rule.

### Workflow per task

```
1. Read the spec lines that define the behavior.
2. Write the test in *_test.go. NOTHING ELSE.
3. Stage and commit:
     git add <test files>
     git commit -m "test(<pkg>): failing test for <behavior>"
   Stage ONLY test files. Do not include implementation files.
4. Run the test. Confirm it FAILS for the expected reason.
   Capture the failure output in your response so the human can verify.
5. Write the production code. NO TEST CHANGES in this step.
6. Stage and commit:
     git add <production files>
     git commit -m "feat(<pkg>): implement <behavior>"
   Stage ONLY production files. Do not include test files.
7. Run the test again. Confirm it PASSES.
8. Run the package tests. Confirm everything is green.
9. Check off the task in the relevant phase file. Commit that as a separate
   docs commit. Do not put the doc update in the red commit or it will appear
   to be done before it actually is.
```

### Self-discipline checks

Before every commit, run:

```bash
git diff --cached --name-only
```

Inspect the output. If the commit is meant to be a **red** commit, every staged file must end in `_test.go` or live under `testdata/`. If it is meant to be a **green** commit, no staged file may end in `_test.go`. If the staged set is wrong, run `git restore --staged <file>` to unstage and fix it before committing.

### What this means in practice

- ❌ You cannot write `runner.go` and `runner_test.go` together and commit them as one. Split into two commits.
- ❌ You cannot edit a test and the implementation in the same commit unless tagged `[refactor]`.
- ✅ You write the test, commit, run it, see red, then write the implementation, commit, run it, see green.
- ✅ Refactors *after* green are allowed in mixed commits when tagged `[refactor]`.

### Reporting back to the human

After every task, your response must include:

1. The commit hash of the red commit.
2. The terminal output showing the test FAILED at the red phase.
3. The commit hash of the green commit.
4. The terminal output showing the test PASSED at the green phase.
5. Output of `make check` (or, if not yet possible, the targeted package test).

If you cannot show all five, the task is not done. Do not check it off in the phase file.

---

## Other workflow rules

- [ ] Before starting a phase, read the phase goal and TDD loop in the phase file.
- [ ] Before implementing, write or identify the failing test (red commit).
- [ ] Do not mark a checkbox complete until the green commit exists and tests pass.
- [ ] Do not mark a checkbox complete if the code is not committed.
- [ ] If implementation reveals a design change, update the relevant spec in the same working session.
- [ ] If a task is deferred, add a short reason in the phase file.
- [ ] If a safety requirement changes, stop and update the safety model first (`PLAN/spec/00-spec.md`).
- [ ] If a golden file changes, document why in the commit message.
- [ ] If an integration test is skipped, document the condition.
- [ ] If a dependency is added, document why in the commit message.
- [ ] If a package boundary changes, document the reason.
- [ ] Never quietly remove a safety requirement.
- [ ] Never claim desired-state comparison works beyond what tests prove.
- [ ] Never add support for `apply`, `destroy`, or arbitrary shell execution in v0.
- [ ] Keep unchecked future work visible.
- [ ] Keep completed work checked only when it is truly done.

---

## Working with GitHub Copilot in VS Code

You are running inside VS Code's Copilot agent mode, backed by Claude Opus or Sonnet. The following constraints apply to your tool usage:

- **Use `run_in_terminal` to run `git` and `go test`.** Do not simulate test runs. Actual terminal output is required for the red/green evidence.
- **Use the editor's file tools (e.g., `replace_string_in_file`, `create_file`) for code changes.** Do not edit through terminal redirection (`echo > file`).
- **Stage files explicitly with `git add <file>`.** Never use `git add -A` or `git add .` — those make it easy to accidentally mix test and implementation files.
- **One logical step per tool exchange.** Do not batch a test write, an implementation write, and a commit into a single edit-and-commit.
- **Always inspect `git diff --cached --name-only` before committing** to confirm the staged set matches the phase (red = tests only; green = production only).

### First actions in any session

The human will likely start a fresh session per phase. Your first actions in any session must be:

1. Read `AGENTS.md` (this file).
2. Read `PLAN/00-PLAN.md`.
3. Read every file in `PLAN/spec/`.
4. Read the current phase file (e.g., `PLAN/01-phase1.md`).
5. Identify the first unchecked task in the phase file.
6. Begin the red/green loop for that task.

### Commit message format

```
test(<pkg>): failing test for <behavior>
feat(<pkg>): implement <behavior>
docs(plan): mark <task> complete in <phase>
refactor(<pkg>): <description> [refactor]
fix(<pkg>): <description>
```

The `<pkg>` is the lowest-level affected package, e.g., `runner`, `safety`, `terraform`, `httpserver`.

---

## Known limitations to preserve until explicitly changed

> These are permanent v0 constraints, not tasks. Do not check them off, and do not implement the forbidden behaviors.

- v0 is local-only.
- v0 does not apply Terraform changes.
- v0 does not destroy Terraform-managed infrastructure.
- v0 does not execute arbitrary shell commands.
- v0 desired-state checking may begin as minimal or stubbed.
- v0 Terraform integration tests may be skipped unless Terraform is installed and explicitly enabled via `TERRAFORMER_RUN_INTEGRATION=1`.
- v0 should prefer conservative safety behavior over convenience.
- Redaction and secret filtering deferred to v0.1.
- Concurrency management deferred to v0.1.
- Secret-looking file exclusion deferred to v0.1.

---

## Failure modes to avoid

These are the specific ways agents fail at this contract. Do not do these.

1. **Writing test and implementation in one edit, then committing them together.** Split into two commits (test first, implementation second).
2. **Claiming a test "would have failed" without running it.** The red phase is observable or it didn't happen. Run the test before writing the implementation.
3. **Stubbing `t.Skip()` or `t.SkipNow()` to fake a red commit.** A skipped test is not a failing test. The human will catch this on audit.
4. **Marking a phase task complete without showing both commit hashes.** A task without a red+green pair is not done.
5. **Editing the spec to make implementation easier.** Specs are the contract. If the spec is wrong, surface it to the human; do not silently rewrite it.
6. **Batching multiple tasks into one red+green pair.** Each meaningful behavior gets its own pair. The phase files list tasks at the right granularity — one pair per task.
7. **Marking a phase complete without running `make check`.** A phase is not done until `make check` passes.
8. **Using `git add -A` or `git add .`.** Always stage files explicitly so you don't accidentally mix tests and implementation.

If you find yourself wanting to bypass the rules, stop and ask the human. The rules exist because agents — including you — reliably fail at TDD when not deliberately constrained.
