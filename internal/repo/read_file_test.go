package repo_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/maazghani/terraformer/internal/repo"
)

// --------------------------------------------------------------------------
// ReadFile
// --------------------------------------------------------------------------

// TestReadFile_BasicRead verifies that a regular file can be read.
func TestReadFile_BasicRead(t *testing.T) {
	root := t.TempDir()
	content := "resource \"local_file\" \"example\" {}"
	writeFile(t, filepath.Join(root, "main.tf"), content)

	svc, _ := repo.New(root)
	resp, err := svc.ReadFile(repo.ReadFileRequest{Path: "main.tf"})
	if err != nil {
		t.Fatalf("ReadFile: unexpected error: %v", err)
	}

	if resp.Content != content {
		t.Errorf("ReadFile: expected content %q, got %q", content, resp.Content)
	}
	if resp.SizeBytes != int64(len(content)) {
		t.Errorf("ReadFile: expected SizeBytes=%d, got %d", len(content), resp.SizeBytes)
	}
	if resp.Truncated {
		t.Error("ReadFile: expected Truncated=false for full read")
	}
}

// TestReadFile_MissingFile verifies that a missing file returns an error.
func TestReadFile_MissingFile(t *testing.T) {
	root := t.TempDir()
	svc, _ := repo.New(root)

	_, err := svc.ReadFile(repo.ReadFileRequest{Path: "nonexistent.tf"})
	if err == nil {
		t.Fatal("ReadFile: expected error for missing file, got nil")
	}
}

// TestReadFile_DirectoryRejected verifies that reading a directory returns an error.
func TestReadFile_DirectoryRejected(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "subdir"), 0o755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	svc, _ := repo.New(root)
	_, err := svc.ReadFile(repo.ReadFileRequest{Path: "subdir"})
	if err == nil {
		t.Fatal("ReadFile: expected error when reading a directory, got nil")
	}
}

// TestReadFile_MaxBytes verifies that content is truncated to max_bytes.
func TestReadFile_MaxBytes(t *testing.T) {
	root := t.TempDir()
	content := "0123456789" // 10 bytes
	writeFile(t, filepath.Join(root, "data.txt"), content)

	svc, _ := repo.New(root)
	resp, err := svc.ReadFile(repo.ReadFileRequest{Path: "data.txt", MaxBytes: 5})
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	if len(resp.Content) != 5 {
		t.Errorf("ReadFile: expected 5 bytes of content, got %d", len(resp.Content))
	}
	if !resp.Truncated {
		t.Error("ReadFile: expected Truncated=true when content is truncated")
	}
}

// TestReadFile_MaxBytes_NoTruncation verifies that Truncated=false when file
// fits within max_bytes.
func TestReadFile_MaxBytes_NoTruncation(t *testing.T) {
	root := t.TempDir()
	content := "hello"
	writeFile(t, filepath.Join(root, "a.tf"), content)

	svc, _ := repo.New(root)
	resp, err := svc.ReadFile(repo.ReadFileRequest{Path: "a.tf", MaxBytes: 1024})
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	if resp.Truncated {
		t.Error("ReadFile: expected Truncated=false when content fits within max_bytes")
	}
	if resp.Content != content {
		t.Errorf("ReadFile: expected %q, got %q", content, resp.Content)
	}
}

// TestReadFile_TraversalRejected verifies that ../ traversal is rejected.
func TestReadFile_TraversalRejected(t *testing.T) {
	root := t.TempDir()
	svc, _ := repo.New(root)

	_, err := svc.ReadFile(repo.ReadFileRequest{Path: "../outside.tf"})
	if err == nil {
		t.Fatal("ReadFile: expected error for traversal path, got nil")
	}
}

// TestReadFile_AbsolutePathRejected verifies that an absolute path is rejected.
func TestReadFile_AbsolutePathRejected(t *testing.T) {
	root := t.TempDir()
	svc, _ := repo.New(root)

	_, err := svc.ReadFile(repo.ReadFileRequest{Path: "/etc/passwd"})
	if err == nil {
		t.Fatal("ReadFile: expected error for absolute path, got nil")
	}
}

// TestReadFile_GitPathRejected verifies that .git paths are rejected.
func TestReadFile_GitPathRejected(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, ".git", "config"), "[core]")

	svc, _ := repo.New(root)
	_, err := svc.ReadFile(repo.ReadFileRequest{Path: ".git/config"})
	if err == nil {
		t.Fatal("ReadFile: expected error reading .git/config, got nil")
	}
}

// TestReadFile_SymlinkEscape verifies that a symlink escaping the repo root is rejected.
func TestReadFile_SymlinkEscape(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	writeFile(t, filepath.Join(outside, "secret.txt"), "top-secret")

	link := filepath.Join(root, "link.txt")
	if err := os.Symlink(filepath.Join(outside, "secret.txt"), link); err != nil {
		t.Skipf("could not create symlink: %v", err)
	}

	svc, _ := repo.New(root)
	_, err := svc.ReadFile(repo.ReadFileRequest{Path: "link.txt"})
	if err == nil {
		t.Fatal("ReadFile: expected error for symlink escape, got nil")
	}
}

// TestReadFile_SizeBytesIsActualSize verifies SizeBytes reflects the full file
// size, not the truncated content length.
func TestReadFile_SizeBytesIsActualSize(t *testing.T) {
	root := t.TempDir()
	content := "0123456789" // 10 bytes
	writeFile(t, filepath.Join(root, "data.txt"), content)

	svc, _ := repo.New(root)
	resp, err := svc.ReadFile(repo.ReadFileRequest{Path: "data.txt", MaxBytes: 5})
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	if resp.SizeBytes != int64(len(content)) {
		t.Errorf("ReadFile: expected SizeBytes=%d (full size), got %d", len(content), resp.SizeBytes)
	}
}
