package terraform

import (
	"testing"

	"github.com/maazghani/terraformer/internal/runner"
)

const validateJSONErrors = `{
  "valid": false,
  "error_count": 1,
  "warning_count": 0,
  "diagnostics": [
    {
      "severity": "error",
      "summary": "Invalid reference",
      "detail": "A reference to a resource type must be followed by at least one attribute access.",
      "range": {
        "filename": "main.tf",
        "start": {"line": 12, "column": 5, "byte": 0},
        "end": {"line": 12, "column": 10, "byte": 0}
      }
    }
  ]
}`

func TestService_Validate_ParsesJSONDiagnostics(t *testing.T) {
	fake := runner.NewFakeRunner()
	fake.Register("terraform", runner.Result{ExitCode: 1, Stdout: validateJSONErrors}, nil)
	svc := NewService(fake, "/repo")

	resp := svc.Validate(ValidateRequest{JSON: true})

	if resp.OK {
		t.Fatalf("expected OK=false on validation errors")
	}
	if len(resp.Diagnostics) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(resp.Diagnostics))
	}
	d := resp.Diagnostics[0]
	if d.Severity != "error" {
		t.Errorf("severity: got %q want error", d.Severity)
	}
	if d.Summary != "Invalid reference" {
		t.Errorf("summary: got %q", d.Summary)
	}
	if d.File != "main.tf" {
		t.Errorf("file: got %q", d.File)
	}
	if d.Line != 12 {
		t.Errorf("line: got %d want 12", d.Line)
	}
	if d.Column != 5 {
		t.Errorf("column: got %d want 5", d.Column)
	}
}

func TestService_Validate_NonJSONStdoutDoesNotPanic(t *testing.T) {
	fake := runner.NewFakeRunner()
	fake.Register("terraform", runner.Result{ExitCode: 1, Stdout: "not json", Stderr: "boom"}, nil)
	svc := NewService(fake, "/repo")

	resp := svc.Validate(ValidateRequest{JSON: true})

	if resp.OK {
		t.Fatalf("expected OK=false")
	}
	if len(resp.Diagnostics) != 0 {
		t.Errorf("expected no parsed diagnostics on malformed JSON, got %d", len(resp.Diagnostics))
	}
}

func TestService_Validate_NoJSONFlagSkipsParsing(t *testing.T) {
	fake := runner.NewFakeRunner()
	fake.Register("terraform", runner.Result{ExitCode: 0, Stdout: validateJSONErrors}, nil)
	svc := NewService(fake, "/repo")

	resp := svc.Validate(ValidateRequest{JSON: false})

	if len(resp.Diagnostics) != 0 {
		t.Errorf("expected no parsed diagnostics when JSON not requested, got %d", len(resp.Diagnostics))
	}
}
