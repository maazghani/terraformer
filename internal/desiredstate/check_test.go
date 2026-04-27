package desiredstate

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheck_StubBehavior(t *testing.T) {
	// Create a temporary directory to act as repo root.
	repoRoot := t.TempDir()

	// Create a dummy plan JSON file.
	planPath := "plan.json"
	fullPlanPath := filepath.Join(repoRoot, planPath)
	planJSON := `{"format_version":"1.0","terraform_version":"1.0.0"}`
	if err := os.WriteFile(fullPlanPath, []byte(planJSON), 0644); err != nil {
		t.Fatalf("failed to create plan JSON: %v", err)
	}

	ds := DesiredState{
		Resources: []ResourceExpectation{
			{
				Address: "local_file.example",
				Actions: []string{"create"},
			},
		},
	}

	result, err := Check(repoRoot, planPath, ds)
	if err != nil {
		t.Fatalf("Check() returned unexpected error: %v", err)
	}

	if !result.OK {
		t.Errorf("Check() result.OK = false, want true")
	}

	if result.Status != "not_implemented" {
		t.Errorf("Check() result.Status = %q, want %q", result.Status, "not_implemented")
	}

	if result.Matched {
		t.Errorf("Check() result.Matched = true, want false for not_implemented")
	}

	if len(result.Warnings) == 0 {
		t.Errorf("Check() result.Warnings is empty, expected warning about stub behavior")
	}
}

func TestCheck_InvalidPath(t *testing.T) {
	repoRoot := t.TempDir()

	ds := DesiredState{
		Resources: []ResourceExpectation{},
	}

	// Try with an unsafe path.
	_, err := Check(repoRoot, "../../../etc/passwd", ds)
	if err == nil {
		t.Errorf("Check() with unsafe path did not return error")
	}
}

func TestCheck_EmptyDesiredState(t *testing.T) {
	repoRoot := t.TempDir()

	// Create a dummy plan JSON file.
	planPath := "plan.json"
	fullPlanPath := filepath.Join(repoRoot, planPath)
	planJSON := `{"format_version":"1.0"}`
	if err := os.WriteFile(fullPlanPath, []byte(planJSON), 0644); err != nil {
		t.Fatalf("failed to create plan JSON: %v", err)
	}

	ds := DesiredState{
		Resources: []ResourceExpectation{},
	}

	result, err := Check(repoRoot, planPath, ds)
	if err != nil {
		t.Fatalf("Check() returned unexpected error: %v", err)
	}

	if !result.OK {
		t.Errorf("Check() result.OK = false, want true")
	}

	// Empty desired state should still return not_implemented stub.
	if result.Status != "not_implemented" {
		t.Errorf("Check() result.Status = %q, want %q", result.Status, "not_implemented")
	}
}

func TestCheck_InvalidDesiredState(t *testing.T) {
	repoRoot := t.TempDir()

	// Create a dummy plan JSON file.
	planPath := "plan.json"
	fullPlanPath := filepath.Join(repoRoot, planPath)
	planJSON := `{"format_version":"1.0"}`
	if err := os.WriteFile(fullPlanPath, []byte(planJSON), 0644); err != nil {
		t.Fatalf("failed to create plan JSON: %v", err)
	}

	// Invalid desired state: resource without address.
	ds := DesiredState{
		Resources: []ResourceExpectation{
			{
				Address: "",
				Actions: []string{"create"},
			},
		},
	}

	_, err := Check(repoRoot, planPath, ds)
	if err == nil {
		t.Errorf("Check() with invalid desired state did not return error")
	}
}
