package httpserver_test

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/maazghani/terraformer/internal/httpserver"
	"github.com/maazghani/terraformer/internal/patch"
	"github.com/maazghani/terraformer/internal/repo"
	"github.com/maazghani/terraformer/internal/runner"
	"github.com/maazghani/terraformer/internal/terraform"

	"bytes"
)

// newE2EServer builds a server with a fake terraform runner for end-to-end tests.
func newE2EServer(t *testing.T) (*httptest.Server, string) {
	t.Helper()
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "main.tf"), "# initial\n")

	repoSvc, err := repo.New(root)
	if err != nil {
		t.Fatalf("repo.New: %v", err)
	}

	fakeRunner := runner.NewFakeRunner()
	// Register terraform with exit 0 for all subcommands.
	fakeRunner.Register("terraform", runner.Result{Stdout: "", Stderr: "", ExitCode: 0}, nil)

	tfSvc, err := terraform.NewService(fakeRunner, root)
	if err != nil {
		t.Fatalf("terraform.NewService: %v", err)
	}

	patchSvc, err := patch.New(root)
	if err != nil {
		t.Fatalf("patch.New: %v", err)
	}

	var logBuf bytes.Buffer
	srv := httpserver.New(httpserver.Config{RepoRoot: root, Port: 0}, repoSvc, tfSvc, patchSvc, &logBuf)
	ts := httptest.NewServer(srv)
	t.Cleanup(ts.Close)
	return ts, root
}

// TestEndToEnd_FullFlow exercises the complete v0 agent workflow:
// list → read → patch → fmt → validate → plan → check_desired_state.
func TestEndToEnd_FullFlow(t *testing.T) {
	ts, root := newE2EServer(t)

	// Step 1: list_repo_files — must see main.tf.
	{
		got, resp := doPost(t, ts.URL+"/tools/list_repo_files", `{"path":".","max_files":100}`)
		if resp.StatusCode != 200 {
			t.Fatalf("list_repo_files: expected 200, got %d", resp.StatusCode)
		}
		if v, _ := got["ok"].(bool); !v {
			t.Fatalf("list_repo_files: ok=false, got %v", got)
		}
		files, _ := got["files"].([]interface{})
		found := false
		for _, f := range files {
			fm, _ := f.(map[string]interface{})
			if fm["path"] == "main.tf" {
				found = true
			}
		}
		if !found {
			t.Error("list_repo_files: expected main.tf in listing")
		}
	}

	// Step 2: read_repo_file — must return content of main.tf.
	{
		got, resp := doPost(t, ts.URL+"/tools/read_repo_file", `{"path":"main.tf"}`)
		if resp.StatusCode != 200 {
			t.Fatalf("read_repo_file: expected 200, got %d", resp.StatusCode)
		}
		if v, _ := got["ok"].(bool); !v {
			t.Fatalf("read_repo_file: ok=false, got %v", got)
		}
		if got["content"] == "" {
			t.Error("read_repo_file: expected non-empty content")
		}
	}

	// Step 3: apply_patch — overwrite main.tf with new content.
	{
		body := `{"files":[{"path":"main.tf","operation":"write","content":"resource \"local_file\" \"x\" {}\n"}]}`
		got, resp := doPost(t, ts.URL+"/tools/apply_patch", body)
		if resp.StatusCode != 200 {
			t.Fatalf("apply_patch: expected 200, got %d", resp.StatusCode)
		}
		if v, _ := got["ok"].(bool); !v {
			t.Fatalf("apply_patch: ok=false, got %v", got)
		}
		// Verify file was actually changed on disk.
		content, err := os.ReadFile(filepath.Join(root, "main.tf"))
		if err != nil {
			t.Fatalf("apply_patch: reading file after patch: %v", err)
		}
		if string(content) != "resource \"local_file\" \"x\" {}\n" {
			t.Errorf("apply_patch: unexpected file content: %q", string(content))
		}
	}

	// Step 4: terraform_fmt — run format (fake runner returns exit 0).
	{
		got, resp := doPost(t, ts.URL+"/tools/terraform_fmt", `{"check":false,"recursive":true}`)
		if resp.StatusCode != 200 {
			t.Fatalf("terraform_fmt: expected 200, got %d", resp.StatusCode)
		}
		if _, ok := got["ok"]; !ok {
			t.Error("terraform_fmt: response missing 'ok' field")
		}
	}

	// Step 5: terraform_validate — run validate (fake runner returns exit 0).
	{
		got, resp := doPost(t, ts.URL+"/tools/terraform_validate", `{"json":true}`)
		if resp.StatusCode != 200 {
			t.Fatalf("terraform_validate: expected 200, got %d", resp.StatusCode)
		}
		if _, ok := got["ok"]; !ok {
			t.Error("terraform_validate: response missing 'ok' field")
		}
	}

	// Step 6: terraform_plan — desired_state_status must always be "not_checked".
	{
		got, resp := doPost(t, ts.URL+"/tools/terraform_plan", `{"detailed_exitcode":true}`)
		if resp.StatusCode != 200 {
			t.Fatalf("terraform_plan: expected 200, got %d", resp.StatusCode)
		}
		if ds, ok := got["desired_state_status"]; !ok {
			t.Error("terraform_plan: response missing 'desired_state_status'")
		} else if ds != "not_checked" {
			t.Errorf("terraform_plan: desired_state_status=%q, want 'not_checked'", ds)
		}
	}

	// Step 7: check_desired_state — status must be "not_implemented" in v0.
	// We first write a dummy plan.json file since the path must exist.
	planPath := filepath.Join(root, "plan.json")
	if err := os.WriteFile(planPath, []byte(`{"format_version":"1.0"}`), 0o644); err != nil {
		t.Fatalf("write plan.json: %v", err)
	}
	{
		body := `{"desired_state":{"resources":[]},"plan_json_path":"plan.json"}`
		got, resp := doPost(t, ts.URL+"/tools/check_desired_state", body)
		if resp.StatusCode != 200 {
			t.Fatalf("check_desired_state: expected 200, got %d", resp.StatusCode)
		}
		if v, _ := got["ok"].(bool); !v {
			t.Fatalf("check_desired_state: ok=false, got %v", got)
		}
		if got["status"] != "not_implemented" {
			t.Errorf("check_desired_state: status=%q, want 'not_implemented'", got["status"])
		}
	}
}

// TestEndToEnd_TraversalRejected verifies that path traversal is blocked at
// the HTTP boundary for all file-path-accepting tools.
func TestEndToEnd_TraversalRejected(t *testing.T) {
	ts, _ := newE2EServer(t)

	cases := []struct {
		tool string
		body string
	}{
		{"/tools/read_repo_file", `{"path":"../../etc/passwd"}`},
		{"/tools/apply_patch", `{"files":[{"path":"../../etc/passwd","operation":"write","content":"x"}]}`},
		{"/tools/list_repo_files", `{"path":"../../etc"}`},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.tool, func(t *testing.T) {
			got, resp := doPost(t, ts.URL+tc.tool, tc.body)
			if resp.StatusCode == 200 {
				// The tool may return 200 with ok=false — that is acceptable.
				if v, _ := got["ok"].(bool); v {
					t.Errorf("%s: traversal path was accepted (ok=true)", tc.tool)
				}
			}
		})
	}
}
