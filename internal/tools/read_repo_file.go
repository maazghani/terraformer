package tools

import (
	"github.com/maazghani/terraformer/internal/repo"
)

// ReadRepoFileRequest is the request struct for the read_repo_file tool.
type ReadRepoFileRequest struct {
	// Path is the repo-relative file path to read.
	Path string `json:"path"`
	// MaxBytes caps the content returned. When zero, the full file is read.
	MaxBytes int64 `json:"max_bytes"`
}

// ReadRepoFileResponse is the response for the read_repo_file tool.
type ReadRepoFileResponse struct {
	// OK is true when the read succeeded.
	OK bool `json:"ok"`
	// Path is the repo-relative path that was read.
	Path string `json:"path"`
	// Content is the file content, possibly truncated.
	Content string `json:"content"`
	// SizeBytes is the actual on-disk file size (not the truncated length).
	SizeBytes int64 `json:"size_bytes"`
	// Truncated is true when the file was larger than MaxBytes.
	Truncated bool `json:"truncated"`
	// Warnings contains non-fatal advisories.
	Warnings []string `json:"warnings"`
}

// ReadRepoFile executes the read_repo_file tool against svc.
func ReadRepoFile(svc *repo.Service, req ReadRepoFileRequest) ReadRepoFileResponse {
	svcResp, err := svc.ReadFile(repo.ReadFileRequest{
		Path:     req.Path,
		MaxBytes: req.MaxBytes,
	})
	if err != nil {
		return ReadRepoFileResponse{
			OK:       false,
			Path:     req.Path,
			Warnings: []string{err.Error()},
		}
	}

	return ReadRepoFileResponse{
		OK:        true,
		Path:      req.Path,
		Content:   svcResp.Content,
		SizeBytes: svcResp.SizeBytes,
		Truncated: svcResp.Truncated,
		Warnings:  []string{},
	}
}
