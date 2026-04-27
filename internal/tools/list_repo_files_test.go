package tools_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/maazghani/terraformer/internal/repo"
	"github.com/maazghani/terraformer/internal/tools"
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
// ListRepoFiles handler
// --------------------------------------------------------------------------

// TestListRepoFiles_OK verifies basic listing returns ok=true and files.
func TestListRepoFiles_OK(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "main.tf"), "# main")

	svc, err := repo.New(root)
	if err != nil {
		t.Fatalf("repo.New: %v", err)
	}

	resp := tools.ListRepoFiles(svc, tools.ListRepoFilesRequest{Path: "."})
	if !resp.OK {
		t.Fatalf("expected OK=true, got OK=false; warnings: %v", resp.Warnings)
	}
	if len(resp.Files) != 1 {
		t.Errorf("expected 1 file, got %d", len(resp.Files))
	}
}

// TestListRepoFiles_Traversal verifies that traversal path returns ok=false.
func TestListRepoFiles_Traversal(t *testing.T) {
	root := t.TempDir()
	svc, _ := repo.New(root)

	resp := tools.ListRepoFiles(svc, tools.ListRepoFilesRequest{Path: "../outside"})
	if resp.OK {
		t.Fatal("expected OK=false for traversal path, got OK=true")
	}
}

// TestListRepoFiles_Truncated verifies that max_files cap is reflected in response.
func TestListRepoFiles_Truncated(t *testing.T) {
	root := t.TempDir()
	for i := 0; i < 5; i++ {
		writeFile(t, filepath.Join(root, "f"+string(rune('a'+i))+".tf"), "# content")
	}

	svc, _ := repo.New(root)
	resp := tools.ListRepoFiles(svc, tools.ListRepoFilesRequest{Path: ".", MaxFiles: 2})
	if !resp.OK {
		t.Fatalf("expected OK=true, got OK=false")
	}
	if len(resp.Files) != 2 {
		t.Errorf("expected 2 files, got %d", len(resp.Files))
	}
	if !resp.Truncated {
		t.Error("expected Truncated=true")
	}
}

// TestListRepoFiles_DefaultMaxFiles verifies that a zero MaxFiles uses a sensible default.
func TestListRepoFiles_DefaultMaxFiles(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "a.tf"), "# a")

	svc, _ := repo.New(root)
	resp := tools.ListRepoFiles(svc, tools.ListRepoFilesRequest{Path: "."})
	if !resp.OK {
		t.Fatalf("expected OK=true, got OK=false")
	}
	if len(resp.Files) == 0 {
		t.Error("expected at least one file")
	}
}

// TestListRepoFiles_ExcludesGit verifies .git is excluded in the tool response.
func TestListRepoFiles_ExcludesGit(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, ".git", "config"), "[core]")
	writeFile(t, filepath.Join(root, "main.tf"), "# main")

	svc, _ := repo.New(root)
	resp := tools.ListRepoFiles(svc, tools.ListRepoFilesRequest{Path: "."})

	for _, f := range resp.Files {
		if f.Path == ".git/config" || f.Path == ".git" {
			t.Errorf("expected .git to be excluded, found %q", f.Path)
		}
	}
}

// TestListRepoFiles_ResponseShape verifies the response has all required fields.
func TestListRepoFiles_ResponseShape(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "main.tf"), "# main")

	svc, _ := repo.New(root)
	resp := tools.ListRepoFiles(svc, tools.ListRepoFilesRequest{Path: "."})

	// Verify required fields exist and have expected zero-value or populated values.
	if resp.Files == nil {
		t.Error("Files must not be nil")
	}
	// Warnings must be non-nil (empty slice, not nil) per structured output contract.
	if resp.Warnings == nil {
		t.Error("Warnings must not be nil")
	}
}
