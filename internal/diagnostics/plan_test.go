package diagnostics

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestParsePlanSummary_GoldenOutput(t *testing.T) {
	// This tests the normalization of Terraform show -json plan output into
	// a stable PlanSummary structure. The golden file ensures the normalized
	// format remains stable.

	input := `{
  "format_version": "1.2",
  "terraform_version": "1.5.0",
  "resource_changes": [
    {
      "address": "local_file.example",
      "mode": "managed",
      "type": "local_file",
      "name": "example",
      "provider_name": "registry.terraform.io/hashicorp/local",
      "change": {
        "actions": ["create"],
        "before": null,
        "after": {
          "content": "hello"
        }
      }
    },
    {
      "address": "local_file.update",
      "mode": "managed",
      "type": "local_file",
      "name": "update",
      "provider_name": "registry.terraform.io/hashicorp/local",
      "change": {
        "actions": ["update"]
      }
    },
    {
      "address": "local_file.replace",
      "mode": "managed",
      "type": "local_file",
      "name": "replace",
      "provider_name": "registry.terraform.io/hashicorp/local",
      "change": {
        "actions": ["delete", "create"]
      }
    }
  ]
}`

	summary := ParsePlanSummary(input)

	// Compare against golden file
	goldenPath := filepath.Join("testdata", "golden", "plan_summary.json")
	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("failed to read golden file: %v", err)
	}

	got, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal plan summary: %v", err)
	}

	if string(got) != string(want) {
		t.Errorf("plan summary output mismatch:\ngot:\n%s\n\nwant:\n%s", got, want)
	}
}

func TestParsePlanSummary_MalformedJSON(t *testing.T) {
	// Test that malformed JSON returns an empty summary without panicking.
	malformedInput := "not json at all"

	summary := ParsePlanSummary(malformedInput)

	if summary.Create != 0 || summary.Update != 0 || summary.Delete != 0 || summary.Replace != 0 || summary.NoOp != 0 {
		t.Errorf("expected zero counts on malformed JSON, got %+v", summary)
	}
}
