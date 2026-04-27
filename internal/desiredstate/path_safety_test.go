package desiredstate

import (
	"testing"
)

func TestValidatePlanJSONPath(t *testing.T) {
	repoRoot := "/repo"

	tests := []struct {
		name     string
		planPath string
		wantErr  bool
	}{
		{
			name:     "valid relative path",
			planPath: ".terraformer/plan.json",
			wantErr:  false,
		},
		{
			name:     "valid nested path",
			planPath: "plans/my-plan.json",
			wantErr:  false,
		},
		{
			name:     "reject absolute path",
			planPath: "/etc/passwd",
			wantErr:  true,
		},
		{
			name:     "reject path traversal with ..",
			planPath: "../../../etc/passwd",
			wantErr:  true,
		},
		{
			name:     "reject path starting with -",
			planPath: "-plan.json",
			wantErr:  true,
		},
		{
			name:     "empty path",
			planPath: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePlanJSONPath(repoRoot, tt.planPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePlanJSONPath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
