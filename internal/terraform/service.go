// Package terraform implements the v0 Terraform command service. All
// invocations go through a runner.Runner. Only the allowlisted Terraform
// subcommands (init, fmt, validate, plan, show -json) are supported.
package terraform

import (
	"time"

	"github.com/maazghani/terraformer/internal/runner"
)

// Service constructs and executes allowlisted Terraform commands inside a
// fixed repo root.
type Service struct {
	runner   runner.Runner
	repoRoot string
}

// NewService returns a Service that runs Terraform commands inside repoRoot
// using r.
func NewService(r runner.Runner, repoRoot string) *Service {
	return &Service{runner: r, repoRoot: repoRoot}
}

// CommandInfo mirrors the "command" object in tool responses.
type CommandInfo struct {
	Name       string   `json:"name"`
	Args       []string `json:"args"`
	WorkingDir string   `json:"working_dir"`
}

// Diagnostic is a normalized diagnostic entry.
type Diagnostic struct {
	Severity string `json:"severity"`
	Summary  string `json:"summary"`
	Detail   string `json:"detail"`
	File     string `json:"file"`
	Line     int    `json:"line"`
	Column   int    `json:"column"`
}

// InitRequest is the request for terraform_init.
type InitRequest struct{}

// InitResponse is the response for terraform_init.
type InitResponse struct {
	OK          bool         `json:"ok"`
	Command     CommandInfo  `json:"command"`
	Stdout      string       `json:"stdout"`
	Stderr      string       `json:"stderr"`
	ExitCode    int          `json:"exit_code"`
	DurationMs  int64        `json:"duration_ms"`
	Diagnostics []Diagnostic `json:"diagnostics"`
	Warnings    []string     `json:"warnings"`
}

// Init runs `terraform init -input=false` inside the repo root.
func (s *Service) Init(_ InitRequest) InitResponse {
	args := []string{"init", "-input=false"}

	cmd := runner.Command{
		Name:       "terraform",
		Args:       args,
		WorkingDir: s.repoRoot,
	}

	res, err := s.runner.Run(cmd)

	resp := InitResponse{
		Command: CommandInfo{
			Name:       cmd.Name,
			Args:       cmd.Args,
			WorkingDir: cmd.WorkingDir,
		},
		Stdout:      res.Stdout,
		Stderr:      res.Stderr,
		ExitCode:    res.ExitCode,
		DurationMs:  durationMs(res.Duration),
		Diagnostics: []Diagnostic{},
		Warnings:    []string{},
	}
	resp.OK = err == nil && res.ExitCode == 0
	return resp
}

func durationMs(d time.Duration) int64 {
	return d.Milliseconds()
}

// FmtRequest is the request for terraform_fmt.
type FmtRequest struct {
	Check     bool `json:"check"`
	Recursive bool `json:"recursive"`
}

// FmtResponse is the response for terraform_fmt.
type FmtResponse struct {
	OK          bool         `json:"ok"`
	Command     CommandInfo  `json:"command"`
	Stdout      string       `json:"stdout"`
	Stderr      string       `json:"stderr"`
	ExitCode    int          `json:"exit_code"`
	DurationMs  int64        `json:"duration_ms"`
	Diagnostics []Diagnostic `json:"diagnostics"`
	Warnings    []string     `json:"warnings"`
}

// Fmt runs `terraform fmt` with optional -check and -recursive flags.
func (s *Service) Fmt(req FmtRequest) FmtResponse {
	args := []string{"fmt"}
	if req.Check {
		args = append(args, "-check")
	}
	if req.Recursive {
		args = append(args, "-recursive")
	}

	cmd := runner.Command{Name: "terraform", Args: args, WorkingDir: s.repoRoot}
	res, err := s.runner.Run(cmd)

	return FmtResponse{
		OK:          err == nil && res.ExitCode == 0,
		Command:     CommandInfo{Name: cmd.Name, Args: cmd.Args, WorkingDir: cmd.WorkingDir},
		Stdout:      res.Stdout,
		Stderr:      res.Stderr,
		ExitCode:    res.ExitCode,
		DurationMs:  durationMs(res.Duration),
		Diagnostics: []Diagnostic{},
		Warnings:    []string{},
	}
}
