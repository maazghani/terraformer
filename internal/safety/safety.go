// Package safety provides path safety helpers and validation guards used
// throughout terraformer. All path resolution must go through this package.
package safety

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// forbiddenPrefixes is the set of path components that are never allowed to be
// accessed through ResolvePath.
var forbiddenPrefixes = []string{".git"}

// defaultExclusions is the set of path component prefixes excluded from
// directory listings and file operations by default.
var defaultExclusions = []string{".git", ".terraform"}

// ValidateRepoRoot checks that root is a valid, existing directory and is an
// absolute path. It returns a non-nil error if any condition is not met.
func ValidateRepoRoot(root string) error {
	if !filepath.IsAbs(root) {
		return fmt.Errorf("safety: repo root must be an absolute path, got %q", root)
	}
	fi, err := os.Stat(root)
	if err != nil {
		return fmt.Errorf("safety: repo root does not exist: %w", err)
	}
	if !fi.IsDir() {
		return fmt.Errorf("safety: repo root is not a directory: %q", root)
	}
	return nil
}

// ResolvePath resolves a user-supplied relative path against repoRoot and
// returns the cleaned absolute path. It rejects:
//   - absolute user paths
//   - paths that traverse outside repoRoot (../)
//   - paths into .git
//   - symlinks whose final resolved target is outside repoRoot
//
// repoRoot must be an absolute, cleaned path (as returned by ValidateRepoRoot).
func ResolvePath(repoRoot, userPath string) (string, error) {
	// Reject absolute user paths immediately.
	if filepath.IsAbs(userPath) {
		return "", fmt.Errorf("safety: absolute paths are not allowed: %q", userPath)
	}

	// Build the candidate path and clean it.
	candidate := filepath.Clean(filepath.Join(repoRoot, userPath))

	// After cleaning, the path must still be inside repoRoot.
	if !isUnderRoot(repoRoot, candidate) {
		return "", fmt.Errorf("safety: path %q escapes repo root %q", userPath, repoRoot)
	}

	// Reject forbidden prefixes (e.g. .git).
	rel, err := filepath.Rel(repoRoot, candidate)
	if err != nil {
		return "", fmt.Errorf("safety: cannot relativize path: %w", err)
	}
	if isForbidden(rel) {
		return "", fmt.Errorf("safety: path %q is forbidden", userPath)
	}

	// Resolve symlinks to detect escapes.
	resolved, err := filepath.EvalSymlinks(candidate)
	if err != nil {
		// The target may not yet exist (e.g., new file being created).
		// In that case we cannot resolve symlinks and we conservatively allow it
		// because there is no real symlink to follow.
		if os.IsNotExist(err) {
			return candidate, nil
		}
		return "", fmt.Errorf("safety: cannot resolve path: %w", err)
	}

	// After symlink resolution, the path must still be inside repoRoot.
	// Resolve the root too, in case it is itself a symlink.
	resolvedRoot, err := filepath.EvalSymlinks(repoRoot)
	if err != nil {
		resolvedRoot = repoRoot
	}
	if !isUnderRoot(resolvedRoot, resolved) {
		return "", fmt.Errorf("safety: symlink %q escapes repo root", userPath)
	}

	// Return the cleaned candidate (not the symlink-resolved path) so that
	// callers receive paths in terms of the repo root they supplied.
	return candidate, nil
}

// IsExcludedByDefault reports whether a relative path should be omitted from
// directory listings by default. It returns true for .git and .terraform and
// anything nested under them.
func IsExcludedByDefault(relPath string) bool {
	cleaned := filepath.Clean(relPath)
	parts := strings.Split(cleaned, string(filepath.Separator))
	if len(parts) == 0 {
		return false
	}
	first := parts[0]
	for _, ex := range defaultExclusions {
		if first == ex {
			return true
		}
	}
	return false
}

// isUnderRoot reports whether p is equal to root or nested under it.
// Both p and root must be cleaned absolute paths.
func isUnderRoot(root, p string) bool {
	// Ensure root ends with separator for prefix matching.
	if root == p {
		return true
	}
	return strings.HasPrefix(p, root+string(filepath.Separator))
}

// isForbidden reports whether a cleaned relative path starts with any
// forbidden prefix component.
func isForbidden(rel string) bool {
	parts := strings.Split(rel, string(filepath.Separator))
	if len(parts) == 0 {
		return false
	}
	first := parts[0]
	for _, f := range forbiddenPrefixes {
		if first == f {
			return true
		}
	}
	return false
}
