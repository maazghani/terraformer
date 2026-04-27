package tools_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/maazghani/terraformer/internal/repo"
	"github.com/maazghani/terraformer/internal/tools"
)

// --------------------------------------------------------------------------
// ReadRepoFile handler
// --------------------------------------------------------------------------

// TestReadRepoFile_OK verifies that a valid file read returns ok=true and content.
func TestReadRepoFile_OK(t *testing.T) {
	root := t.TempDir()
	content := "resource \"aws_s3_bucket\" \"b\" {}"
	writeFile(t, filepath.Join(root, "main.tf"), content)

	svc, _ := repo.New(root)
	resp := tools.ReadRepoFile(svc, tools.ReadRepoFileRequest{Path: "main.tf"})

	if !resp.OK {
		t.Fatalf("expected OK=true, got OK=false; warnings: %v", resp.Warnings)
	}
	if resp.Content != content {
		t.Errorf("expected content %q, got %q", content, resp.Content)
	}
	if resp.Path != "main.tf" {
		t.Errorf("expected Path=main.tf, got %q", resp.Path)
	}
}

// TestReadRepoFile_MissingFile verifies that missing file returns ok=false.
func TestReadRepoFile_MissingFile(t *testing.T) {
	root := t.TempDir()
	svc, _ := repo.New(root)

	resp := tools.ReadRepoFile(svc, tools.ReadRepoFileRequest{Path: "nonexistent.tf"})
	if resp.OK {
		t.Fatal("expected OK=false for missing file, got OK=true")
	}
}

// TestReadRepoFile_DirectoryRejected verifies that reading a directory returns ok=false.
func TestReadRepoFile_DirectoryRejected(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "subdir"), 0o755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	svc, _ := repo.New(root)
	resp := tools.ReadRepoFile(svc, tools.ReadRepoFileRequest{Path: "subdir"})
	if resp.OK {
		t.Fatal("expected OK=false for directory read, got OK=true")
	}
}

// TestReadRepoFile_MaxBytes verifies that content is truncated and Truncated=true.
func TestReadRepoFile_MaxBytes(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "big.tf"), "0123456789")

	svc, _ := repo.New(root)
	resp := tools.ReadRepoFile(svc, tools.ReadRepoFileRequest{Path: "big.tf", MaxBytes: 5})

	if !resp.OK {
		t.Fatalf("expected OK=true, got OK=false")
	}
	if len(resp.Content) != 5 {
		t.Errorf("expected 5 bytes of content, got %d", len(resp.Content))
	}
	if !resp.Truncated {
		t.Error("expected Truncated=true when max_bytes is exceeded")
	}
}

// TestReadRepoFile_Traversal verifies that traversal path returns ok=false.
func TestReadRepoFile_Traversal(t *testing.T) {
	root := t.TempDir()
	svc, _ := repo.New(root)

	resp := tools.ReadRepoFile(svc, tools.ReadRepoFileRequest{Path: "../outside.tf"})
	if resp.OK {
		t.Fatal("expected OK=false for traversal, got OK=true")
	}
}

// TestReadRepoFile_GitPathRejected verifies that .git paths return ok=false.
func TestReadRepoFile_GitPathRejected(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, ".git", "config"), "[core]")

	svc, _ := repo.New(root)
	resp := tools.ReadRepoFile(svc, tools.ReadRepoFileRequest{Path: ".git/config"})
	if resp.OK {
		t.Fatal("expected OK=false for .git path, got OK=true")
	}
}

// TestReadRepoFile_ResponseShape verifies all response fields are populated.
func TestReadRepoFile_ResponseShape(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "a.tf"), "hello")

	svc, _ := repo.New(root)
	resp := tools.ReadRepoFile(svc, tools.ReadRepoFileRequest{Path: "a.tf"})

	if resp.Warnings == nil {
		t.Error("Warnings must not be nil")
	}
	if resp.SizeBytes == 0 {
		t.Error("expected non-zero SizeBytes")
	}
}
