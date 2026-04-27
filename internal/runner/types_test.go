package runner_test

import (
	"testing"
	"time"

	"github.com/maazghani/terraformer/internal/runner"
)

// TestCommandShape verifies the Command struct fields exist and are typed correctly.
func TestCommandShape(t *testing.T) {
	cmd := runner.Command{
		Name:       "terraform",
		Args:       []string{"validate"},
		WorkingDir: "/tmp/repo",
		Env:        []string{"TF_LOG=DEBUG"},
	}

	if cmd.Name != "terraform" {
		t.Errorf("Name: got %q, want %q", cmd.Name, "terraform")
	}
	if len(cmd.Args) != 1 || cmd.Args[0] != "validate" {
		t.Errorf("Args: got %v, want [validate]", cmd.Args)
	}
	if cmd.WorkingDir != "/tmp/repo" {
		t.Errorf("WorkingDir: got %q, want %q", cmd.WorkingDir, "/tmp/repo")
	}
	if len(cmd.Env) != 1 || cmd.Env[0] != "TF_LOG=DEBUG" {
		t.Errorf("Env: got %v, want [TF_LOG=DEBUG]", cmd.Env)
	}
}

// TestResultShape verifies the Result struct fields exist and are typed correctly.
func TestResultShape(t *testing.T) {
	result := runner.Result{
		Stdout:   "output",
		Stderr:   "error output",
		ExitCode: 0,
		Duration: 42 * time.Millisecond,
	}

	if result.Stdout != "output" {
		t.Errorf("Stdout: got %q, want %q", result.Stdout, "output")
	}
	if result.Stderr != "error output" {
		t.Errorf("Stderr: got %q, want %q", result.Stderr, "error output")
	}
	if result.ExitCode != 0 {
		t.Errorf("ExitCode: got %d, want 0", result.ExitCode)
	}
	if result.Duration != 42*time.Millisecond {
		t.Errorf("Duration: got %v, want 42ms", result.Duration)
	}
}

// TestResultNonZeroExitCode verifies Result captures non-zero exit codes.
func TestResultNonZeroExitCode(t *testing.T) {
	result := runner.Result{
		ExitCode: 1,
		Stderr:   "something went wrong",
	}
	if result.ExitCode != 1 {
		t.Errorf("ExitCode: got %d, want 1", result.ExitCode)
	}
}

// TestRunnerInterface verifies the Runner interface is satisfied by a minimal implementation.
func TestRunnerInterface(t *testing.T) {
	// Compile-time check: a type that implements runner.Runner can be used as one.
	var _ runner.Runner = (*testStubRunner)(nil)
}

// testStubRunner is a minimal implementation used only to verify the interface shape.
type testStubRunner struct{}

func (s *testStubRunner) Run(cmd runner.Command) (runner.Result, error) {
	return runner.Result{}, nil
}
