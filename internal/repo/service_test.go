package repo_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/maazghani/terraformer/internal/repo"
)

// helper: create a file at path with content.
func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
}

// --------------------------------------------------------------------------
// repo.New
// --------------------------------------------------------------------------

// TestNew_AcceptsValidRoot verifies that New succeeds with an existing directory.
func TestNew_AcceptsValidRoot(t *testing.T) {
	root := t.TempDir()
	if _, err := repo.New(root); err != nil {
		t.Fatalf("New: unexpected error: %v", err)
	}
}

// TestNew_RejectsRelativeRoot verifies that New rejects a relative root path.
func TestNew_RejectsRelativeRoot(t *testing.T) {
	if _, err := repo.New("relative/path"); err == nil {
		t.Fatal("New: expected error for relative root, got nil")
	}
}

// TestNew_RejectsNonexistentRoot verifies that New rejects a missing directory.
func TestNew_RejectsNonexistentRoot(t *testing.T) {
	if _, err := repo.New("/nonexistent/path/abc123xyz"); err == nil {
		t.Fatal("New: expected error for nonexistent root, got nil")
	}
}

// --------------------------------------------------------------------------
// ListFiles
// --------------------------------------------------------------------------

// TestListFiles_BasicListing verifies that files are returned and paths are relative.
func TestListFiles_BasicListing(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "main.tf"), "# main")
	writeFile(t, filepath.Join(root, "vars.tf"), "# vars")

	svc, err := repo.New(root)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	resp, err := svc.ListFiles(repo.ListFilesRequest{Path: "."})
	if err != nil {
		t.Fatalf("ListFiles: %v", err)
	}

	if len(resp.Files) != 2 {
		t.Errorf("expected 2 files, got %d: %v", len(resp.Files), resp.Files)
	}
	// All paths must be relative (no leading /).
	for _, f := range resp.Files {
		if filepath.IsAbs(f.Path) {
			t.Errorf("file path is absolute: %q", f.Path)
		}
	}
}

// TestListFiles_ExcludesGitDir verifies that .git is excluded from listing.
func TestListFiles_ExcludesGitDir(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, ".git", "config"), "[core]")
	writeFile(t, filepath.Join(root, "main.tf"), "# main")

	svc, _ := repo.New(root)
	resp, err := svc.ListFiles(repo.ListFilesRequest{Path: "."})
	if err != nil {
		t.Fatalf("ListFiles: %v", err)
	}

	for _, f := range resp.Files {
		if f.Path == ".git/config" || f.Path == ".git" {
			t.Errorf("expected .git to be excluded, got %q", f.Path)
		}
	}
}

// TestListFiles_ExcludesTerraformDir verifies that .terraform is excluded from listing.
func TestListFiles_ExcludesTerraformDir(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, ".terraform", "providers", "lock.hcl"), "# lock")
	writeFile(t, filepath.Join(root, "main.tf"), "# main")

	svc, _ := repo.New(root)
	resp, err := svc.ListFiles(repo.ListFilesRequest{Path: "."})
	if err != nil {
		t.Fatalf("ListFiles: %v", err)
	}

	for _, f := range resp.Files {
		if f.Path == ".terraform" || filepath.HasPrefix(f.Path, ".terraform/") {
			t.Errorf("expected .terraform to be excluded, got %q", f.Path)
		}
	}
}

// TestListFiles_MaxFiles verifies that max_files limits the listing and
// sets Truncated=true when the limit is hit.
func TestListFiles_MaxFiles(t *testing.T) {
	root := t.TempDir()
	for i := 0; i < 5; i++ {
		name := filepath.Join(root, "f"+string(rune('a'+i))+".tf")
		writeFile(t, name, "# content")
	}

	svc, _ := repo.New(root)
	resp, err := svc.ListFiles(repo.ListFilesRequest{Path: ".", MaxFiles: 3})
	if err != nil {
		t.Fatalf("ListFiles: %v", err)
	}

	if len(resp.Files) != 3 {
		t.Errorf("expected 3 files (capped by max_files), got %d", len(resp.Files))
	}
	if !resp.Truncated {
		t.Error("expected Truncated=true when max_files is exceeded")
	}
}

