> **Phase 5 of 9** | [← Phase 4](04-phase4.md) | [Index: 00-PLAN.md](00-PLAN.md) | Next: [Phase 6 →](06-phase6.md)

# Phase 5: Plan JSON parsing and normalized diagnostics

Parse Terraform JSON outputs (`terraform validate -json`, `terraform show -json`) into stable, normalized structures. Handle malformed JSON gracefully. Redaction deferred to v0.1.

## Phase 5: Plan JSON parsing and normalized diagnostics

### Goal

Parse Terraform JSON outputs into stable, model-friendly structures.

### TDD loop

- [x] Write failing tests for Terraform validate JSON diagnostics.
- [x] Write failing tests for malformed JSON fallback behavior.
- [x] Write failing tests for plan JSON summary parsing.
- [x] Write failing golden tests for normalized diagnostics.
- [x] Implement smallest diagnostics and plan parsing code.
- [x] Run targeted tests.
- [x] Refactor after tests pass.
- [x] Update this phase file's checklist and the status tracker in [00-PLAN.md](00-PLAN.md).
- [x] Leave code committed or ready to commit.

### Tasks

- [x] Define `diagnostics.Diagnostic`.
- [x] Parse `terraform validate -json` diagnostics.
- [x] Normalize severity.
- [x] Normalize summary.
- [x] Normalize detail.
- [x] Normalize file/range if available.
- [x] Add fallback diagnostic for plain stderr.
- [x] Define plan JSON summary type.
- [x] Parse resource changes from `terraform show -json`.
- [x] Count create, update, delete, replace, and no-op actions.
- [ ] Add response size limit behavior.
- [x] Add golden files for representative diagnostics.
- [x] Add golden files for representative plan summaries.
- [x] Run `go test ./internal/diagnostics ./internal/terraform`.

### Completion criteria

- [x] Diagnostics are stable and tested.
- [x] Plan summaries are stable and tested.
- [x] Malformed JSON produces a safe response, not a panic.
- [x] Golden outputs are intentional and documented.
- [x] Redaction deferred to v0.1.
