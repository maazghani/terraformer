> **Phase 7 of 9** | [← Phase 6](06-phase6.md) | [Index: 00-PLAN.md](00-PLAN.md) | Next: [Phase 8 →](08-phase8.md)

# Phase 7: Desired-state spec and plan comparison

Introduce the `check_desired_state` tool contract and minimal implementation. v0 honesty: report `not_implemented` or stubbed behavior; never declare success based on plan alone. Comparison logic can be expanded behind tests in v0.1+.

## Phase 7: Desired-state spec and plan comparison

### Goal

Introduce a minimal desired-state contract and comparison flow. Do not overclaim.

### TDD loop

- [x] Write failing tests for minimal desired-state request validation.
- [x] Write failing tests for stubbed `not_implemented` or minimal comparison behavior.
- [x] Write failing tests for unsafe plan JSON paths.
- [x] Write failing tests for mismatched plan summary behavior if implemented.
- [x] Implement smallest desired-state code.
- [x] Run targeted tests.
- [x] Refactor after tests pass.
- [x] Update this phase file's checklist and the status tracker in [00-PLAN.md](00-PLAN.md).
- [x] Leave code committed or ready to commit.

### Initial desired-state model

```json
{
  "resources": [
    {
      "address": "local_file.example",
      "actions": ["create"]
    }
  ],
  "forbidden_actions": ["delete"]
}
```

Initial comparison result
```json
{
  "ok": true,
  "status": "mismatched",
  "matched": false,
  "mismatches": [
    {
      "address": "local_file.example",
      "reason": "Expected create but plan contains delete."
    }
  ],
  "warnings": []
}
```

### Tasks

- [x] Define desired-state schema.
- [x] Define comparison result schema.
- [x] Implement request validation.
- [x] Implement stub response if full comparison is not ready.
- [x] Implement basic comparison for resource address and actions if practical.
- [x] Support forbidden actions such as delete.
- [x] Ensure terraform_plan response remains separate from desired-state result.
- [x] Ensure check_desired_state is the only tool that declares desired-state status.
- [x] Add golden tests for matched, mismatched, and not implemented responses.
- [x] Run go test ./internal/desiredstate ./internal/terraform.

### Completion criteria

- [x] Desired-state tool exists internally.
- [x] It does not pretend to do more than implemented.
- [x] Plan success alone cannot produce desired-state success.
- [x] Tests prove matched and mismatched behavior if comparison is implemented.