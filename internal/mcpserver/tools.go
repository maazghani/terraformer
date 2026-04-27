package mcpserver

import (
	"encoding/json"
	"fmt"

	"github.com/maazghani/terraformer/internal/patch"
	"github.com/maazghani/terraformer/internal/repo"
	"github.com/maazghani/terraformer/internal/terraform"
	"github.com/maazghani/terraformer/internal/tools"
)

// ToolDef describes a single MCP tool with its name and JSON Schema input schema.
type ToolDef struct {
	Name        string          `json:"name"`
	InputSchema json.RawMessage `json:"inputSchema"`
}

// AllTools returns all 9 v0 tool definitions with their JSON Schema inputSchema.
// The schemas match the contracts defined in PLAN/spec/02-mcp-tool-contracts.md.
func AllTools() []ToolDef {
	return []ToolDef{
		{
			Name: "terraform_init",
			InputSchema: rawJSON(`{"type":"object","properties":{"upgrade":{"type":"boolean","default":false},"backend":{"type":"boolean","default":true}}}`),
		},
		{
			Name: "terraform_fmt",
			InputSchema: rawJSON(`{"type":"object","properties":{"check":{"type":"boolean","default":false},"recursive":{"type":"boolean","default":true}}}`),
		},
		{
			Name: "terraform_validate",
			InputSchema: rawJSON(`{"type":"object","properties":{"json":{"type":"boolean","default":true}}}`),
		},
		{
			Name: "terraform_plan",
			InputSchema: rawJSON(`{"type":"object","properties":{"out":{"type":"string","default":".terraformer/plan.tfplan"},"detailed_exitcode":{"type":"boolean","default":true},"refresh":{"type":"boolean","default":false}}}`),
		},
		{
			Name: "terraform_show_json",
			InputSchema: rawJSON(`{"type":"object","properties":{"plan_path":{"type":"string","default":".terraformer/plan.tfplan"}},"required":["plan_path"]}`),
		},
		{
			Name: "list_repo_files",
			InputSchema: rawJSON(`{"type":"object","properties":{"path":{"type":"string","default":"."},"include_globs":{"type":"array","items":{"type":"string"}},"exclude_globs":{"type":"array","items":{"type":"string"}},"max_files":{"type":"integer","default":200}}}`),
		},
		{
			Name: "read_repo_file",
			InputSchema: rawJSON(`{"type":"object","properties":{"path":{"type":"string","default":"main.tf"},"max_bytes":{"type":"integer","default":65536}},"required":["path"]}`),
		},
		{
			Name: "apply_patch",
			InputSchema: rawJSON(`{"type":"object","properties":{"files":{"type":"array","items":{"type":"object","properties":{"path":{"type":"string"},"operation":{"type":"string","enum":["write","delete"]},"content":{"type":"string"}},"required":["path","operation"]}}},"required":["files"]}`),
		},
		{
			Name: "check_desired_state",
			InputSchema: rawJSON(`{"type":"object","properties":{"desired_state":{"type":"object","properties":{"resources":{"type":"array"}}},"plan_json_path":{"type":"string"}}}`),
		},
	}
}

// applyPatchArgs mirrors the JSON shape expected for apply_patch arguments.
type applyPatchArgs struct {
	Files []struct {
		Path      string `json:"path"`
		Operation string `json:"operation"`
		Content   string `json:"content,omitempty"`
	} `json:"files"`
}

// applyPatchResp is the structured response returned for apply_patch tool calls.
type applyPatchResp struct {
	OK            bool     `json:"ok"`
	ChangedFiles  []string `json:"changed_files"`
	RejectedFiles []string `json:"rejected_files"`
	Warnings      []string `json:"warnings"`
}

// CallTool dispatches a tools/call request to the appropriate internal handler
// and returns the JSON-serialized response as a string. It returns an error for
// unknown tool names or invalid arguments. Forbidden tool names (apply, destroy,
// arbitrary shell) are not in the allowlist and therefore return errors.
func CallTool(name string, args json.RawMessage, repoSvc *repo.Service, tfSvc *terraform.Service, patchSvc *patch.Service, repoRoot string) (string, error) {
	if args == nil {
		args = json.RawMessage(`{}`)
	}
	switch name {
	case "terraform_init":
		var req terraform.InitRequest
		if err := json.Unmarshal(args, &req); err != nil {
			return "", fmt.Errorf("invalid arguments for terraform_init: %w", err)
		}
		return jsonString(tfSvc.Init(req))

	case "terraform_fmt":
		var req terraform.FmtRequest
		if err := json.Unmarshal(args, &req); err != nil {
			return "", fmt.Errorf("invalid arguments for terraform_fmt: %w", err)
		}
		return jsonString(tfSvc.Fmt(req))

	case "terraform_validate":
		var req terraform.ValidateRequest
		if err := json.Unmarshal(args, &req); err != nil {
			return "", fmt.Errorf("invalid arguments for terraform_validate: %w", err)
		}
		return jsonString(tfSvc.Validate(req))

	case "terraform_plan":
		var req terraform.PlanRequest
		if err := json.Unmarshal(args, &req); err != nil {
			return "", fmt.Errorf("invalid arguments for terraform_plan: %w", err)
		}
		return jsonString(tfSvc.Plan(req))

	case "terraform_show_json":
		var req terraform.ShowJSONRequest
		if err := json.Unmarshal(args, &req); err != nil {
			return "", fmt.Errorf("invalid arguments for terraform_show_json: %w", err)
		}
		return jsonString(tfSvc.ShowJSON(req))

	case "list_repo_files":
		var req tools.ListRepoFilesRequest
		if err := json.Unmarshal(args, &req); err != nil {
			return "", fmt.Errorf("invalid arguments for list_repo_files: %w", err)
		}
		return jsonString(tools.ListRepoFiles(repoSvc, req))

	case "read_repo_file":
		var req tools.ReadRepoFileRequest
		if err := json.Unmarshal(args, &req); err != nil {
			return "", fmt.Errorf("invalid arguments for read_repo_file: %w", err)
		}
		return jsonString(tools.ReadRepoFile(repoSvc, req))

	case "apply_patch":
		var pArgs applyPatchArgs
		if err := json.Unmarshal(args, &pArgs); err != nil {
			return "", fmt.Errorf("invalid arguments for apply_patch: %w", err)
		}
		ops := make([]patch.FileOperation, len(pArgs.Files))
		for i, f := range pArgs.Files {
			ops[i] = patch.FileOperation{
				Path:      f.Path,
				Operation: f.Operation,
				Content:   f.Content,
			}
		}
		result, err := patchSvc.ApplyPatch(patch.ApplyPatchRequest{Files: ops})
		if err != nil {
			return "", fmt.Errorf("apply_patch: %w", err)
		}
		resp := applyPatchResp{
			OK:            result.OK,
			ChangedFiles:  result.ChangedFiles,
			RejectedFiles: result.RejectedFiles,
			Warnings:      result.Warnings,
		}
		return jsonString(resp)

	case "check_desired_state":
		var req tools.CheckDesiredStateRequest
		if err := json.Unmarshal(args, &req); err != nil {
			return "", fmt.Errorf("invalid arguments for check_desired_state: %w", err)
		}
		return jsonString(tools.CheckDesiredState(repoRoot, req))

	default:
		return "", fmt.Errorf("unknown tool: %q", name)
	}
}

// jsonString marshals v to a JSON string.
func jsonString(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("json.Marshal: %w", err)
	}
	return string(b), nil
}

// rawJSON returns s as a json.RawMessage. s must be valid JSON.
func rawJSON(s string) json.RawMessage {
	return json.RawMessage(s)
}
