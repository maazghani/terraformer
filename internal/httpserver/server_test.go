// Package httpserver_test contains tests for the HTTP/JSON server wiring.
package httpserver_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/maazghani/terraformer/internal/httpserver"
	"github.com/maazghani/terraformer/internal/patch"
	"github.com/maazghani/terraformer/internal/repo"
	"github.com/maazghani/terraformer/internal/runner"
	"github.com/maazghani/terraformer/internal/terraform"
)

// writeFile is a test helper that creates a file with given content.
func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
}

// newTestServer builds a Server backed by a temp repo with a fake runner and
// returns the httptest.Server and a buffer capturing log output.
func newTestServer(t *testing.T) (*httptest.Server, *bytes.Buffer, string) {
	t.Helper()
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "main.tf"), "# main\n")

	repoSvc, err := repo.New(root)
	if err != nil {
		t.Fatalf("repo.New: %v", err)
	}

	fakeRunner := runner.NewFakeRunner()
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
	cfg := httpserver.Config{RepoRoot: root, Port: 0}
	srv := httpserver.New(cfg, repoSvc, tfSvc, patchSvc, &logBuf)

	ts := httptest.NewServer(srv)
	t.Cleanup(ts.Close)
	return ts, &logBuf, root
}

// doPost is a helper that POSTs JSON and returns the decoded response map.
func doPost(t *testing.T, url, body string) (map[string]interface{}, *http.Response) {
	t.Helper()
	resp, err := http.Post(url, "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST %s: %v", url, err)
	}
	t.Cleanup(func() { resp.Body.Close() })
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode response from %s: %v", url, err)
	}
	return result, resp
}

// ---------------------------------------------------------------------------
// Tool registration
// ---------------------------------------------------------------------------

// TestServer_AllToolsRegistered verifies all 9 required v0 tool endpoints exist
// and do not return 404.
func TestServer_AllToolsRegistered(t *testing.T) {
	ts, _, _ := newTestServer(t)

	endpoints := []string{
		"/tools/terraform_init",
		"/tools/terraform_fmt",
		"/tools/terraform_validate",
		"/tools/terraform_plan",
		"/tools/terraform_show_json",
		"/tools/list_repo_files",
		"/tools/read_repo_file",
		"/tools/apply_patch",
		"/tools/check_desired_state",
	}

	for _, ep := range endpoints {
		ep := ep
		t.Run(ep, func(t *testing.T) {
			resp, err := http.Post(ts.URL+ep, "application/json", bytes.NewBufferString("{}"))
			if err != nil {
				t.Fatalf("POST %s: %v", ep, err)
			}
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusNotFound {
				t.Errorf("endpoint %s not registered (got 404)", ep)
			}
		})
	}
}

// TestServer_OnlyPOSTAccepted verifies GET requests are rejected (405 or 404).
func TestServer_OnlyPOSTAccepted(t *testing.T) {
	ts, _, _ := newTestServer(t)

	resp, err := http.Get(ts.URL + "/tools/list_repo_files")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		t.Error("expected non-200 for GET on a POST-only endpoint")
	}
}

// ---------------------------------------------------------------------------
// Request decoding
// ---------------------------------------------------------------------------

// TestServer_RequestDecoding_ListRepoFiles verifies the list_repo_files endpoint
// decodes the JSON request and returns a structured response.
func TestServer_RequestDecoding_ListRepoFiles(t *testing.T) {
	ts, _, _ := newTestServer(t)

	got, resp := doPost(t, ts.URL+"/tools/list_repo_files", `{"path":".","max_files":10}`)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if _, ok := got["ok"]; !ok {
		t.Error("response missing 'ok' field")
	}
	if _, ok := got["files"]; !ok {
		t.Error("response missing 'files' field")
	}
}

// TestServer_RequestDecoding_ReadRepoFile verifies the read_repo_file endpoint.
func TestServer_RequestDecoding_ReadRepoFile(t *testing.T) {
	ts, _, _ := newTestServer(t)

	got, resp := doPost(t, ts.URL+"/tools/read_repo_file", `{"path":"main.tf"}`)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if v, _ := got["ok"].(bool); !v {
		t.Errorf("expected ok=true, got %v", got)
	}
	if _, ok := got["content"]; !ok {
		t.Error("response missing 'content' field")
	}
}

// TestServer_RequestDecoding_TerraformInit verifies terraform_init decodes and
// dispatches to the service.
func TestServer_RequestDecoding_TerraformInit(t *testing.T) {
	ts, _, _ := newTestServer(t)

	got, resp := doPost(t, ts.URL+"/tools/terraform_init", `{"upgrade":false}`)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if _, ok := got["ok"]; !ok {
		t.Error("response missing 'ok' field")
	}
	if _, ok := got["command"]; !ok {
		t.Error("response missing 'command' field")
	}
}

// ---------------------------------------------------------------------------
// Response encoding
// ---------------------------------------------------------------------------

