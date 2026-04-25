# Specification: HTTP/JSON Tool Contracts

This document defines the structured request and response contracts for each v0 tool, including per-tool TDD requirements. All tool implementations must satisfy these contracts.

## HTTP/JSON tool contracts for v0
All tool responses must use structured JSON-compatible objects.
{
  "ok": true,
  "command": {
    "name": "terraform",
    "args": ["validate"],
    "working_dir": "."
  },
  "stdout": "",
  "stderr": "",
  "exit_code": 0,
  "duration_ms": 123,
  "diagnostics": [],
  "warnings": []
}
On failure:
{
  "ok": false,
  "command": {
    "name": "terraform",
    "args": ["validate"],
    "working_dir": "."
  },
  "stdout": "",
  "stderr": "",
  "exit_code": 1,
  "duration_ms": 123,
  "diagnostics": [
    {
      "severity": "error",
      "summary": "Invalid reference",
      "detail": "A reference to a resource type must be followed by at least one attribute access.",
      "file": "main.tf",
      "line": 12,
      "column": 5
    }
  ],
  "warnings": []
}
Secret values must be redacted before response serialization (v0.1).

### Tool: terraform_init
Purpose:
Run terraform init inside the configured repo root.
Request:
{
  "upgrade": false,
  "backend": true
}
Allowed behavior:
 May run terraform init.
 May include safe flags such as -input=false.
 Must run inside repo root.
 Must not expose credentials.
 Must not accept arbitrary extra args in v0.
Response:
{
  "ok": true,
  "command": {
    "name": "terraform",
    "args": ["init", "-input=false"],
    "working_dir": "."
  },
  "stdout": "",
  "stderr": "",
  "exit_code": 0,
  "duration_ms": 0,
  "diagnostics": [],
  "warnings": []
}
TDD requirements:
 Test fake runner receives exactly the expected args.
 Test arbitrary args are rejected or impossible to pass.
 Test command runs at repo root.
 Test stdout and stderr redaction.

### Tool: terraform_fmt
Purpose:
Run terraform fmt.
Request:
{
  "check": false,
  "recursive": true
}
Allowed behavior:
 May run terraform fmt.
 May support -check.
 May support -recursive.
 Must not accept arbitrary paths outside repo root.
 Must not accept arbitrary extra args.
Response:
{
  "ok": true,
  "command": {
    "name": "terraform",
    "args": ["fmt", "-recursive"],
    "working_dir": "."
  },
  "stdout": "",
  "stderr": "",
  "exit_code": 0,
  "duration_ms": 0,
  "diagnostics": [],
  "warnings": []
}
TDD requirements:
 Test default args.
 Test check=true.
 Test recursive=true.
 Test output structure.
 Test fake runner failure maps to ok=false.

### Tool: terraform_validate
Purpose:
Run terraform validate.
Request:
{
  "json": true
}
Allowed behavior:
 Prefer terraform validate -json for structured diagnostics.
 Must run inside repo root.
 Must normalize diagnostics where possible.
 Must not accept arbitrary extra args.
Response:
{
  "ok": false,
  "command": {
    "name": "terraform",
    "args": ["validate", "-json"],
    "working_dir": "."
  },
  "stdout": "",
  "stderr": "",
  "exit_code": 1,
  "duration_ms": 0,
  "diagnostics": [
    {
      "severity": "error",
      "summary": "",
      "detail": "",
      "file": "",
      "line": 0,
      "column": 0
    }
  ],
  "warnings": []
}
TDD requirements:
 Test validate command args.
 Test JSON diagnostics parsing.
 Test non-JSON fallback diagnostics.
 Test redaction.
 Test failed validation does not panic.

### Tool: terraform_plan
Purpose:
Run terraform plan.
Request:
{
  "out": ".terraformer/plan.tfplan",
  "detailed_exitcode": true,
  "refresh": false
}
Allowed behavior:
 May run terraform plan.
 Should use -input=false.
 Should support -detailed-exitcode.
 Should support a safe repo-local -out path.
 Must reject plan output paths outside the repo.
 Must normalize exit codes:
0: success with no changes
2: success with changes when -detailed-exitcode is used
1: failure
 Must not treat plan success as desired-state success.
Response:
{
  "ok": true,
  "plan_status": "changes_present",
  "desired_state_status": "not_checked",
  "command": {
    "name": "terraform",
    "args": ["plan", "-input=false", "-detailed-exitcode", "-out=.terraformer/plan.tfplan"],
    "working_dir": "."
  },
  "stdout": "",
  "stderr": "",
  "exit_code": 2,
  "duration_ms": 0,
  "diagnostics": [],
  "warnings": []
}
TDD requirements:
 Test exit code 0 maps to no_changes.
 Test exit code 2 maps to changes_present.
 Test exit code 1 maps to failure.
 Test unsafe out path is rejected.
 Test plan does not call desired-state success automatically.
 Test fake runner args exactly.

