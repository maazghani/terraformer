// Package httpserver implements the HTTP/JSON server for the terraformer MCP.
// It registers all v0 tools as POST endpoints under /tools/<name>, decodes
// JSON request bodies, delegates to the appropriate service, and encodes
// structured JSON responses. All state is derived from Config at construction;
// the repo root cannot be changed by any HTTP request.
package httpserver

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/maazghani/terraformer/internal/patch"
	"github.com/maazghani/terraformer/internal/repo"
	"github.com/maazghani/terraformer/internal/terraform"
	"github.com/maazghani/terraformer/internal/tools"
)

// Config holds the startup configuration for the server.
type Config struct {
	// RepoRoot is the absolute path to the repository root. Immutable after startup.
	RepoRoot string
	// Port is the TCP port to listen on (default 9001). Used only by ListenAndServe.
	Port int
	// LogLevel is the minimum log level (debug|info|warn|error). Default is "info".
	LogLevel string
	// MaxResponseBytes is the maximum size of response bodies in bytes. Default is 1048576 (1 MiB).
	MaxResponseBytes int
}

// Server is an HTTP/JSON server that exposes v0 terraformer tools.
type Server struct {
	cfg      Config
	mux      *http.ServeMux
	repoSvc  *repo.Service
	tfSvc    *terraform.Service
	patchSvc *patch.Service
	logger   *jsonLogger
}

// New creates a Server with the given configuration and services. log receives
// structured JSON log lines (one per request). It must not be nil.
func New(cfg Config, repoSvc *repo.Service, tfSvc *terraform.Service, patchSvc *patch.Service, log io.Writer) *Server {
	logLevel := cfg.LogLevel
	if logLevel == "" {
		logLevel = "info"
	}
	s := &Server{
		cfg:      cfg,
		mux:      http.NewServeMux(),
		repoSvc:  repoSvc,
		tfSvc:    tfSvc,
		patchSvc: patchSvc,
		logger:   newJSONLogger(log, logLevel),
	}
	s.registerRoutes()
	return s
}

// ServeHTTP implements http.Handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// ListenAndServe starts the HTTP server on addr (e.g. ":9001").
func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, s)
}

// registerRoutes wires all v0 tool endpoints. Only POST is accepted; any other
// method receives a 405 response.
func (s *Server) registerRoutes() {
	s.mux.HandleFunc("POST /tools/terraform_init", s.handleTerraformInit)
	s.mux.HandleFunc("POST /tools/terraform_fmt", s.handleTerraformFmt)
	s.mux.HandleFunc("POST /tools/terraform_validate", s.handleTerraformValidate)
	s.mux.HandleFunc("POST /tools/terraform_plan", s.handleTerraformPlan)
	s.mux.HandleFunc("POST /tools/terraform_show_json", s.handleTerraformShowJSON)
	s.mux.HandleFunc("POST /tools/list_repo_files", s.handleListRepoFiles)
	s.mux.HandleFunc("POST /tools/read_repo_file", s.handleReadRepoFile)
	s.mux.HandleFunc("POST /tools/apply_patch", s.handleApplyPatch)
	s.mux.HandleFunc("POST /tools/check_desired_state", s.handleCheckDesiredState)
}

// writeJSON encodes v as JSON and writes it to w with a 200 status and the
// application/json content type.
func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		// Encoding errors after the header is sent cannot be corrected.
		_ = err
	}
}

// errorJSON writes a JSON error response with the given HTTP status code.
func errorJSON(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

// decodeBody decodes the JSON request body into dst. Returns false and writes a
// 400 response when decoding fails.
func decodeBody(w http.ResponseWriter, r *http.Request, dst interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return false
	}
	return true
}

// ---------------------------------------------------------------------------
// Tool handlers
// ---------------------------------------------------------------------------

func (s *Server) handleTerraformInit(w http.ResponseWriter, r *http.Request) {
	var req terraform.InitRequest
	if !decodeBody(w, r, &req) {
		s.logger.log("warn", "terraform_init: invalid request body")
		return
	}
	resp := s.tfSvc.Init(req)
	s.logger.log("info", "terraform_init")
	writeJSON(w, resp)
}

func (s *Server) handleTerraformFmt(w http.ResponseWriter, r *http.Request) {
	var req terraform.FmtRequest
	if !decodeBody(w, r, &req) {
		s.logger.log("warn", "terraform_fmt: invalid request body")
		return
	}
	resp := s.tfSvc.Fmt(req)
	s.logger.log("info", "terraform_fmt")
	writeJSON(w, resp)
}

