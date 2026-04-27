package safety_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/maazghani/terraformer/internal/safety"
)

// --------------------------------------------------------------------------
// Repo root validation helpers
// --------------------------------------------------------------------------

// TestRepoRootValidation_Exists verifies that a directory is accepted.
func TestRepoRootValidation_Exists(t *testing.T) {
	dir := t.TempDir()
	if err := safety.ValidateRepoRoot(dir); err != nil {
		t.Errorf("expected valid dir to be accepted, got error: %v", err)
	}
}

// TestRepoRootValidation_MustBeAbsolute verifies that a relative path is rejected.
func TestRepoRootValidation_MustBeAbsolute(t *testing.T) {
	if err := safety.ValidateRepoRoot("relative/path"); err == nil {
		t.Error("expected relative path to be rejected, got nil error")
	}
}

// TestRepoRootValidation_MustExist verifies that a non-existent path is rejected.
func TestRepoRootValidation_MustExist(t *testing.T) {
	if err := safety.ValidateRepoRoot("/nonexistent/path/that/should/not/exist/abc123"); err == nil {
		t.Error("expected non-existent path to be rejected, got nil error")
	}
}

// TestRepoRootValidation_MustBeDirectory verifies that a regular file is rejected.
func TestRepoRootValidation_MustBeDirectory(t *testing.T) {
	dir := t.TempDir()
	f, err := os.CreateTemp(dir, "notadir")
	if err != nil {
		t.Fatalf("setup: %v", err)
	}
	f.Close()
	if err := safety.ValidateRepoRoot(f.Name()); err == nil {
		t.Error("expected file path to be rejected as repo root, got nil error")
	}
}

// --------------------------------------------------------------------------
// Safe path resolution helpers
// --------------------------------------------------------------------------

// TestResolvePath_SimpleRelative verifies that a simple relative path resolves
// to repoRoot/relative.
func TestResolvePath_SimpleRelative(t *testing.T) {
	root := t.TempDir()
	got, err := safety.ResolvePath(root, "subdir/file.tf")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := filepath.Join(root, "subdir", "file.tf")
	if got != want {
		t.Errorf("ResolvePath: got %q, want %q", got, want)
	}
}

// TestResolvePath_DotSlashRelative verifies that a ./ prefix is handled.
func TestResolvePath_DotSlashRelative(t *testing.T) {
	root := t.TempDir()
	got, err := safety.ResolvePath(root, "./main.tf")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := filepath.Join(root, "main.tf")
	if got != want {
		t.Errorf("ResolvePath: got %q, want %q", got, want)
	}
}

// TestResolvePath_TraversalDoubleDot verifies that ../ traversal is rejected.
func TestResolvePath_TraversalDoubleDot(t *testing.T) {
	root := t.TempDir()
	if _, err := safety.ResolvePath(root, "../outside"); err == nil {
		t.Error("expected ../ traversal to be rejected, got nil error")
	}
}

// TestResolvePath_TraversalEmbedded verifies that embedded ../ is rejected.
func TestResolvePath_TraversalEmbedded(t *testing.T) {
	root := t.TempDir()
	if _, err := safety.ResolvePath(root, "subdir/../../outside"); err == nil {
		t.Error("expected embedded ../ traversal to be rejected, got nil error")
	}
}

// TestResolvePath_AbsolutePathRejected verifies that absolute user paths are rejected.
func TestResolvePath_AbsolutePathRejected(t *testing.T) {
	root := t.TempDir()
	if _, err := safety.ResolvePath(root, "/etc/passwd"); err == nil {
		t.Error("expected absolute path to be rejected, got nil error")
	}
}

// TestResolvePath_RepoRootItself verifies that the empty string resolves to root.
func TestResolvePath_RepoRootItself(t *testing.T) {
	root := t.TempDir()
	got, err := safety.ResolvePath(root, ".")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != root {
		t.Errorf("ResolvePath('.'): got %q, want %q", got, root)
	}
}

// --------------------------------------------------------------------------
// Symlink escape tests
// --------------------------------------------------------------------------

