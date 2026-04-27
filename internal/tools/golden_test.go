package tools_test

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/maazghani/terraformer/internal/repo"
	"github.com/maazghani/terraformer/internal/tools"
)

// update controls whether golden files are regenerated. Run with -update to update.
var update = flag.Bool("update", false, "update golden files")

// goldenPath returns the path to a golden file by name.
func goldenPath(name string) string {
	return filepath.Join("..", "..", "testdata", "golden", name)
}

// checkGolden compares got (as JSON) against the named golden file.
// When -update is set, the golden file is written with got.
func checkGolden(t *testing.T, name string, got interface{}) {
	t.Helper()

	data, err := json.MarshalIndent(got, "", "  ")
	if err != nil {
		t.Fatalf("golden: failed to marshal response: %v", err)
	}
	actual := string(data) + "\n"

	path := goldenPath(name)

	if *update {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("golden: mkdir: %v", err)
		}
		if err := os.WriteFile(path, []byte(actual), 0o644); err != nil {
			t.Fatalf("golden: write: %v", err)
		}
		t.Logf("golden: updated %s", path)
		return
	}

	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("golden: missing golden file %q — run with -update to create it: %v", path, err)
	}

	if actual != string(want) {
		t.Errorf("golden: response shape mismatch for %q\ngot:\n%s\nwant:\n%s", name, actual, want)
	}
}

// TestGolden_ListRepoFilesResponse verifies the list_repo_files response shape.
func TestGolden_ListRepoFilesResponse(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "main.tf"), "# main terraform")
	writeFile(t, filepath.Join(root, "vars.tf"), "# variables")
	writeFile(t, filepath.Join(root, ".git", "config"), "[core]")
	writeFile(t, filepath.Join(root, ".terraform", "lock.hcl"), "# lock")

	svc, err := repo.New(root)
	if err != nil {
		t.Fatalf("repo.New: %v", err)
	}

	resp := tools.ListRepoFiles(svc, tools.ListRepoFilesRequest{Path: ".", MaxFiles: 200})

	// Replace size_bytes with a stable value for golden comparison since
	// actual sizes depend on content written above.
	// We compare structural shape rather than exact byte counts.
	type stableFileInfo struct {
		Path string `json:"path"`
		Kind string `json:"kind"`
	}
	type stableListResponse struct {
		OK        bool             `json:"ok"`
		Files     []stableFileInfo `json:"files"`
		Truncated bool             `json:"truncated"`
		Warnings  []string         `json:"warnings"`
	}

	stable := stableListResponse{
		OK:        resp.OK,
		Truncated: resp.Truncated,
		Warnings:  resp.Warnings,
	}
	for _, f := range resp.Files {
		stable.Files = append(stable.Files, stableFileInfo{Path: f.Path, Kind: f.Kind})
	}
	if stable.Files == nil {
		stable.Files = []stableFileInfo{}
	}

	checkGolden(t, "list_repo_files_response.json", stable)
}

// TestGolden_ReadRepoFileResponse verifies the read_repo_file response shape.
func TestGolden_ReadRepoFileResponse(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "main.tf"), "resource \"local_file\" \"example\" {}\n")

	svc, err := repo.New(root)
	if err != nil {
		t.Fatalf("repo.New: %v", err)
	}

	resp := tools.ReadRepoFile(svc, tools.ReadRepoFileRequest{Path: "main.tf"})

	// Use a stable shape (normalize size_bytes).
	type stableReadResponse struct {
		OK        bool     `json:"ok"`
		Path      string   `json:"path"`
		Content   string   `json:"content"`
		Truncated bool     `json:"truncated"`
		Warnings  []string `json:"warnings"`
	}

	stable := stableReadResponse{
		OK:        resp.OK,
		Path:      resp.Path,
		Content:   resp.Content,
		Truncated: resp.Truncated,
		Warnings:  resp.Warnings,
	}

	checkGolden(t, "read_repo_file_response.json", stable)
}

// TestGolden_ReadRepoFileResponse_Truncated verifies the truncated response shape.
func TestGolden_ReadRepoFileResponse_Truncated(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "big.tf"), "0123456789ABCDEF")

	svc, err := repo.New(root)
	if err != nil {
		t.Fatalf("repo.New: %v", err)
	}

	resp := tools.ReadRepoFile(svc, tools.ReadRepoFileRequest{Path: "big.tf", MaxBytes: 8})

	type stableReadResponse struct {
		OK        bool     `json:"ok"`
		Path      string   `json:"path"`
		Content   string   `json:"content"`
		Truncated bool     `json:"truncated"`
		Warnings  []string `json:"warnings"`
	}

	stable := stableReadResponse{
		OK:        resp.OK,
		Path:      resp.Path,
		Content:   resp.Content,
		Truncated: resp.Truncated,
		Warnings:  resp.Warnings,
	}

	checkGolden(t, "read_repo_file_response_truncated.json", stable)
}
