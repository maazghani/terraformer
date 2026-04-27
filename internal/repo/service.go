// Package repo provides safe, repo-local file access operations.
// All path access goes through the safety package to prevent traversal
// and symlink escapes.
package repo

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/maazghani/terraformer/internal/safety"
)

const defaultMaxFiles = 200

// Service provides repo-local file operations bounded by the configured root.
type Service struct {
	root string
}

// New creates a Service rooted at root. root must be an absolute path to an
// existing directory; it is validated via safety.ValidateRepoRoot.
func New(root string) (*Service, error) {
	if err := safety.ValidateRepoRoot(root); err != nil {
		return nil, fmt.Errorf("repo.New: %w", err)
	}
	return &Service{root: root}, nil
}

// ListFilesRequest describes the parameters for a file listing operation.
type ListFilesRequest struct {
	// Path is a repo-relative directory to list. "." means repo root.
	Path string
	// MaxFiles is the maximum number of file entries to return. When zero,
	// defaultMaxFiles is used. When the limit is exceeded, Truncated is set.
	MaxFiles int
}

// FileEntry represents a single entry in a file listing.
type FileEntry struct {
	// Path is the repo-relative, slash-separated path to the entry.
	Path string
	// SizeBytes is the size of the file in bytes (0 for directories).
	SizeBytes int64
	// Kind is "file" or "dir".
	Kind string
}

// ListFilesResponse is returned by ListFiles.
type ListFilesResponse struct {
	// Files is the list of file entries under the requested path.
	Files []FileEntry
	// Truncated is true when the result was capped by MaxFiles.
	Truncated bool
	// Warnings contains non-fatal advisories.
	Warnings []string
}

// ListFiles walks the repo-relative Path and returns file entries, applying
// default exclusions (.git, .terraform) and the MaxFiles cap.
func (s *Service) ListFiles(req ListFilesRequest) (ListFilesResponse, error) {
	path := req.Path
	if path == "" {
		path = "."
	}

	// Resolve the requested path safely.
	// "." resolves to the repo root itself, which is valid.
	var startDir string
	if path == "." {
		startDir = s.root
	} else {
		resolved, err := safety.ResolvePath(s.root, path)
		if err != nil {
			return ListFilesResponse{}, fmt.Errorf("repo.ListFiles: %w", err)
		}
		startDir = resolved
	}

	maxFiles := req.MaxFiles
	if maxFiles <= 0 {
		maxFiles = defaultMaxFiles
	}

	var entries []FileEntry
	truncated := false

	err := filepath.WalkDir(startDir, func(absPath string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			// Skip entries we cannot access.
			return nil
		}

		// Compute the relative path from the repo root for exclusion checks.
		relFromRoot, err := filepath.Rel(s.root, absPath)
		if err != nil {
			return nil
		}

		// Skip the root directory itself.
		if relFromRoot == "." {
			return nil
		}

		// Apply default exclusions to any entry under an excluded directory.
		if safety.IsExcludedByDefault(relFromRoot) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// For symlinks pointing outside the repo, skip safely.
		if d.Type()&fs.ModeSymlink != 0 {
			// Try to resolve and validate; skip if it escapes.
			if _, err := safety.ResolvePath(s.root, relFromRoot); err != nil {
				return nil
			}
		}

		// Only record regular files (not directories).
		if d.IsDir() {
			return nil
		}

		// Enforce max_files limit.
		if len(entries) >= maxFiles {
			truncated = true
			return filepath.SkipAll
		}

		info, err := d.Info()
		if err != nil {
			return nil
		}

		// Normalize path to use forward slashes.
		normalized := filepath.ToSlash(relFromRoot)

		entries = append(entries, FileEntry{
			Path:      normalized,
			SizeBytes: info.Size(),
			Kind:      "file",
		})

		return nil
	})
	if err != nil {
		return ListFilesResponse{}, fmt.Errorf("repo.ListFiles: walk error: %w", err)
	}

	if entries == nil {
		entries = []FileEntry{}
	}

	return ListFilesResponse{
		Files:     entries,
		Truncated: truncated,
		Warnings:  []string{},
	}, nil
}