// TestServer_ResponseContentType verifies every tool response uses
// Content-Type: application/json.
func TestServer_ResponseContentType(t *testing.T) {
	ts, _, _ := newTestServer(t)

	endpoints := []struct {
		path string
		body string
	}{
		{"/tools/list_repo_files", `{}`},
		{"/tools/read_repo_file", `{"path":"main.tf"}`},
		{"/tools/terraform_init", `{}`},
		{"/tools/terraform_fmt", `{}`},
		{"/tools/terraform_validate", `{}`},
		{"/tools/terraform_plan", `{}`},
	}

	for _, tc := range endpoints {
		tc := tc
		t.Run(tc.path, func(t *testing.T) {
			resp, err := http.Post(ts.URL+tc.path, "application/json", strings.NewReader(tc.body))
			if err != nil {
				t.Fatalf("POST: %v", err)
			}
			defer resp.Body.Close()
			ct := resp.Header.Get("Content-Type")
			if !strings.Contains(ct, "application/json") {
				t.Errorf("expected Content-Type application/json, got %q", ct)
			}
		})
	}
}

// TestServer_ResponseIsValidJSON verifies every tool produces valid JSON output.
func TestServer_ResponseIsValidJSON(t *testing.T) {
	ts, _, _ := newTestServer(t)

	endpoints := []struct {
		path string
		body string
	}{
		{"/tools/list_repo_files", `{}`},
		{"/tools/read_repo_file", `{"path":"main.tf"}`},
		{"/tools/apply_patch", `{"files":[]}`},
		{"/tools/check_desired_state", `{"desired_state":{"resources":[]},"plan_json_path":"plan.json"}`},
	}

	for _, tc := range endpoints {
		tc := tc
		t.Run(tc.path, func(t *testing.T) {
			resp, err := http.Post(ts.URL+tc.path, "application/json", strings.NewReader(tc.body))
			if err != nil {
				t.Fatalf("POST: %v", err)
			}
			defer resp.Body.Close()
			data, _ := io.ReadAll(resp.Body)
			if !json.Valid(data) {
				t.Errorf("response is not valid JSON: %s", data)
			}
		})
	}
}

// TestServer_TerraformResponseShape verifies terraform tool responses include
// the required fields (ok, command, stdout, stderr, exit_code, duration_ms,
// diagnostics, warnings).
func TestServer_TerraformResponseShape(t *testing.T) {
	ts, _, _ := newTestServer(t)

	got, resp := doPost(t, ts.URL+"/tools/terraform_validate", `{"json":false}`)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	required := []string{"ok", "command", "stdout", "stderr", "exit_code", "duration_ms", "diagnostics", "warnings"}
	for _, field := range required {
		if _, ok := got[field]; !ok {
			t.Errorf("terraform_validate response missing field %q", field)
		}
	}
}

// TestServer_RepoRootImmutable verifies that the repo root used by the server
// cannot be overridden via an HTTP request.
func TestServer_RepoRootImmutable(t *testing.T) {
	ts, _, root := newTestServer(t)

	// Even if a request body tries to supply a repo root it should be ignored.
	got, resp := doPost(t, ts.URL+"/tools/list_repo_files",
		`{"path":".","repo_root":"/etc"}`)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	// The response should reflect the configured root, not /etc.
	// We verify this by checking that ok=true (listing /etc would fail or
	// return different files).
	if v, _ := got["ok"].(bool); !v {
		t.Errorf("expected ok=true, got ok=%v; warnings: %v", got["ok"], got["warnings"])
	}

	// The files should be from the configured temp root, not /etc.
	files, _ := got["files"].([]interface{})
	for _, f := range files {
		fmap, _ := f.(map[string]interface{})
		if p, _ := fmap["path"].(string); strings.HasPrefix(p, "/") {
			t.Errorf("file path %q looks absolute; repo root may have been overridden", p)
		}
	}
	_ = root
}

// ---------------------------------------------------------------------------
// Invalid request handling
// ---------------------------------------------------------------------------

