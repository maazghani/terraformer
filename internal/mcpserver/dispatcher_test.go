// Package mcpserver_test contains tests for the MCP JSON-RPC 2.0 dispatcher.
package mcpserver_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/maazghani/terraformer/internal/mcpserver"
	"github.com/maazghani/terraformer/internal/patch"
	"github.com/maazghani/terraformer/internal/repo"
	"github.com/maazghani/terraformer/internal/runner"
	"github.com/maazghani/terraformer/internal/terraform"
)

// newTestDispatcher creates a Dispatcher backed by temporary services for testing.
func newTestDispatcher(t *testing.T) *mcpserver.Dispatcher {
	t.Helper()
	root := t.TempDir()

	if err := os.WriteFile(filepath.Join(root, "main.tf"), []byte("# test\n"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

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

	return mcpserver.New(repoSvc, tfSvc, patchSvc, root)
}

// doRPC posts a JSON-RPC body directly to the dispatcher and decodes the response map.
func doRPC(t *testing.T, d *mcpserver.Dispatcher, body string) (map[string]interface{}, *http.Response) {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	d.ServeHTTP(w, req)
	resp := w.Result()
	t.Cleanup(func() { resp.Body.Close() })

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return result, resp
}

// ---------------------------------------------------------------------------
// TDD task 1: parse valid JSON-RPC request / malformed JSON → -32700
// ---------------------------------------------------------------------------

// TestDispatcher_ParseValidRequest verifies that a valid JSON-RPC body is
// accepted and produces a result (not an error).
func TestDispatcher_ParseValidRequest(t *testing.T) {
	d := newTestDispatcher(t)
	got, resp := doRPC(t, d, `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}`)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if got["result"] == nil {
		t.Errorf("expected non-nil result for valid request, got: %v", got)
	}
	if got["error"] != nil {
		t.Errorf("expected no error for valid request, got: %v", got["error"])
	}
}

// TestDispatcher_ParseError_MalformedJSON verifies that malformed JSON returns
// a JSON-RPC error with code -32700.
func TestDispatcher_ParseError_MalformedJSON(t *testing.T) {
	d := newTestDispatcher(t)
	got, resp := doRPC(t, d, `not-json{{{`)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 even on parse error, got %d", resp.StatusCode)
	}
	errObj, ok := got["error"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected error object, got: %v", got)
	}
	code, _ := errObj["code"].(float64)
	if code != -32700 {
		t.Errorf("expected error code -32700, got %v", code)
	}
}

// ---------------------------------------------------------------------------
// TDD task 2: initialize handshake
// ---------------------------------------------------------------------------

// TestDispatcher_Initialize verifies the initialize response shape.
func TestDispatcher_Initialize(t *testing.T) {
	d := newTestDispatcher(t)
	got, _ := doRPC(t, d, `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"codex","version":"1.0"}}}`)

	result, ok := got["result"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected result object, got: %v", got)
	}

	if v, _ := result["protocolVersion"].(string); v != "2024-11-05" {
		t.Errorf("expected protocolVersion 2024-11-05, got %q", v)
	}

	caps, _ := result["capabilities"].(map[string]interface{})
	if caps == nil {
		t.Error("expected capabilities object")
	}
	if _, ok := caps["tools"]; !ok {
		t.Error("expected capabilities.tools field")
	}

	srvInfo, _ := result["serverInfo"].(map[string]interface{})
	if name, _ := srvInfo["name"].(string); name != "terraformer" {
		t.Errorf("expected serverInfo.name=terraformer, got %q", name)
	}
}

// ---------------------------------------------------------------------------
// TDD task 3: tools/list — all 9 tools with non-empty inputSchema
// ---------------------------------------------------------------------------

// TestDispatcher_ToolsList verifies that tools/list returns all 9 v0 tools.
func TestDispatcher_ToolsList(t *testing.T) {
	d := newTestDispatcher(t)
	got, _ := doRPC(t, d, `{"jsonrpc":"2.0","id":2,"method":"tools/list"}`)

	result, ok := got["result"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected result object, got: %v", got)
	}

	toolsArr, ok := result["tools"].([]interface{})
	if !ok {
		t.Fatalf("expected tools array, got %T: %v", result["tools"], result["tools"])
	}

	if len(toolsArr) != 9 {
		t.Errorf("expected 9 tools, got %d", len(toolsArr))
	}

	wantNames := []string{
		"terraform_init", "terraform_fmt", "terraform_validate",
		"terraform_plan", "terraform_show_json", "list_repo_files",
		"read_repo_file", "apply_patch", "check_desired_state",
	}
	gotNames := make(map[string]bool)
	for _, item := range toolsArr {
		toolMap, _ := item.(map[string]interface{})
		name, _ := toolMap["name"].(string)
		gotNames[name] = true
		if toolMap["inputSchema"] == nil {
			t.Errorf("tool %q has nil inputSchema", name)
		}
	}
	for _, want := range wantNames {
		if !gotNames[want] {
			t.Errorf("missing tool %q in tools/list response", want)
		}
	}
}

// ---------------------------------------------------------------------------
// TDD task 4: tools/call dispatch — known tool, content[0].type == "text"
// ---------------------------------------------------------------------------

// TestDispatcher_ToolsCall_Dispatch verifies that a known tool call returns
// content[0].type == "text" with valid JSON in content[0].text.
func TestDispatcher_ToolsCall_Dispatch(t *testing.T) {
	d := newTestDispatcher(t)
	body := `{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"list_repo_files","arguments":{}}}`
	got, _ := doRPC(t, d, body)

	result, ok := got["result"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected result object, got: %v", got)
	}

	content, ok := result["content"].([]interface{})
	if !ok || len(content) == 0 {
		t.Fatalf("expected non-empty content array, got: %v", result["content"])
	}

	item, _ := content[0].(map[string]interface{})
	if typ, _ := item["type"].(string); typ != "text" {
		t.Errorf("expected content[0].type=text, got %q", typ)
	}

	text, _ := item["text"].(string)
	if !json.Valid([]byte(text)) {
		t.Errorf("content[0].text is not valid JSON: %s", text)
	}
}

// ---------------------------------------------------------------------------
// TDD task 5: tools/call unknown tool → -32602
// ---------------------------------------------------------------------------

// TestDispatcher_ToolsCall_UnknownTool verifies that calling an unknown tool
// returns JSON-RPC error code -32602.
func TestDispatcher_ToolsCall_UnknownTool(t *testing.T) {
	d := newTestDispatcher(t)
	body := `{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"nonexistent_tool","arguments":{}}}`
	got, _ := doRPC(t, d, body)

	errObj, ok := got["error"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected error object, got: %v", got)
	}
	code, _ := errObj["code"].(float64)
	if code != -32602 {
		t.Errorf("expected error code -32602, got %v", code)
	}
}

// ---------------------------------------------------------------------------
// TDD task 6: unknown method → -32601
// ---------------------------------------------------------------------------

// TestDispatcher_UnknownMethod verifies that an unknown method returns -32601.
func TestDispatcher_UnknownMethod(t *testing.T) {
	d := newTestDispatcher(t)
	body := `{"jsonrpc":"2.0","id":5,"method":"some/unknown/method"}`
	got, _ := doRPC(t, d, body)

	errObj, ok := got["error"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected error object, got: %v", got)
	}
	code, _ := errObj["code"].(float64)
	if code != -32601 {
		t.Errorf("expected error code -32601, got %v", code)
	}
}

// ---------------------------------------------------------------------------
// TDD task 7: Content-Type: application/json on all responses
// ---------------------------------------------------------------------------

// TestDispatcher_ContentType verifies that every response has
// Content-Type: application/json, regardless of success or error.
func TestDispatcher_ContentType(t *testing.T) {
	d := newTestDispatcher(t)

	cases := []string{
		`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}`,
		`{"jsonrpc":"2.0","id":2,"method":"tools/list"}`,
		`{"jsonrpc":"2.0","id":3,"method":"unknown"}`,
		`not-json`,
	}

	for _, body := range cases {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		w := httptest.NewRecorder()
		d.ServeHTTP(w, req)
		ct := w.Header().Get("Content-Type")
		if !strings.Contains(ct, "application/json") {
			t.Errorf("body %q: expected Content-Type application/json, got %q", body, ct)
		}
	}
}

// ---------------------------------------------------------------------------
// TDD task 8: Integration — Codex startup sequence
// ---------------------------------------------------------------------------

// TestDispatcher_Integration_CodexStartupSequence exercises the full
// initialize → tools/list → tools/call sequence through an httptest.Server.
func TestDispatcher_Integration_CodexStartupSequence(t *testing.T) {
	d := newTestDispatcher(t)
	ts := httptest.NewServer(d)
	t.Cleanup(ts.Close)

	post := func(body string) map[string]interface{} {
		t.Helper()
		resp, err := http.Post(ts.URL, "application/json", strings.NewReader(body))
		if err != nil {
			t.Fatalf("http.Post: %v", err)
		}
		defer resp.Body.Close()
		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("decode: %v", err)
		}
		return result
	}

	// Step 1: initialize
	init1 := post(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"codex","version":"1.0"}}}`)
	if init1["result"] == nil {
		t.Fatalf("initialize: expected result, got nil; response: %v", init1)
	}

	// Step 2: tools/list
	list := post(`{"jsonrpc":"2.0","id":2,"method":"tools/list"}`)
	listResult, _ := list["result"].(map[string]interface{})
	toolsArr, _ := listResult["tools"].([]interface{})
	if len(toolsArr) < 9 {
		t.Fatalf("tools/list: expected >= 9 tools, got %d", len(toolsArr))
	}

	// Step 3: tools/call
	call := post(`{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"list_repo_files","arguments":{}}}`)
	callResult, _ := call["result"].(map[string]interface{})
	content, _ := callResult["content"].([]interface{})
	if len(content) == 0 {
		t.Fatalf("tools/call: expected non-empty content, got: %v", call)
	}
	item, _ := content[0].(map[string]interface{})
	if typ, _ := item["type"].(string); typ != "text" {
		t.Errorf("tools/call: expected content[0].type=text, got %q", typ)
	}
}

// ---------------------------------------------------------------------------
// Safety: forbidden tool names must not be callable
// ---------------------------------------------------------------------------

// TestDispatcher_ToolsCall_ForbiddenTools verifies that tools/call rejects
// names like terraform_apply, terraform_destroy, or shell execution.
func TestDispatcher_ToolsCall_ForbiddenTools(t *testing.T) {
	d := newTestDispatcher(t)
	forbidden := []string{"terraform_apply", "terraform_destroy", "exec", "shell", "run"}
	for _, name := range forbidden {
		name := name
		t.Run(name, func(t *testing.T) {
			body := `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"` + name + `","arguments":{}}}`
			got, _ := doRPC(t, d, body)
			if got["error"] == nil {
				t.Errorf("expected error for forbidden tool %q, got result: %v", name, got)
			}
		})
	}
}