// TestListFiles_MaxFiles_NoTruncation verifies that Truncated=false when
// the file count is under the max_files limit.
func TestListFiles_MaxFiles_NoTruncation(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "a.tf"), "# a")
	writeFile(t, filepath.Join(root, "b.tf"), "# b")

	svc, _ := repo.New(root)
	resp, err := svc.ListFiles(repo.ListFilesRequest{Path: ".", MaxFiles: 100})
	if err != nil {
		t.Fatalf("ListFiles: %v", err)
	}

	if resp.Truncated {
		t.Error("expected Truncated=false when count is under max_files")
	}
}

// TestListFiles_NormalizedRelativePaths verifies that nested file paths use
// slash-separated relative paths from repo root.
func TestListFiles_NormalizedRelativePaths(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "subdir", "nested.tf"), "# nested")

	svc, _ := repo.New(root)
	resp, err := svc.ListFiles(repo.ListFilesRequest{Path: "."})
	if err != nil {
		t.Fatalf("ListFiles: %v", err)
	}

	found := false
	for _, f := range resp.Files {
		if f.Path == "subdir/nested.tf" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected subdir/nested.tf in listing, got %v", resp.Files)
	}
}

// TestListFiles_TraversalRejected verifies that a traversal path is rejected.
func TestListFiles_TraversalRejected(t *testing.T) {
	root := t.TempDir()
	svc, _ := repo.New(root)
	_, err := svc.ListFiles(repo.ListFilesRequest{Path: "../outside"})
	if err == nil {
		t.Fatal("expected traversal to be rejected, got nil error")
	}
}

// TestListFiles_AbsolutePathRejected verifies that an absolute path is rejected.
func TestListFiles_AbsolutePathRejected(t *testing.T) {
	root := t.TempDir()
	svc, _ := repo.New(root)
	_, err := svc.ListFiles(repo.ListFilesRequest{Path: "/etc"})
	if err == nil {
		t.Fatal("expected absolute path to be rejected, got nil error")
	}
}

// TestListFiles_SymlinkEscapeRejected verifies that a symlink escaping the repo
// root is skipped or causes an error.
func TestListFiles_SymlinkEscapeRejected(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	writeFile(t, filepath.Join(outside, "secret.txt"), "secret")

	// Create a symlink inside root pointing to outside.
	link := filepath.Join(root, "escape_link")
	if err := os.Symlink(outside, link); err != nil {
		t.Skipf("could not create symlink (may require elevated privileges): %v", err)
	}

	svc, _ := repo.New(root)
	resp, err := svc.ListFiles(repo.ListFilesRequest{Path: "."})
	if err != nil {
		// An error is acceptable—symlink detection failing safely.
		return
	}

	// The symlink target contents must not appear in the listing.
	for _, f := range resp.Files {
		if f.Path == "escape_link/secret.txt" {
			t.Errorf("symlink escape: escape_link/secret.txt should not appear in listing")
		}
	}
}

// TestListFiles_FileEntryHasSizeBytes verifies that FileEntry.SizeBytes is populated.
func TestListFiles_FileEntryHasSizeBytes(t *testing.T) {
	root := t.TempDir()
	content := "hello world"
	writeFile(t, filepath.Join(root, "a.tf"), content)

	svc, _ := repo.New(root)
	resp, err := svc.ListFiles(repo.ListFilesRequest{Path: "."})
	if err != nil {
		t.Fatalf("ListFiles: %v", err)
	}

	if len(resp.Files) == 0 {
		t.Fatal("expected at least one file")
	}
	if resp.Files[0].SizeBytes != int64(len(content)) {
		t.Errorf("expected SizeBytes=%d, got %d", len(content), resp.Files[0].SizeBytes)
	}
}

// TestListFiles_FileEntryKindIsFile verifies that regular files have Kind="file".
func TestListFiles_FileEntryKindIsFile(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "a.tf"), "# content")

	svc, _ := repo.New(root)
	resp, err := svc.ListFiles(repo.ListFilesRequest{Path: "."})
	if err != nil {
		t.Fatalf("ListFiles: %v", err)
	}

	for _, f := range resp.Files {
		if f.Kind != "file" {
			t.Errorf("expected Kind=file, got %q for %q", f.Kind, f.Path)
		}
	}
}
