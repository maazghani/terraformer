> **Phase 4 of 9** | [← Phase 3](03-phase3.md) | [Index: 00-PLAN.md](00-PLAN.md) | Next: [Phase 5 →](05-phase5.md)

# Phase 4: Real Terraform integration tests using local fixtures

Add integration tests that exercise real Terraform against credential-free, committed fixtures. Tests are skipped by default unless `TERRAFORMER_RUN_INTEGRATION=1` is set. No provider credentials required.

## Phase 4: Real Terraform integration tests using local fixtures

### Goal

Add integration tests that run real Terraform against local fixtures, skipped by default.

### TDD loop

- [x] Write failing integration test scaffold that skips when not enabled.
- [x] Add fixtures copied into temporary dirs before execution.
- [x] Write integration test for valid basic Terraform.
- [x] Write integration test for invalid HCL.
- [x] Write integration test for missing provider or init/validate failure.
- [x] Implement any missing fixture helpers.
- [x] Run targeted integration tests with opt-in environment variable.
- [x] Refactor after tests pass.
- [x] Update this phase file's checklist and the status tracker in [00-PLAN.md](00-PLAN.md).
- [x] Leave code committed or ready to commit.

### Tasks

- [x] Add `testdata/fixtures/valid-basic`.
- [x] Add `testdata/fixtures/invalid-hcl`.
- [x] Add `testdata/fixtures/missing-provider`.
- [x] Add `testdata/fixtures/plan-basic`.
- [x] Add fixture copy helper.
- [x] Add Terraform availability check helper.
- [x] Add integration test skip behavior when `TERRAFORMER_RUN_INTEGRATION` is unset.
- [x] Test `terraform fmt` on fixture.
- [x] Test `terraform validate` on valid fixture.
- [x] Test `terraform validate` on invalid fixture.
- [x] Test `terraform plan` behavior where practical.
- [x] Ensure integration tests do not require provider credentials.
- [x] Ensure integration tests do not mutate committed fixtures.
- [x] Run `TERRAFORMER_RUN_INTEGRATION=1 go test ./... -run Integration`.

### Completion criteria

- [x] Integration tests are skipped by default.
- [x] Integration tests pass when enabled in an environment with Terraform installed.
- [x] Fixtures are local and credential-free.
- [x] No integration test calls `apply` or `destroy`.
