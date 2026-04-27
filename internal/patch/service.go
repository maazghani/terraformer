// Package patch provides safe file modification operations within a repo root.
// All operations validate paths through the safety package to prevent traversal
// and symlink escapes. No command execution is performed during patching.
package patch

import (
	"fmt"
	"os"
	"path/filepath"

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

// FileOperation represents a single file modification operation within a patch.
type FileOperation struct {
	// Path is the repo-relative path to the file.
	Path string
	// Operation is the type of operation: "write" or "delete".
	Operation string
	// Content is the file content for "write" operations (omitted for "delete").
	Content string
}

// ApplyPatchRequest describes a batch of file operations to apply atomically.
type ApplyPatchRequest struct {
	// Files is the list of file operations to perform.
	Files []FileOperation
}

// ApplyPatchResponse is returned by ApplyPatch.
type ApplyPatchResponse struct {
	// OK is true when all operations succeeded.
	OK bool
	// ChangedFiles lists the paths of files that were successfully modified.
	ChangedFiles []string
	// RejectedFiles lists the paths of files that were rejected due to safety.
	RejectedFiles []string
	// Warnings contains non-fatal advisories.
	Warnings []string
}

// ApplyPatch applies the requested file operations to the repo root. Each
// operation's path is validated through safety.ResolvePath before any file
// system changes are made. If any path is unsafe, it is added to RejectedFiles
// and not applied. OK is true only when all operations succeed.
func (s *Service) ApplyPatch(req ApplyPatchRequest) (ApplyPatchResponse, error) {
	var changedFiles []string
	var rejectedFiles []string
	var warnings []string

	// Validate and apply each operation.
	for _, op := range req.Files {
		// Validate the path through safety.
		absPath, err := safety.ResolvePath(s.root, op.Path)
		if err != nil {
			// Path is unsafe; reject it.
			rejectedFiles = append(rejectedFiles, op.Path)
			continue
		}

		switch op.Operation {
		case "write":
			// Ensure parent directory exists and is safe.
			parentDir := filepath.Dir(absPath)

			// Validate parent directory is also safe (in case it's a symlink).
			relParent := filepath.Dir(op.Path)
			if relParent != "." {
				if _, err := safety.ResolvePath(s.root, relParent); err != nil {
					rejectedFiles = append(rejectedFiles, op.Path)
					continue
				}
			}

			if err := os.MkdirAll(parentDir, 0755); err != nil {
				rejectedFiles = append(rejectedFiles, op.Path)
				continue
			}

			// Write the file.
			if err := os.WriteFile(absPath, []byte(op.Content), 0644); err != nil {
				rejectedFiles = append(rejectedFiles, op.Path)
				continue
			}
			changedFiles = append(changedFiles, op.Path)

		case "delete":
			// Delete the file if it exists.
			if err := os.Remove(absPath); err != nil {
				if !os.IsNotExist(err) {
					rejectedFiles = append(rejectedFiles, op.Path)
					continue
				}
				// If file doesn't exist, we treat it as successful (idempotent).
			}
			changedFiles = append(changedFiles, op.Path)

		default:
			// Unknown operation; reject it.
			rejectedFiles = append(rejectedFiles, op.Path)
		}
	}

	// Ensure empty slices are represented as empty arrays in JSON.
	if changedFiles == nil {
		changedFiles = []string{}
	}
	if rejectedFiles == nil {
		rejectedFiles = []string{}
	}
	if warnings == nil {
		warnings = []string{}
	}

	// OK is true only when all operations succeeded and no files were rejected.
	ok := len(rejectedFiles) == 0 && len(req.Files) == len(changedFiles)

	return ApplyPatchResponse{
		OK:            ok,
		ChangedFiles:  changedFiles,
		RejectedFiles: rejectedFiles,
		Warnings:      warnings,
	}, nil
}
