package desiredstate

import (
	"fmt"
	"strings"

	"github.com/maazghani/terraformer/internal/safety"
)

// ValidatePlanJSONPath validates that a plan JSON path is safe to access.
// It rejects:
//   - empty paths
//   - paths starting with "-" (option injection)
//   - absolute paths
//   - paths that escape the repo root
func ValidatePlanJSONPath(repoRoot, planPath string) error {
	if planPath == "" {
		return fmt.Errorf("plan JSON path cannot be empty")
	}

	// Reject paths starting with "-" to prevent option injection.
	if strings.HasPrefix(planPath, "-") {
		return fmt.Errorf("plan JSON path cannot start with '-': %q", planPath)
	}

	// Use safety.ResolvePath to validate the path is within repo root.
	_, err := safety.ResolvePath(repoRoot, planPath)
	if err != nil {
		return fmt.Errorf("invalid plan JSON path: %w", err)
	}

	return nil
}
