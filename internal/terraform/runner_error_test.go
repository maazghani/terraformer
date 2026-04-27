package terraform

import (
	"errors"
	"strings"
	"testing"

	"github.com/maazghani/terraformer/internal/runner"
)

// runnerErr is a sentinel error used across runner-failure tests.
var runnerErr = errors.New("exec: terraform binary not found")

// hasWarningContainingStr returns true when any of the warnings contains sub
// (case-insensitive). The helper below shadows the one in plan_safety_test.go
// but both are in the same package so only one copy is needed; the
// plan_safety_test.go version is renamed to avoid a duplicate.

func TestService_Init_RunnerErrorSurfacedInWarnings(t *testing.T) {
	fake := runner.NewFakeRunner()
	fake.Register("terraform", runner.Result{ExitCode: 0}, runnerErr)
	svc, err := NewService(fake, t.TempDir())
	if err != nil {
		t.Fatalf("NewService: %v", err)
	}

	resp := svc.Init(InitRequest{})

	if resp.OK {
		t.Fatalf("expected OK=false when runner returns an error")
	}
	if !warningContains(resp.Warnings, runnerErr.Error()) {
		t.Errorf("expected runner error in Warnings, got %v", resp.Warnings)
	}
}

func TestService_Fmt_RunnerErrorSurfacedInWarnings(t *testing.T) {
	fake := runner.NewFakeRunner()
	fake.Register("terraform", runner.Result{ExitCode: 0}, runnerErr)
	svc, err := NewService(fake, t.TempDir())
	if err != nil {
		t.Fatalf("NewService: %v", err)
	}

	resp := svc.Fmt(FmtRequest{})

	if resp.OK {
		t.Fatalf("expected OK=false when runner returns an error")
	}
	if !warningContains(resp.Warnings, runnerErr.Error()) {
		t.Errorf("expected runner error in Warnings, got %v", resp.Warnings)
	}
}

func TestService_Validate_RunnerErrorSurfacedInWarnings(t *testing.T) {
	fake := runner.NewFakeRunner()
	fake.Register("terraform", runner.Result{ExitCode: 0}, runnerErr)
	svc, err := NewService(fake, t.TempDir())
	if err != nil {
		t.Fatalf("NewService: %v", err)
	}

	resp := svc.Validate(ValidateRequest{})

	if resp.OK {
		t.Fatalf("expected OK=false when runner returns an error")
	}
	if !warningContains(resp.Warnings, runnerErr.Error()) {
		t.Errorf("expected runner error in Warnings, got %v", resp.Warnings)
	}
}

func TestService_Plan_RunnerErrorSurfacedInWarnings(t *testing.T) {
	fake := runner.NewFakeRunner()
	fake.Register("terraform", runner.Result{ExitCode: 0}, runnerErr)
	svc, err := NewService(fake, t.TempDir())
	if err != nil {
		t.Fatalf("NewService: %v", err)
	}

	resp := svc.Plan(PlanRequest{Refresh: true})

	if resp.OK {
		t.Fatalf("expected OK=false when runner returns an error")
	}
	if !warningContains(resp.Warnings, runnerErr.Error()) {
		t.Errorf("expected runner error in Warnings, got %v", resp.Warnings)
	}
}

func TestService_ShowJSON_RunnerErrorSurfacedInWarnings(t *testing.T) {
	fake := runner.NewFakeRunner()
	fake.Register("terraform", runner.Result{ExitCode: 0}, runnerErr)
	svc, err := NewService(fake, t.TempDir())
	if err != nil {
		t.Fatalf("NewService: %v", err)
	}

	resp := svc.ShowJSON(ShowJSONRequest{PlanPath: ".terraformer/plan.tfplan"})

	if resp.OK {
		t.Fatalf("expected OK=false when runner returns an error")
	}
	if !warningContains(resp.Warnings, runnerErr.Error()) {
		t.Errorf("expected runner error in Warnings, got %v", resp.Warnings)
	}
}

// warningContains reports whether any element of warnings contains sub
// (case-insensitive).
func warningContains(warnings []string, sub string) bool {
	for _, w := range warnings {
		if strings.Contains(strings.ToLower(w), strings.ToLower(sub)) {
			return true
		}
	}
	return false
}
