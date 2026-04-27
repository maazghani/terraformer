// Package diagnostics provides normalized diagnostic structures for Terraform
// validation errors, plan failures, and other structured feedback. Diagnostics
// are parsed from JSON outputs where available (terraform validate -json,
// terraform show -json) and normalized into a stable format for tool responses.
package diagnostics

// Diagnostic is a normalized diagnostic entry used across all Terraform tool
// responses. It represents an error, warning, or informational message with
// optional file location information.
type Diagnostic struct {
	Severity string `json:"severity"` // "error", "warning", "info"
	Summary  string `json:"summary"`  // Short description
	Detail   string `json:"detail"`   // Detailed explanation
	File     string `json:"file"`     // Source file name (repo-relative)
	Line     int    `json:"line"`     // Line number (1-indexed, 0 if unknown)
	Column   int    `json:"column"`   // Column number (1-indexed, 0 if unknown)
}
