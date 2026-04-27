// Package terraform implements the v0 Terraform command service. All
// invocations go through a runner.Runner. Only the allowlisted Terraform
// subcommands (init, fmt, validate, plan, show -json) are supported.
package terraform

import (
	"encoding/json"
	"time"

	"github.com/maazghani/terraformer/internal/runner"
	"github.com/maazghani/terraformer/internal/safety"
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

// ValidateRequest is the request for terraform_validate.
type ValidateRequest struct {
	JSON bool `json:"json"`
}

// ValidateResponse is the response for terraform_validate.
type ValidateResponse struct {
	OK          bool         `json:"ok"`
	Command     CommandInfo  `json:"command"`
	Stdout      string       `json:"stdout"`
	Stderr      string       `json:"stderr"`
	ExitCode    int          `json:"exit_code"`
	DurationMs  int64        `json:"duration_ms"`
	Diagnostics []Diagnostic `json:"diagnostics"`
	Warnings    []string     `json:"warnings"`
}

// Validate runs `terraform validate`, optionally with -json for structured
// diagnostics.
func (s *Service) Validate(req ValidateRequest) ValidateResponse {
	args := []string{"validate"}
	if req.JSON {
		args = append(args, "-json")
	}

	cmd := runner.Command{Name: "terraform", Args: args, WorkingDir: s.repoRoot}
	res, err := s.runner.Run(cmd)

	diags := []Diagnostic{}
	if req.JSON {
		diags = parseValidateJSON(res.Stdout)
	}

	return ValidateResponse{
		OK:          err == nil && res.ExitCode == 0,
		Command:     CommandInfo{Name: cmd.Name, Args: cmd.Args, WorkingDir: cmd.WorkingDir},
		Stdout:      res.Stdout,
		Stderr:      res.Stderr,
		ExitCode:    res.ExitCode,
		DurationMs:  durationMs(res.Duration),
		Diagnostics: diags,
		Warnings:    []string{},
	}
}

// validateJSONOutput mirrors the relevant subset of `terraform validate -json`.
type validateJSONOutput struct {
	Diagnostics []validateJSONDiagnostic `json:"diagnostics"`
}

type validateJSONDiagnostic struct {
	Severity string             `json:"severity"`
	Summary  string             `json:"summary"`
	Detail   string             `json:"detail"`
	Range    *validateJSONRange `json:"range,omitempty"`
}

type validateJSONRange struct {
	Filename string             `json:"filename"`
	Start    validateJSONOffset `json:"start"`
}

type validateJSONOffset struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

func parseValidateJSON(stdout string) []Diagnostic {
	var out validateJSONOutput
	if err := json.Unmarshal([]byte(stdout), &out); err != nil {
		return []Diagnostic{}
	}
	diags := make([]Diagnostic, 0, len(out.Diagnostics))
	for _, d := range out.Diagnostics {
		nd := Diagnostic{
			Severity: d.Severity,
			Summary:  d.Summary,
			Detail:   d.Detail,
		}
		if d.Range != nil {
			nd.File = d.Range.Filename
			nd.Line = d.Range.Start.Line
			nd.Column = d.Range.Start.Column
		}
		diags = append(diags, nd)
	}
	return diags
}

// PlanRequest is the request for terraform_plan.
type PlanRequest struct {
	Out              string `json:"out"`
	DetailedExitCode bool   `json:"detailed_exitcode"`
	Refresh          bool   `json:"refresh"`
}

// PlanResponse is the response for terraform_plan.
type PlanResponse struct {
	OK                 bool         `json:"ok"`
	PlanStatus         string       `json:"plan_status"`
	DesiredStateStatus string       `json:"desired_state_status"`
	Command            CommandInfo  `json:"command"`
	Stdout             string       `json:"stdout"`
	Stderr             string       `json:"stderr"`
	ExitCode           int          `json:"exit_code"`
	DurationMs         int64        `json:"duration_ms"`
	Diagnostics        []Diagnostic `json:"diagnostics"`
	Warnings           []string     `json:"warnings"`
}

// Plan runs `terraform plan` with safe defaults. It always passes -input=false
// and treats plan success as necessary-but-not-sufficient: DesiredStateStatus
// is always reported as "not_checked". When Out is provided, it is validated
// to be a repo-relative path inside the repo root before being passed to
// Terraform; unsafe paths cause Plan to return without invoking the runner.
func (s *Service) Plan(req PlanRequest) PlanResponse {
	args := []string{"plan", "-input=false"}
	if !req.Refresh {
		args = append(args, "-refresh=false")
	}
	if req.DetailedExitCode {
		args = append(args, "-detailed-exitcode")
	}
	if req.Out != "" {
		if _, err := safety.ResolvePath(s.repoRoot, req.Out); err != nil {
			return PlanResponse{
				OK:                 false,
				PlanStatus:         "failure",
				DesiredStateStatus: "not_checked",
				Command: CommandInfo{
					Name:       "terraform",
					Args:       args,
					WorkingDir: s.repoRoot,
				},
				Diagnostics: []Diagnostic{},
				Warnings:    []string{"unsafe plan out path rejected: " + err.Error()},
			}
		}
		args = append(args, "-out="+req.Out)
	}

	cmd := runner.Command{Name: "terraform", Args: args, WorkingDir: s.repoRoot}
	res, err := s.runner.Run(cmd)

	resp := PlanResponse{
		Command:            CommandInfo{Name: cmd.Name, Args: cmd.Args, WorkingDir: cmd.WorkingDir},
		Stdout:             res.Stdout,
		Stderr:             res.Stderr,
		ExitCode:           res.ExitCode,
		DurationMs:         durationMs(res.Duration),
		Diagnostics:        []Diagnostic{},
		Warnings:           []string{},
		DesiredStateStatus: "not_checked",
	}

	switch {
	case err != nil:
		resp.OK = false
		resp.PlanStatus = "failure"
	case req.DetailedExitCode && res.ExitCode == 0:
		resp.OK = true
		resp.PlanStatus = "no_changes"
	case req.DetailedExitCode && res.ExitCode == 2:
		resp.OK = true
		resp.PlanStatus = "changes_present"
	case res.ExitCode == 0:
		resp.OK = true
		resp.PlanStatus = "ok"
	default:
		resp.OK = false
		resp.PlanStatus = "failure"
	}
	return resp
}
