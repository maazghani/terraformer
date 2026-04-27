package desiredstate

// Check compares a Terraform plan against the desired state expectations.
// In v0, this returns a stubbed "not_implemented" response to establish the
// contract without overclaiming functionality. Future versions can expand
// this to perform actual comparison logic.
func Check(repoRoot, planJSONPath string, desired DesiredState) (ComparisonResult, error) {
	// Validate the desired state.
	if err := desired.Validate(); err != nil {
		return ComparisonResult{}, err
	}

	// Validate the plan JSON path is safe.
	if err := ValidatePlanJSONPath(repoRoot, planJSONPath); err != nil {
		return ComparisonResult{}, err
	}

	// For v0, return a stubbed not_implemented response.
	// This establishes the contract without overclaiming functionality.
	result := ComparisonResult{
		OK:         true,
		Status:     "not_implemented",
		Matched:    false,
		Mismatches: []Mismatch{},
		Warnings: []string{
			"Desired-state comparison is stubbed in this version.",
		},
	}

	return result, nil
}