// TestServer_InvalidJSON returns 400 for malformed request bodies.
func TestServer_InvalidJSON(t *testing.T) {
	ts, _, _ := newTestServer(t)

	endpoints := []string{
		"/tools/terraform_init",
		"/tools/terraform_fmt",
		"/tools/terraform_validate",
		"/tools/terraform_plan",
		"/tools/terraform_show_json",
		"/tools/list_repo_files",
		"/tools/read_repo_file",
		"/tools/apply_patch",
		"/tools/check_desired_state",
	}

	for _, ep := range endpoints {
		ep := ep
		t.Run(ep, func(t *testing.T) {
			resp, err := http.Post(ts.URL+ep, "application/json", strings.NewReader("not-json{{{"))
			if err != nil {
				t.Fatalf("POST %s: %v", ep, err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusBadRequest {
				t.Errorf("expected 400 for bad JSON at %s, got %d", ep, resp.StatusCode)
			}
		})
	}
}

// TestServer_InvalidJSON_ResponseIsJSON verifies the 400 error body is JSON.
func TestServer_InvalidJSON_ResponseIsJSON(t *testing.T) {
	ts, _, _ := newTestServer(t)

	resp, err := http.Post(ts.URL+"/tools/list_repo_files", "application/json",
		strings.NewReader("bad{{{"))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	if !json.Valid(data) {
		t.Errorf("error response body is not valid JSON: %s", data)
	}
}

// ---------------------------------------------------------------------------
// Structured JSON logging
// ---------------------------------------------------------------------------

// TestServer_StructuredLogging verifies that each request emits at least one
// valid JSON log line to the configured writer.
func TestServer_StructuredLogging(t *testing.T) {
	ts, logBuf, _ := newTestServer(t)

	_, resp := doPost(t, ts.URL+"/tools/list_repo_files", `{}`)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status: %d", resp.StatusCode)
	}

	lines := bytes.Split(bytes.TrimSpace(logBuf.Bytes()), []byte("\n"))
	if len(lines) == 0 || (len(lines) == 1 && len(lines[0]) == 0) {
		t.Fatal("expected at least one log line, got none")
	}
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		var entry map[string]interface{}
		if err := json.Unmarshal(line, &entry); err != nil {
			t.Errorf("log line is not valid JSON: %q, err: %v", line, err)
		}
	}
}

// TestServer_StructuredLogging_HasRequiredFields verifies each log entry
// contains level, message, and method fields.
func TestServer_StructuredLogging_HasRequiredFields(t *testing.T) {
	ts, logBuf, _ := newTestServer(t)

	doPost(t, ts.URL+"/tools/list_repo_files", `{}`)

	lines := bytes.Split(bytes.TrimSpace(logBuf.Bytes()), []byte("\n"))
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		var entry map[string]interface{}
		if err := json.Unmarshal(line, &entry); err != nil {
			t.Fatalf("log line is not valid JSON: %v", err)
		}
		for _, field := range []string{"level", "message"} {
			if _, ok := entry[field]; !ok {
				t.Errorf("log entry missing field %q: %s", field, line)
			}
		}
	}
}

// TestServer_NoUnsafeCommandsReachable verifies that the httpserver does not
// expose terraform apply, terraform destroy, or arbitrary shell execution.
func TestServer_NoUnsafeCommandsReachable(t *testing.T) {
	ts, _, _ := newTestServer(t)

	forbidden := []string{
		"/tools/terraform_apply",
		"/tools/terraform_destroy",
		"/tools/exec",
		"/tools/shell",
		"/tools/run",
	}

	for _, ep := range forbidden {
		ep := ep
		t.Run(ep, func(t *testing.T) {
			resp, err := http.Post(ts.URL+ep, "application/json", bytes.NewBufferString("{}"))
			if err != nil {
				t.Fatalf("POST %s: %v", ep, err)
			}
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				t.Errorf("forbidden endpoint %s returned 200", ep)
			}
		})
	}
}

// TestServer_ResponseTruncation verifies that large responses are truncated
// when maxResponseBytes is configured.
func TestServer_ResponseTruncation(t *testing.T) {
	root := t.TempDir()
	// Create a large file that exceeds the max response bytes limit
	largeContent := strings.Repeat("x", 2000)
	writeFile(t, filepath.Join(root, "large.txt"), largeContent)

	repoSvc, err := repo.New(root)
	if err != nil {
		t.Fatalf("repo.New: %v", err)
	}

	fakeRunner := runner.NewFakeRunner()
	tfSvc, err := terraform.NewService(fakeRunner, root)
	if err != nil {
		t.Fatalf("terraform.NewService: %v", err)
	}

	patchSvc, err := patch.New(root)
	if err != nil {
		t.Fatalf("patch.New: %v", err)
	}

	var logBuf bytes.Buffer
	cfg := httpserver.Config{
		RepoRoot:         root,
		Port:             0,
		MaxResponseBytes: 100, // Set a small limit
	}
	srv := httpserver.New(cfg, repoSvc, tfSvc, patchSvc, &logBuf)

	ts := httptest.NewServer(srv)
	defer ts.Close()

	// Read the large file
	body := `{"path": "large.txt", "max_bytes": 10000}`
	result, _ := doPost(t, ts.URL+"/tools/read_repo_file", body)

	// Check that the response indicates truncation
	if truncated, ok := result["truncated"].(bool); !ok || !truncated {
		t.Errorf("expected truncated=true for large response, got %v", result["truncated"])
	}

	// Check that content is actually truncated
	content, ok := result["content"].(string)
	if !ok {
		t.Fatalf("content field missing or not a string")
	}
	if len(content) > 100 {
		t.Errorf("content length = %d, want <= 100", len(content))
	}
}
