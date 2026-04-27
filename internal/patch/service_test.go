package patch_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/maazghani/terraformer/internal/patch"
)

// TestNew verifies that New creates a Service for a valid repo root.
func TestNew(t *testing.T) {
	tmpDir := t.TempDir()

	svc, err := patch.New(tmpDir)
	if err != nil {
		t.Fatalf("patch.New(%q) returned error: %v", tmpDir, err)
	}
	if svc == nil {
		t.Fatal("patch.New() returned nil Service")
	}
}

// TestNewRejectsNonExistent verifies that New rejects a non-existent root.
func TestNewRejectsNonExistent(t *testing.T) {
	nonExistent := filepath.Join(t.TempDir(), "does-not-exist")

	_, err := patch.New(nonExistent)
	if err == nil {
		t.Fatal("patch.New() with non-existent root should fail")
	}
}

// TestNewRejectsRelativePath verifies that New rejects a relative path.
func TestNewRejectsRelativePath(t *testing.T) {
	_, err := patch.New("relative/path")
	if err == nil {
		t.Fatal("patch.New() with relative path should fail")
	}
}

// TestNewRejectsFile verifies that New rejects a file path as the root.
func TestNewRejectsFile(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "file.txt")
	if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := patch.New(filePath)
	if err == nil {
		t.Fatal("patch.New() with file path should fail")
	}
}
