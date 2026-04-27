// Package patch provides safe file modification operations within a repo root.
// All operations validate paths through the safety package to prevent traversal
// and symlink escapes. No command execution is performed during patching.
package patch

import (
	"fmt"

	"github.com/maazghani/terraformer/internal/safety"
)

// Service provides patch application operations bounded by the configured root.
type Service struct {
	root string
}

// New creates a Service rooted at root. root must be an absolute path to an
// existing directory; it is validated via safety.ValidateRepoRoot.
func New(root string) (*Service, error) {
	if err := safety.ValidateRepoRoot(root); err != nil {
		return nil, fmt.Errorf("patch.New: %w", err)
	}
	return &Service{root: root}, nil
}
