> **Phase 4 of 9** | [← Phase 3](03-phase3.md) | [Index: 00-PLAN.md](00-PLAN.md) | Next: [Phase 5 →](05-phase5.md)

# Phase 4: Real Terraform integration tests using local fixtures

Add integration tests that exercise real Terraform against credential-free, committed fixtures. Tests are skipped by default unless `TERRAFORMER_RUN_INTEGRATION=1` is set. No provider credentials required.

## Phase 4: Real Terraform integration tests using local fixtures

### Goal

Add integration tests that run real Terraform against local fixtures, skipped by default.

### TDD loop

- [ ] Write failing integration test scaffold that skips when not enabled.
- [ ] Add fixtures copied into temporary dirs before execution.
- [ ] Write integration test for valid basic Terraform.
- [ ] Write integration test for invalid HCL.
- [ ] Write integration test for missing provider or init/validate failure.
- [ ] Implement any missing fixture helpers.
- [ ] Run targeted integration tests with opt-in environment variable.
- [ ] Refactor after tests pass.
- [ ] Update this phase file's checklist and the status tracker in [00-PLAN.md](00-PLAN.md).
- [ ] Leave code committed or ready to commit.

### Tasks

- [ ] Add `testdata/fixtures/valid-basic`.
- [ ] Add `testdata/fixtures/invalid-hcl`.
- [ ] Add `testdata/fixtures/missing-provider`.
- [ ] Add `testdata/fixtures/plan-basic`.
- [ ] Add fixture copy helper.
- [ ] Add Terraform availability check helper.
- [ ] Add integration test skip behavior when `TERRAFORMER_RUN_INTEGRATION` is unset.
- [ ] Test `terraform fmt` on fixture.
- [ ] Test `terraform validate` on valid fixture.
- [ ] Test `terraform validate` on invalid fixture.
- [ ] Test `terraform plan` behavior where practical.
- [ ] Ensure integration tests do not require provider credentials.
- [ ] Ensure integration tests do not mutate committed fixtures.
- [ ] Run `TERRAFORMER_RUN_INTEGRATION=1 go test ./... -run Integration`.

### Completion criteria

- [ ] Integration tests are skipped by default.
- [ ] Integration tests pass when enabled in an environment with Terraform installed.
- [ ] Fixtures are local and credential-free.
- [ ] No integration test calls `apply` or `destroy`.
