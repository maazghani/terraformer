package tools

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/maazghani/terraformer/internal/desiredstate"
)

func TestCheckDesiredState_Success(t *testing.T) {
	repoRoot := t.TempDir()

	// Create a dummy plan JSON file.
	planPath := "plan.json"
	fullPlanPath := filepath.Join(repoRoot, planPath)
	planJSON := `{"format_version":"1.0","terraform_version":"1.0.0"}`
	if err := os.WriteFile(fullPlanPath, []byte(planJSON), 0644); err != nil {
		t.Fatalf("failed to create plan JSON: %v", err)
	}

	req := CheckDesiredStateRequest{
		DesiredState: desiredstate.DesiredState{
			Resources: []desiredstate.ResourceExpectation{
				{
					Address: "local_file.example",
					Actions: []string{"create"},
				},
			},
		},
		PlanJSONPath: planPath,
	}

	resp := CheckDesiredState(repoRoot, req)

	if !resp.OK {
		t.Errorf("CheckDesiredState() OK = false, want true")
	}

	if resp.Status != "not_implemented" {
		t.Errorf("CheckDesiredState() Status = %q, want %q", resp.Status, "not_implemented")
	}

	if resp.Matched {
		t.Errorf("CheckDesiredState() Matched = true, want false for not_implemented")
	}
}

func TestCheckDesiredState_InvalidPath(t *testing.T) {
	repoRoot := t.TempDir()

	req := CheckDesiredStateRequest{
		DesiredState: desiredstate.DesiredState{
			Resources: []desiredstate.ResourceExpectation{},
		},
		PlanJSONPath: "../../../etc/passwd",
	}

	resp := CheckDesiredState(repoRoot, req)

	if resp.OK {
		t.Errorf("CheckDesiredState() OK = true, want false for invalid path")
	}

	if len(resp.Warnings) == 0 {
		t.Errorf("CheckDesiredState() Warnings is empty, expected error message")
	}
}

func TestCheckDesiredState_InvalidDesiredState(t *testing.T) {
	repoRoot := t.TempDir()

	// Create a dummy plan JSON file.
	planPath := "plan.json"
	fullPlanPath := filepath.Join(repoRoot, planPath)
	planJSON := `{"format_version":"1.0"}`
	if err := os.WriteFile(fullPlanPath, []byte(planJSON), 0644); err != nil {
		t.Fatalf("failed to create plan JSON: %v", err)
	}

	// Invalid: resource without address.
	req := CheckDesiredStateRequest{
		DesiredState: desiredstate.DesiredState{
			Resources: []desiredstate.ResourceExpectation{
				{
					Address: "",
					Actions: []string{"create"},
				},
			},
		},
		PlanJSONPath: planPath,
	}

	resp := CheckDesiredState(repoRoot, req)

	if resp.OK {
		t.Errorf("CheckDesiredState() OK = true, want false for invalid desired state")
	}

	if len(resp.Warnings) == 0 {
		t.Errorf("CheckDesiredState() Warnings is empty, expected error message")
	}
}
