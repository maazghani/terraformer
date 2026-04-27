// Package tools provides internal tool handlers for the MCP server.
// Each handler accepts a structured request, delegates to the appropriate
// service, and returns a structured response.
package tools

import (
	"github.com/maazghani/terraformer/internal/repo"
)

const defaultMaxFiles = 200

// ListRepoFilesRequest is the request struct for the list_repo_files tool.
type ListRepoFilesRequest struct {
	// Path is the repo-relative directory to list. Defaults to ".".
	Path string `json:"path"`
	// MaxFiles caps the number of file entries returned. When zero, a default is used.
	MaxFiles int `json:"max_files"`
}

// FileInfo is a single file entry returned by list_repo_files.
type FileInfo struct {
	// Path is the repo-relative, slash-separated path to the file.
	Path string `json:"path"`
	// SizeBytes is the file size in bytes.
	SizeBytes int64 `json:"size_bytes"`
	// Kind is "file" or "dir".
	Kind string `json:"kind"`
}

// ListRepoFilesResponse is the response for the list_repo_files tool.
type ListRepoFilesResponse struct {
	// OK is true when the operation succeeded.
	OK bool `json:"ok"`
	// Files is the list of file entries.
	Files []FileInfo `json:"files"`
	// Truncated is true when the result was capped by MaxFiles.
	Truncated bool `json:"truncated"`
	// Warnings contains non-fatal advisories.
	Warnings []string `json:"warnings"`
}

// ListRepoFiles executes the list_repo_files tool against svc.
func ListRepoFiles(svc *repo.Service, req ListRepoFilesRequest) ListRepoFilesResponse {
	maxFiles := req.MaxFiles
	if maxFiles <= 0 {
		maxFiles = defaultMaxFiles
	}

	svcResp, err := svc.ListFiles(repo.ListFilesRequest{
		Path:     req.Path,
		MaxFiles: maxFiles,
	})
	if err != nil {
		return ListRepoFilesResponse{
			OK:       false,
			Files:    []FileInfo{},
			Warnings: []string{err.Error()},
		}
	}

	files := make([]FileInfo, len(svcResp.Files))
	for i, f := range svcResp.Files {
		files[i] = FileInfo{
			Path:      f.Path,
			SizeBytes: f.SizeBytes,
			Kind:      f.Kind,
		}
	}

	return ListRepoFilesResponse{
		OK:        true,
		Files:     files,
		Truncated: svcResp.Truncated,
		Warnings:  []string{},
	}
}