func (s *Server) handleTerraformValidate(w http.ResponseWriter, r *http.Request) {
	var req terraform.ValidateRequest
	if !decodeBody(w, r, &req) {
		s.logger.log("warn", "terraform_validate: invalid request body")
		return
	}
	resp := s.tfSvc.Validate(req)
	s.logger.log("info", "terraform_validate")
	writeJSON(w, resp)
}

func (s *Server) handleTerraformPlan(w http.ResponseWriter, r *http.Request) {
	var req terraform.PlanRequest
	if !decodeBody(w, r, &req) {
		s.logger.log("warn", "terraform_plan: invalid request body")
		return
	}
	resp := s.tfSvc.Plan(req)
	s.logger.log("info", "terraform_plan")
	writeJSON(w, resp)
}

func (s *Server) handleTerraformShowJSON(w http.ResponseWriter, r *http.Request) {
	var req terraform.ShowJSONRequest
	if !decodeBody(w, r, &req) {
		s.logger.log("warn", "terraform_show_json: invalid request body")
		return
	}
	resp := s.tfSvc.ShowJSON(req)
	s.logger.log("info", "terraform_show_json")
	writeJSON(w, resp)
}

func (s *Server) handleListRepoFiles(w http.ResponseWriter, r *http.Request) {
	var req tools.ListRepoFilesRequest
	if !decodeBody(w, r, &req) {
		s.logger.log("warn", "list_repo_files: invalid request body")
		return
	}
	resp := tools.ListRepoFiles(s.repoSvc, req)
	s.logger.log("info", "list_repo_files")
	writeJSON(w, resp)
}

func (s *Server) handleReadRepoFile(w http.ResponseWriter, r *http.Request) {
	var req tools.ReadRepoFileRequest
	if !decodeBody(w, r, &req) {
		s.logger.log("warn", "read_repo_file: invalid request body")
		return
	}
	// Apply global maxResponseBytes limit if configured
	if s.cfg.MaxResponseBytes > 0 {
		maxBytes := int64(s.cfg.MaxResponseBytes)
		if req.MaxBytes == 0 || req.MaxBytes > maxBytes {
			req.MaxBytes = maxBytes
		}
	}
	resp := tools.ReadRepoFile(s.repoSvc, req)
	s.logger.log("info", "read_repo_file")
	writeJSON(w, resp)
}

// applyPatchRequest mirrors the JSON shape for the apply_patch tool.
type applyPatchRequest struct {
	Files []applyPatchFileOp `json:"files"`
}

// applyPatchFileOp is a single file operation within an apply_patch request.
type applyPatchFileOp struct {
	Path      string `json:"path"`
	Operation string `json:"operation"`
	Content   string `json:"content,omitempty"`
}

// applyPatchResponse is the structured response for apply_patch.
type applyPatchResponse struct {
	OK            bool     `json:"ok"`
	ChangedFiles  []string `json:"changed_files"`
	RejectedFiles []string `json:"rejected_files"`
	Warnings      []string `json:"warnings"`
}

func (s *Server) handleApplyPatch(w http.ResponseWriter, r *http.Request) {
	var req applyPatchRequest
	if !decodeBody(w, r, &req) {
		s.logger.log("warn", "apply_patch: invalid request body")
		return
	}

	ops := make([]patch.FileOperation, len(req.Files))
	for i, f := range req.Files {
		ops[i] = patch.FileOperation{
			Path:      f.Path,
			Operation: f.Operation,
			Content:   f.Content,
		}
	}

	result, err := s.patchSvc.ApplyPatch(patch.ApplyPatchRequest{Files: ops})
	var resp applyPatchResponse
	if err != nil {
		resp = applyPatchResponse{
			OK:            false,
			ChangedFiles:  []string{},
			RejectedFiles: []string{},
			Warnings:      []string{err.Error()},
		}
	} else {
		resp = applyPatchResponse{
			OK:            result.OK,
			ChangedFiles:  result.ChangedFiles,
			RejectedFiles: result.RejectedFiles,
			Warnings:      result.Warnings,
		}
	}

	s.logger.log("info", "apply_patch")
	writeJSON(w, resp)
}

func (s *Server) handleCheckDesiredState(w http.ResponseWriter, r *http.Request) {
	var req tools.CheckDesiredStateRequest
	if !decodeBody(w, r, &req) {
		s.logger.log("warn", "check_desired_state: invalid request body")
		return
	}
	resp := tools.CheckDesiredState(s.cfg.RepoRoot, req)
	s.logger.log("info", "check_desired_state")
	writeJSON(w, resp)
}
