# PLAN.md: terraformer

**This document orchestrates the entire v0 implementation.** It defines the sequential phases (1-9), references the non-negotiable specification documents, and tracks completion status. All agents must read the specifications before starting any phase and re-consult them when making architectural or implementation decisions.

## Instructions for Coding Agent

Work through the phases below **in strict order**, completing each phase fully before moving to the next.

### Rules

All work must comply with the specifications in the `spec/` directory. Treat every file in `spec/` as a non-negotiable set of rules and constraints. Read them before starting any phase and re-consult them whenever making architectural or implementation decisions.

- [spec/00-spec.md](spec/00-spec.md)
- [spec/01-testing.md](spec/01-testing.md)
- [spec/02-mcp-tool-contracts.md](spec/02-mcp-tool-contracts.md)
- [spec/03-development.md](spec/03-development.md)

### Phases

Execute phases in this order. Do not begin a phase until the previous one is complete and verified.

1. [01-phase1.md](01-phase1.md)
2. [02-phase2.md](02-phase2.md)
3. [03-phase3.md](03-phase3.md)
4. [04-phase4.md](04-phase4.md)
5. [05-phase5.md](05-phase5.md)
6. [06-phase6.md](06-phase6.md)
7. [07-phase7.md](07-phase7.md)
8. [08-phase8.md](08-phase8.md)
9. [09-phase9.md](09-phase9.md)

## Detailed TDD checklist template for every task

> **Template only.** Copy this checklist into the relevant phase file for each capability. Do not check off these boxes here.

Use this checklist for every meaningful capability, no exceptions.

- [ ] Identify the behavior.
- [ ] Write the smallest failing test for that behavior.
- [ ] Confirm the test fails for the expected reason.
- [ ] Implement the smallest code needed.
- [ ] Run the targeted test.
- [ ] Add edge-case tests.
- [ ] Run the package tests.
- [ ] Refactor only after green tests.
- [ ] Run the package tests again.
- [ ] Update docs or comments if the behavior is part of the public contract.
- [ ] Update this phase file's checklist and the status tracker in [PLAN/00-PLAN.md](PLAN/00-PLAN.md).
- [ ] Leave code committed or ready to commit.

Do not batch large untested changes. Large untested changes rot fast. Keep the work tight.
## Definition of done for a capability

A capability is done only when all of the following are true:

- [ ] Tests were written before or alongside the implementation.
- [ ] Tests cover normal behavior.
- [ ] Tests cover failure behavior.
- [ ] Safety-relevant edge cases are covered.
- [ ] The implementation is small and understandable.
- [ ] No forbidden command or path escape is introduced.
- [ ] Outputs are structured.
- [ ] Secrets are redacted where output may reach the model.
- [ ] Targeted tests pass.
- [ ] Relevant package tests pass.
- [ ] go test ./... passes unless a documented temporary exception exists.
- [ ] Code is gofmt-formatted.
- [ ] Code is committed or ready to commit.
- [ ] This phase file's checklist and the status tracker in [PLAN/00-PLAN.md](PLAN/00-PLAN.md) are updated.


# Definition of done for v0

terraformer v0 is done when:

- [ ] The MCP server can be started for a configured local repo root.
- [ ] The repo root is immutable after startup.
- [ ] The server exposes all expected v0 tools:
  - [ ] terraform_init
  - [ ] terraform_fmt
  - [ ] terraform_validate
  - [ ] terraform_plan
  - [ ] terraform_show_json
  - [ ] list_repo_files
  - [ ] read_repo_file
  - [ ] apply_patch
  - [ ] check_desired_state
- [ ] All tools use structured request and response objects.
- [ ] Terraform commands return stdout, stderr, exit code, duration, and diagnostics where possible.
- [ ] Terraform command execution uses the runner abstraction.
- [ ] Unit tests do not require Terraform.
- [ ] Integration tests are opt-in and use local fixtures.
- [ ] Path traversal is prevented.
- [ ] Symlink escapes are prevented.
- [ ] Forbidden Terraform commands cannot be executed.
- [ ] Arbitrary shell commands cannot be executed.
- [ ] Redaction and secret filtering deferred to v0.1.
- [ ] terraform plan success is not treated as desired-state success.
- [ ] check_desired_state exists and honestly reports its implemented level of functionality.
- [ ] make check passes.
- [ ] README.md documents actual behavior.
- [ ] PLAN.md reflects actual completed and remaining work.

# OVERALL STATUS

# Current status

- [ ] Phase 1 in progress.
- [ ] Phase 2 pending.
- [x] Phase 3 complete.
- [x] Phase 4 complete.
- [x] Phase 5 complete.
- [x] Phase 6 complete.
- [x] Phase 7 complete.
- [ ] Phase 8 in progress (all tasks done except Docker; concurrency deferred to v0.1).
- [ ] Phase 9 pending.

Update this section after every meaningful implementation session.
