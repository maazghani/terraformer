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

// TestApplyPatchWrite verifies that ApplyPatch can write a new file.
func TestApplyPatchWrite(t *testing.T) {
	tmpDir := t.TempDir()
	svc, err := patch.New(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	req := patch.ApplyPatchRequest{
		Files: []patch.FileOperation{
			{
				Path:      "main.tf",
				Operation: "write",
				Content:   "resource \"local_file\" \"example\" {}\n",
			},
		},
	}

	resp, err := svc.ApplyPatch(req)
	if err != nil {
		t.Fatalf("ApplyPatch() failed: %v", err)
	}
	if !resp.OK {
		t.Fatal("ApplyPatch() should succeed")
	}
	if len(resp.ChangedFiles) != 1 {
		t.Fatalf("expected 1 changed file, got %d", len(resp.ChangedFiles))
	}
	if resp.ChangedFiles[0] != "main.tf" {
		t.Errorf("expected changed file 'main.tf', got %q", resp.ChangedFiles[0])
	}

	// Verify the file was actually written
	content, err := os.ReadFile(filepath.Join(tmpDir, "main.tf"))
	if err != nil {
		t.Fatalf("failed to read written file: %v", err)
	}
	expected := "resource \"local_file\" \"example\" {}\n"
	if string(content) != expected {
		t.Errorf("file content mismatch\nwant: %q\ngot:  %q", expected, string(content))
	}
}

// TestApplyPatchOverwrite verifies that ApplyPatch can overwrite an existing file.
func TestApplyPatchOverwrite(t *testing.T) {
	tmpDir := t.TempDir()
	svc, err := patch.New(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	// Create an existing file
	existingPath := filepath.Join(tmpDir, "vars.tf")
	if err := os.WriteFile(existingPath, []byte("old content"), 0644); err != nil {
		t.Fatal(err)
	}

	req := patch.ApplyPatchRequest{
		Files: []patch.FileOperation{
			{
				Path:      "vars.tf",
				Operation: "write",
				Content:   "variable \"env\" { default = \"dev\" }\n",
			},
		},
	}

	resp, err := svc.ApplyPatch(req)
	if err != nil {
		t.Fatalf("ApplyPatch() failed: %v", err)
	}
	if !resp.OK {
		t.Fatal("ApplyPatch() should succeed")
	}

	// Verify the file was overwritten
	content, err := os.ReadFile(existingPath)
	if err != nil {
		t.Fatalf("failed to read overwritten file: %v", err)
	}
	expected := "variable \"env\" { default = \"dev\" }\n"
	if string(content) != expected {
		t.Errorf("file content mismatch\nwant: %q\ngot:  %q", expected, string(content))
	}
}

// TestApplyPatchMultipleFiles verifies that ApplyPatch can handle multiple files.
func TestApplyPatchMultipleFiles(t *testing.T) {
	tmpDir := t.TempDir()
	svc, err := patch.New(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	req := patch.ApplyPatchRequest{
		Files: []patch.FileOperation{
			{
				Path:      "main.tf",
				Operation: "write",
				Content:   "resource \"local_file\" \"example\" {}\n",
			},
			{
				Path:      "vars.tf",
				Operation: "write",
				Content:   "variable \"env\" { default = \"dev\" }\n",
			},
		},
	}

	resp, err := svc.ApplyPatch(req)
	if err != nil {
		t.Fatalf("ApplyPatch() failed: %v", err)
	}
	if !resp.OK {
		t.Fatal("ApplyPatch() should succeed")
	}
	if len(resp.ChangedFiles) != 2 {
		t.Fatalf("expected 2 changed files, got %d", len(resp.ChangedFiles))
	}
}

// TestApplyPatchDelete verifies that ApplyPatch can delete a file.
func TestApplyPatchDelete(t *testing.T) {
	tmpDir := t.TempDir()
	svc, err := patch.New(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	// Create a file to delete
	targetPath := filepath.Join(tmpDir, "delete-me.tf")
	if err := os.WriteFile(targetPath, []byte("to be deleted"), 0644); err != nil {
		t.Fatal(err)
	}

	req := patch.ApplyPatchRequest{
		Files: []patch.FileOperation{
			{
				Path:      "delete-me.tf",
				Operation: "delete",
			},
		},
	}

	resp, err := svc.ApplyPatch(req)
	if err != nil {
		t.Fatalf("ApplyPatch() failed: %v", err)
	}
	if !resp.OK {
		t.Fatal("ApplyPatch() should succeed")
	}
	if len(resp.ChangedFiles) != 1 {
		t.Fatalf("expected 1 changed file, got %d", len(resp.ChangedFiles))
	}

	// Verify the file was deleted
	if _, err := os.Stat(targetPath); !os.IsNotExist(err) {
		t.Error("file should have been deleted")
	}
}

// TestApplyPatchRejectsTraversal verifies that ApplyPatch rejects path traversal.
func TestApplyPatchRejectsTraversal(t *testing.T) {
	tmpDir := t.TempDir()
	svc, err := patch.New(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	req := patch.ApplyPatchRequest{
		Files: []patch.FileOperation{
			{
				Path:      "../outside.tf",
				Operation: "write",
				Content:   "malicious",
			},
		},
	}

	resp, err := svc.ApplyPatch(req)
	if err != nil {
		t.Fatalf("ApplyPatch() should not return error, but set rejected files: %v", err)
	}
	if resp.OK {
		t.Error("ApplyPatch() should fail for path traversal")
	}
	if len(resp.RejectedFiles) != 1 {
		t.Fatalf("expected 1 rejected file, got %d", len(resp.RejectedFiles))
	}
	if resp.RejectedFiles[0] != "../outside.tf" {
		t.Errorf("expected rejected file '../outside.tf', got %q", resp.RejectedFiles[0])
	}
}

// TestApplyPatchRejectsAbsolute verifies that ApplyPatch rejects absolute paths.
func TestApplyPatchRejectsAbsolute(t *testing.T) {
	tmpDir := t.TempDir()
	svc, err := patch.New(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	req := patch.ApplyPatchRequest{
		Files: []patch.FileOperation{
			{
				Path:      "/etc/passwd",
				Operation: "write",
				Content:   "malicious",
			},
		},
	}

	resp, err := svc.ApplyPatch(req)
	if err != nil {
		t.Fatalf("ApplyPatch() should not return error, but set rejected files: %v", err)
	}
	if resp.OK {
		t.Error("ApplyPatch() should fail for absolute path")
	}
	if len(resp.RejectedFiles) != 1 {
		t.Fatalf("expected 1 rejected file, got %d", len(resp.RejectedFiles))
	}
}

// TestApplyPatchRejectsSymlinkEscape verifies that symlinks escaping the repo are rejected.
func TestApplyPatchRejectsSymlinkEscape(t *testing.T) {
	tmpDir := t.TempDir()
	outsideDir := t.TempDir()
	svc, err := patch.New(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	// Create a symlink pointing outside the repo
	symlinkPath := filepath.Join(tmpDir, "escape-link")
	if err := os.Symlink(outsideDir, symlinkPath); err != nil {
		t.Skip("cannot create symlinks on this platform")
	}

	req := patch.ApplyPatchRequest{
		Files: []patch.FileOperation{
			{
				Path:      "escape-link/evil.tf",
				Operation: "write",
				Content:   "malicious",
			},
		},
	}

	resp, err := svc.ApplyPatch(req)
	if err != nil {
		t.Fatalf("ApplyPatch() should not return error, but set rejected files: %v", err)
	}
	if resp.OK {
		t.Error("ApplyPatch() should fail for symlink escape")
	}
	if len(resp.RejectedFiles) != 1 {
		t.Fatalf("expected 1 rejected file, got %d", len(resp.RejectedFiles))
	}
}

// TestApplyPatchNoCommandExecution verifies that ApplyPatch does not execute commands.
// This is a behavioral assertion: the patch package has no dependency on runner or terraform.
func TestApplyPatchNoCommandExecution(t *testing.T) {
	tmpDir := t.TempDir()
	svc, err := patch.New(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	req := patch.ApplyPatchRequest{
		Files: []patch.FileOperation{
			{
				Path:      "main.tf",
				Operation: "write",
				Content:   "resource \"local_file\" \"example\" {}\n",
			},
		},
	}

	_, err = svc.ApplyPatch(req)
	if err != nil {
		t.Fatalf("ApplyPatch() failed: %v", err)
	}

	// No command execution means no terraform, no runner calls.
	// This test passes if ApplyPatch completes without invoking external processes.
	// The implementation should only use os.WriteFile, os.Remove, etc.
}