### Tool: terraform_show_json
Purpose:
Run terraform show -json against a repo-local plan file.
Request:
{
  "plan_path": ".terraformer/plan.tfplan"
}
Allowed behavior:
 Must run terraform show -json <plan_path>.
 Must reject paths outside the repo.
 Must parse JSON when possible.
 Must return raw JSON only if redacted and size-limited.
 Should return a normalized summary.
Response:
{
  "ok": true,
  "command": {
    "name": "terraform",
    "args": ["show", "-json", ".terraformer/plan.tfplan"],
    "working_dir": "."
  },
  "stdout": "{...}",
  "stderr": "",
  "exit_code": 0,
  "duration_ms": 0,
  "plan_summary": {
    "create": 0,
    "update": 0,
    "delete": 0,
    "replace": 0,
    "no_op": 0
  },
  "diagnostics": [],
  "warnings": []
}
TDD requirements:
 Test path safety.
 Test command args.
 Test plan summary parsing.
 Test malformed JSON handling.
 Test redaction and response size limits.

### Tool: list_repo_files
Purpose:
List repo-local files visible to the agent.
Request:
{
  "path": ".",
  "include_globs": ["*.tf", "*.tfvars", "*.md"],
  "exclude_globs": [".terraform/**", ".git/**"],
  "max_files": 200
}
Allowed behavior:
 Must list only inside repo root.
 Must exclude .git.
 Must exclude .terraform by default.
 Must exclude likely secret files by default.
 Must return normalized relative paths.
 Must enforce max file count.
Response:
{
  "ok": true,
  "files": [
    {
      "path": "main.tf",
      "size_bytes": 1234,
      "kind": "file"
    }
  ],
  "truncated": false,
  "warnings": []
}
TDD requirements:
 Test basic listing.
 Test exclude defaults.
 Test max file limit.
 Test path traversal rejection.
 Test symlink escape rejection.

### Tool: read_repo_file
Purpose:
Read a repo-local file.
Request:
{
  "path": "main.tf",
  "max_bytes": 65536
}
Allowed behavior:
 Must read only files inside repo root.
 Must reject directories.
 Must reject unsafe paths.
 Must enforce max bytes.
 Must redact secrets.
 Must return content and metadata.
Response:
{
  "ok": true,
  "path": "main.tf",
  "content": "",
  "size_bytes": 0,
  "truncated": false,
  "warnings": []
}
TDD requirements:
 Test reading a valid file.
 Test missing file behavior.
 Test directory rejection.
 Test max bytes truncation.
 Test secret redaction.
 Test traversal and symlink escape rejection.

### Tool: apply_patch
Purpose:
Apply a structured JSON patch to repo-local files.
Request:
{
  "files": [
    {
      "path": "main.tf",
      "operation": "write",
      "content": "resource \"local_file\" \"example\" {}\n"
    },
    {
      "path": "vars.tf",
      "operation": "write",
      "content": "variable \"env\" { default = \"dev\" }\n"
    }
  ]
}
Request field descriptions:
- `path`: repo-relative file path
- `operation`: "write" to create/overwrite, "delete" to remove (delete only if explicitly tested)
- `content`: file content (required for write, omitted for delete)
Allowed behavior:
 Must apply only repo-local file changes.
 Must reject patch paths outside the repo.
 Must reject symlink escapes.
 Must return changed files.
 Must fail cleanly on invalid patches.
 Must not run Terraform automatically.
Response:
{
  "ok": true,
  "changed_files": ["main.tf", "vars.tf"],
  "rejected_files": [],
  "warnings": []
}
TDD requirements:
 Test valid write operations.
 Test invalid file paths.
 Test patch traversal rejection.
 Test symlink escape rejection.
 Test patch does not execute commands.
 Test changed file reporting.

### Tool: check_desired_state
Purpose:
Compare the current Terraform plan or plan JSON against a desired-state specification.
Initial v0 behavior may be minimal or stubbed, but it must be behind tests.
Request:
{
  "desired_state": {
    "resources": []
  },
  "plan_json_path": ".terraformer/plan.json"
}
Allowed behavior:
 Must not claim full desired-state matching before implemented.
 Must return not_implemented, not_checked, matched, or mismatched.
 Must be designed so future comparison logic can be added without changing the tool contract unnecessarily.
 Must parse enough plan JSON to support future checks.
 Must reject unsafe plan JSON paths.
Response:
{
  "ok": true,
  "status": "not_implemented",
  "matched": false,
  "mismatches": [],
  "warnings": [
    "Desired-state comparison is stubbed in this version."
  ]
}
TDD requirements:
 Test stub behavior.
 Test unsafe path rejection.
 Test empty desired-state behavior.
 Test malformed desired-state behavior.
 Test future-compatible response shape.

## See also

- [00-spec.md](00-spec.md) — Safety rules that every tool must uphold
- [01-testing.md](01-testing.md) — Test categories and TDD loop that apply to every tool implementation
- [03-development.md](03-development.md) — Commands for running the test suite
