> **Phase 3 of 9** | [← Phase 2](02-phase2.md) | [Index: 00-PLAN.md](00-PLAN.md) | Next: [Phase 4 →](04-phase4.md)

# Phase 3: Terraform command tools using fake runner first

Implement all v0 Terraform tool logic (`init`, `fmt`, `validate`, `plan`, `show -json`) using the fake runner from Phase 1. This phase proves the command structure and semantics without requiring Terraform to be installed.

## Phase 3: Terraform command tools using fake runner first

### Goal

Implement Terraform tool logic using fake runner tests before any real Terraform execution.

### TDD loop

- [x] Write failing tests for `terraform_init` command construction.
- [x] Write failing tests for `terraform_fmt` command construction.
- [x] Write failing tests for `terraform_validate` command construction.
- [x] Write failing tests for `terraform_plan` command construction and exit-code normalization.
- [x] Write failing tests for `terraform_show_json` command construction.
- [x] Write failing tests proving forbidden commands cannot be constructed.
- [x] Implement smallest Terraform service code.
- [x] Run targeted Terraform tests.
- [x] Refactor after tests pass.
- [x] Update this phase file's checklist and the status tracker in [00-PLAN.md](00-PLAN.md).
- [x] Leave code committed or ready to commit.

### Tasks

- [x] Create `terraform.Service`.
- [x] Add typed request structs for init, fmt, validate, plan, and show JSON.
- [x] Add typed response structs for command results.
- [x] Ensure all Terraform commands use `runner.Runner`.
- [x] Implement `terraform_init`.
- [x] Implement `terraform_fmt`.
- [x] Implement `terraform_validate`.
- [x] Implement `terraform_plan`.
- [x] Implement `terraform_show_json`.
- [x] Normalize plan detailed exit codes.
- [x] Reject unsafe plan output paths.
- [x] Reject unsafe show JSON plan paths.
- [ ] Prove `apply` cannot be invoked through the service.
- [ ] Prove `destroy` cannot be invoked through the service.
- [ ] Prove arbitrary shell commands cannot be invoked through the service.
- [ ] Return stdout, stderr, exit code, duration, diagnostics, and warnings.
- [ ] Run `go test ./internal/terraform ./internal/runner`.

### Completion criteria

- [ ] All v0 Terraform commands are implemented behind fake runner tests.
- [ ] No unit test requires Terraform installed.
- [ ] Forbidden commands are impossible or rejected by tested code.
- [ ] Plan success is not treated as desired-state success.
- [ ] Redaction deferred to v0.1.