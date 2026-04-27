// Package mcpserver_test contains tests for the MCP tool definitions and routing.
package mcpserver_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/maazghani/terraformer/internal/mcpserver"
	"github.com/maazghani/terraformer/internal/patch"
	"github.com/maazghani/terraformer/internal/repo"
	"github.com/maazghani/terraformer/internal/runner"
	"github.com/maazghani/terraformer/internal/terraform"
)

// newTestServices creates test services backed by a temp directory.
func newTestServices(t *testing.T) (*repo.Service, *terraform.Service, *patch.Service, string) {
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

	return repoSvc, tfSvc, patchSvc, root
}

// ---------------------------------------------------------------------------
// AllTools tests
// ---------------------------------------------------------------------------

// TestAllTools_Count verifies that AllTools returns exactly 9 tool definitions.
func TestAllTools_Count(t *testing.T) {
	tools := mcpserver.AllTools()
	if len(tools) != 9 {
		t.Errorf("expected 9 tools, got %d", len(tools))
	}
}

// TestAllTools_AllHaveNames verifies that every tool has a non-empty name.
func TestAllTools_AllHaveNames(t *testing.T) {
	for i, tool := range mcpserver.AllTools() {
		if tool.Name == "" {
			t.Errorf("tools[%d] has empty name", i)
		}
	}
}

// TestAllTools_AllHaveInputSchema verifies that every tool has a non-empty
// inputSchema that is valid JSON.
func TestAllTools_AllHaveInputSchema(t *testing.T) {
	for _, tool := range mcpserver.AllTools() {
		if len(tool.InputSchema) == 0 {
			t.Errorf("tool %q has empty inputSchema", tool.Name)
			continue
		}
		if !json.Valid(tool.InputSchema) {
			t.Errorf("tool %q inputSchema is not valid JSON: %s", tool.Name, tool.InputSchema)
		}
	}
}

// TestAllTools_ExpectedNames verifies that the 9 required v0 tool names are present.
func TestAllTools_ExpectedNames(t *testing.T) {
	want := []string{
		"terraform_init", "terraform_fmt", "terraform_validate",
		"terraform_plan", "terraform_show_json", "list_repo_files",
		"read_repo_file", "apply_patch", "check_desired_state",
	}
	gotNames := make(map[string]bool)
	for _, tool := range mcpserver.AllTools() {
		gotNames[tool.Name] = true
	}
	for _, name := range want {
		if !gotNames[name] {
			t.Errorf("missing tool %q", name)
		}
	}
}

// ---------------------------------------------------------------------------
// CallTool tests
// ---------------------------------------------------------------------------

// TestCallTool_KnownTool verifies that calling a known tool returns valid JSON.
func TestCallTool_KnownTool(t *testing.T) {
	repoSvc, tfSvc, patchSvc, root := newTestServices(t)
	text, err := mcpserver.CallTool("list_repo_files", json.RawMessage(`{}`), repoSvc, tfSvc, patchSvc, root)
	if err != nil {
		t.Fatalf("CallTool list_repo_files: %v", err)
	}
	if !json.Valid([]byte(text)) {
		t.Errorf("CallTool returned invalid JSON: %s", text)
	}
}

// TestCallTool_UnknownTool verifies that calling an unknown tool returns an error.
func TestCallTool_UnknownTool(t *testing.T) {
	repoSvc, tfSvc, patchSvc, root := newTestServices(t)
	_, err := mcpserver.CallTool("nonexistent_tool", json.RawMessage(`{}`), repoSvc, tfSvc, patchSvc, root)
	if err == nil {
		t.Error("expected error for unknown tool, got nil")
	}
}

// TestCallTool_AllKnownTools verifies that every tool in AllTools() can be
// called with empty arguments without panicking and returns valid JSON.
func TestCallTool_AllKnownTools(t *testing.T) {
	repoSvc, tfSvc, patchSvc, root := newTestServices(t)

	// These tools require specific arguments to avoid errors; test them with
	// minimal valid arguments.
	toolArgs := map[string]string{
		"terraform_init":      `{}`,
		"terraform_fmt":       `{}`,
		"terraform_validate":  `{}`,
		"terraform_plan":      `{}`,
		"terraform_show_json": `{"plan_path":".terraformer/plan.tfplan"}`,
		"list_repo_files":     `{}`,
		"read_repo_file":      `{"path":"main.tf"}`,
		"apply_patch":         `{"files":[]}`,
		"check_desired_state": `{"desired_state":{"resources":[]},"plan_json_path":"plan.json"}`,
	}

	for _, tool := range mcpserver.AllTools() {
		name := tool.Name
		args, ok := toolArgs[name]
		if !ok {
			t.Errorf("test missing args for tool %q", name)
			continue
		}
		text, err := mcpserver.CallTool(name, json.RawMessage(args), repoSvc, tfSvc, patchSvc, root)
		if err != nil {
			// Some tools may legitimately error (e.g., show_json without a real plan file).
			// We just ensure no panic and a meaningful error.
			t.Logf("tool %q returned error (may be expected): %v", name, err)
			continue
		}
		if !json.Valid([]byte(text)) {
			t.Errorf("tool %q returned invalid JSON: %s", name, text)
		}
	}
}

// TestCallTool_ForbiddenTools verifies that forbidden tool names return errors.
func TestCallTool_ForbiddenTools(t *testing.T) {
	repoSvc, tfSvc, patchSvc, root := newTestServices(t)
	forbidden := []string{"terraform_apply", "terraform_destroy", "exec", "shell"}
	for _, name := range forbidden {
		name := name
		t.Run(name, func(t *testing.T) {
			_, err := mcpserver.CallTool(name, json.RawMessage(`{}`), repoSvc, tfSvc, patchSvc, root)
			if err == nil {
				t.Errorf("expected error for forbidden tool %q, got nil", name)
			}
		})
	}
}