// TestResolvePath_SymlinkEscape verifies that a symlink pointing outside the
// repo root is rejected.
func TestResolvePath_SymlinkEscape(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()

	linkPath := filepath.Join(root, "escape")
	if err := os.Symlink(outside, linkPath); err != nil {
		t.Fatalf("setup symlink: %v", err)
	}

	// Create a real file outside so the symlink target resolves.
	if err := os.WriteFile(filepath.Join(outside, "secret"), []byte("top secret"), 0600); err != nil {
		t.Fatalf("setup file: %v", err)
	}

	if _, err := safety.ResolvePath(root, "escape/secret"); err == nil {
		t.Error("expected symlink escape to be rejected, got nil error")
	}
}

// TestResolvePath_SymlinkInsideRepoOK verifies that symlinks whose target stays
// within the repo root are accepted.
func TestResolvePath_SymlinkInsideRepoOK(t *testing.T) {
	root := t.TempDir()

	// Create a real target inside the repo.
	targetDir := filepath.Join(root, "real")
	if err := os.MkdirAll(targetDir, 0750); err != nil {
		t.Fatalf("setup: %v", err)
	}
	if err := os.WriteFile(filepath.Join(targetDir, "main.tf"), []byte(""), 0600); err != nil {
		t.Fatalf("setup: %v", err)
	}

	// Create a symlink inside the repo pointing to the target inside the repo.
	linkDir := filepath.Join(root, "link")
	if err := os.Symlink(targetDir, linkDir); err != nil {
		t.Fatalf("setup symlink: %v", err)
	}

	if _, err := safety.ResolvePath(root, "link/main.tf"); err != nil {
		t.Errorf("expected in-repo symlink to be accepted, got error: %v", err)
	}
}

// --------------------------------------------------------------------------
// Forbidden path tests
// --------------------------------------------------------------------------

// TestResolvePath_DotGitRejected verifies paths into .git are rejected.
func TestResolvePath_DotGitRejected(t *testing.T) {
	root := t.TempDir()
	if _, err := safety.ResolvePath(root, ".git/config"); err == nil {
		t.Error("expected .git path to be rejected, got nil error")
	}
}

// TestResolvePath_DotGitDirectRejected verifies the .git directory itself is rejected.
func TestResolvePath_DotGitDirectRejected(t *testing.T) {
	root := t.TempDir()
	if _, err := safety.ResolvePath(root, ".git"); err == nil {
		t.Error("expected .git path to be rejected, got nil error")
	}
}

// --------------------------------------------------------------------------
// Default exclusion tests (.terraform)
// --------------------------------------------------------------------------

// TestIsExcludedByDefault_DotTerraform verifies .terraform is excluded by default.
func TestIsExcludedByDefault_DotTerraform(t *testing.T) {
	if !safety.IsExcludedByDefault(".terraform") {
		t.Error("expected .terraform to be excluded by default")
	}
}

// TestIsExcludedByDefault_DotTerraformSubpath verifies subpaths of .terraform are excluded.
func TestIsExcludedByDefault_DotTerraformSubpath(t *testing.T) {
	if !safety.IsExcludedByDefault(".terraform/providers/plugin") {
		t.Error("expected .terraform subpath to be excluded by default")
	}
}

// TestIsExcludedByDefault_DotGit verifies .git is excluded by default.
func TestIsExcludedByDefault_DotGit(t *testing.T) {
	if !safety.IsExcludedByDefault(".git") {
		t.Error("expected .git to be excluded by default")
	}
}

// TestIsExcludedByDefault_NormalFile verifies normal files are not excluded.
func TestIsExcludedByDefault_NormalFile(t *testing.T) {
	if safety.IsExcludedByDefault("main.tf") {
		t.Error("expected main.tf to NOT be excluded by default")
	}
}

// TestIsExcludedByDefault_NormalSubdir verifies normal subdirectories are not excluded.
func TestIsExcludedByDefault_NormalSubdir(t *testing.T) {
	if safety.IsExcludedByDefault("modules/vpc/main.tf") {
		t.Error("expected modules/vpc/main.tf to NOT be excluded by default")
	}
}
