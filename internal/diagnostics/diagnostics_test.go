package diagnostics

import (
	"testing"
)

func TestDiagnostic_Structure(t *testing.T) {
	// Test that Diagnostic type exists and has the expected fields.
	d := Diagnostic{
		Severity: "error",
		Summary:  "test summary",
		Detail:   "test detail",
		File:     "main.tf",
		Line:     10,
		Column:   5,
	}

	if d.Severity != "error" {
		t.Errorf("Severity: got %q, want %q", d.Severity, "error")
	}
	if d.Summary != "test summary" {
		t.Errorf("Summary: got %q, want %q", d.Summary, "test summary")
	}
	if d.Detail != "test detail" {
		t.Errorf("Detail: got %q, want %q", d.Detail, "test detail")
	}
	if d.File != "main.tf" {
		t.Errorf("File: got %q, want %q", d.File, "main.tf")
	}
	if d.Line != 10 {
		t.Errorf("Line: got %d, want %d", d.Line, 10)
	}
	if d.Column != 5 {
		t.Errorf("Column: got %d, want %d", d.Column, 5)
	}
}
