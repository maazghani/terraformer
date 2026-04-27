// Package mcpserver implements the MCP Streamable HTTP transport protocol
// as a JSON-RPC 2.0 dispatcher. It handles the initialize handshake, tool
// discovery (tools/list), and tool invocation (tools/call) as required by
// MCP-compliant clients such as Codex.
//
// The dispatcher is mounted at POST / alongside the existing /tools/* endpoints
// in internal/httpserver so that MCP clients can discover and call all v0 tools.
package mcpserver

import (
	"encoding/json"
	"net/http"

	"github.com/maazghani/terraformer/internal/patch"
	"github.com/maazghani/terraformer/internal/repo"
	"github.com/maazghani/terraformer/internal/terraform"
)

// Standard JSON-RPC 2.0 error codes.
const (
	errCodeParseError     = -32700
	errCodeMethodNotFound = -32601
	errCodeInvalidParams  = -32602
)

// JSONRPCRequest is a JSON-RPC 2.0 request object.
type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// JSONRPCResponse is a JSON-RPC 2.0 response object.
type JSONRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Result  interface{}     `json:"result,omitempty"`
	Error   *JSONRPCError   `json:"error,omitempty"`
}

// JSONRPCError is a JSON-RPC 2.0 error object.
type JSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Dispatcher handles JSON-RPC 2.0 requests for the MCP Streamable HTTP transport.
// It is safe to use concurrently; all state is read-only after construction.
type Dispatcher struct {
	repoSvc  *repo.Service
	tfSvc    *terraform.Service
	patchSvc *patch.Service
	repoRoot string
}

// New creates a new Dispatcher. All parameters are required.
func New(repoSvc *repo.Service, tfSvc *terraform.Service, patchSvc *patch.Service, repoRoot string) *Dispatcher {
	return &Dispatcher{
		repoSvc:  repoSvc,
		tfSvc:    tfSvc,
		patchSvc: patchSvc,
		repoRoot: repoRoot,
	}
}

// ServeHTTP implements http.Handler. It parses the JSON-RPC 2.0 request,
// dispatches to the appropriate method handler, and writes a JSON-RPC 2.0
// response. All responses carry Content-Type: application/json.
func (d *Dispatcher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req JSONRPCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, nil, errCodeParseError, "Parse error")
		return
	}

	switch req.Method {
	case "initialize":
		d.handleInitialize(w, req)
	case "tools/list":
		d.handleToolsList(w, req)
	case "tools/call":
		d.handleToolsCall(w, req)
	default:
		writeError(w, req.ID, errCodeMethodNotFound, "Method not found")
	}
}

// handleInitialize responds to the MCP initialize handshake.
func (d *Dispatcher) handleInitialize(w http.ResponseWriter, req JSONRPCRequest) {
	result := map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities":    map[string]interface{}{"tools": map[string]interface{}{}},
		"serverInfo":      map[string]interface{}{"name": "terraformer", "version": "0.1.0"},
	}
	writeResult(w, req.ID, result)
}

// handleToolsList responds with all registered v0 tool definitions.
func (d *Dispatcher) handleToolsList(w http.ResponseWriter, req JSONRPCRequest) {
	writeResult(w, req.ID, map[string]interface{}{"tools": AllTools()})
}

// handleToolsCall dispatches a tools/call request to the appropriate handler.
func (d *Dispatcher) handleToolsCall(w http.ResponseWriter, req JSONRPCRequest) {
	var params struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}
	if err := json.Unmarshal(req.Params, &params); err != nil {
		writeError(w, req.ID, errCodeInvalidParams, "Invalid params: "+err.Error())
		return
	}

	text, err := CallTool(params.Name, params.Arguments, d.repoSvc, d.tfSvc, d.patchSvc, d.repoRoot)
	if err != nil {
		writeError(w, req.ID, errCodeInvalidParams, err.Error())
		return
	}

	result := map[string]interface{}{
		"content": []map[string]interface{}{
			{"type": "text", "text": text},
		},
	}
	writeResult(w, req.ID, result)
}

// writeResult encodes a successful JSON-RPC 2.0 response.
func writeResult(w http.ResponseWriter, id json.RawMessage, result interface{}) {
	resp := JSONRPCResponse{JSONRPC: "2.0", ID: id, Result: result}
	_ = json.NewEncoder(w).Encode(resp)
}

// writeError encodes a JSON-RPC 2.0 error response.
func writeError(w http.ResponseWriter, id json.RawMessage, code int, message string) {
	resp := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error:   &JSONRPCError{Code: code, Message: message},
	}
	_ = json.NewEncoder(w).Encode(resp)
}
