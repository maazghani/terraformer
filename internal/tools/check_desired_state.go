package tools

import (
	"github.com/maazghani/terraformer/internal/desiredstate"
)

// CheckDesiredStateRequest is the request struct for the check_desired_state tool.
type CheckDesiredStateRequest struct {
	// DesiredState is the desired-state specification.
	DesiredState desiredstate.DesiredState `json:"desired_state"`
	// PlanJSONPath is the repo-relative path to a Terraform plan JSON file.
	PlanJSONPath string `json:"plan_json_path"`
}

// CheckDesiredStateResponse is the response for the check_desired_state tool.
type CheckDesiredStateResponse struct {
	// OK is true when the check completed without error.
	OK bool `json:"ok"`
	// Status is one of: "not_implemented", "not_checked", "matched", "mismatched".
	Status string `json:"status"`
	// Matched is true when the plan matches the desired state.
	Matched bool `json:"matched"`
	// Mismatches contains discrepancies between desired state and plan.
	Mismatches []desiredstate.Mismatch `json:"mismatches"`
	// Warnings contains non-fatal advisories.
	Warnings []string `json:"warnings"`
}

// CheckDesiredState executes the check_desired_state tool.
// repoRoot is the absolute path to the repository root.
func CheckDesiredState(repoRoot string, req CheckDesiredStateRequest) CheckDesiredStateResponse {
	result, err := desiredstate.Check(repoRoot, req.PlanJSONPath, req.DesiredState)
	if err != nil {
		return CheckDesiredStateResponse{
			OK:         false,
			Status:     "error",
			Matched:    false,
			Mismatches: []desiredstate.Mismatch{},
			Warnings:   []string{err.Error()},
		}
	}

	return CheckDesiredStateResponse{
		OK:         result.OK,
		Status:     result.Status,
		Matched:    result.Matched,
		Mismatches: result.Mismatches,
		Warnings:   result.Warnings,
	}
}
