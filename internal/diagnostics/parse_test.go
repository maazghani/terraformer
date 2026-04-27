package diagnostics

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestParseValidateJSON_GoldenOutput(t *testing.T) {
	// This tests the normalization of Terraform validate -json output into
	// stable Diagnostic structures. The golden file ensures the normalized
	// format remains stable and model-friendly.

	input := `{
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

	diags := ParseValidateJSON(input, "")

	if len(diags) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(diags))
	}

	// Compare against golden file
	goldenPath := filepath.Join("testdata", "golden", "validate_diagnostics.json")
	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("failed to read golden file: %v", err)
	}

	got, err := json.MarshalIndent(diags, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal diagnostics: %v", err)
	}

	if string(got) != string(want) {
		t.Errorf("diagnostics output mismatch:\ngot:\n%s\n\nwant:\n%s", got, want)
	}
}

func TestParseFallbackDiagnostic_GoldenOutput(t *testing.T) {
	// Test that malformed JSON produces a safe fallback diagnostic.
	malformedInput := "not json at all"
	stderr := "Error: Invalid configuration"

	diags := ParseValidateJSON(malformedInput, stderr)

	if len(diags) == 0 {
		t.Fatal("expected at least one fallback diagnostic")
	}

	// Compare against golden file
	goldenPath := filepath.Join("testdata", "golden", "validate_fallback_diagnostic.json")
	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("failed to read golden file: %v", err)
	}

	got, err := json.MarshalIndent(diags, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal diagnostics: %v", err)
	}

	if string(got) != string(want) {
		t.Errorf("fallback diagnostic output mismatch:\ngot:\n%s\n\nwant:\n%s", got, want)
	}
}
