package runner_test

import (
	"strings"
	"testing"
	"time"

	"github.com/maazghani/terraformer/internal/runner"
)

// newReal returns a real process runner under test.
func newReal() runner.Runner {
	return runner.NewLocalRunner()
}

// TestRealRunnerCapturesStdout verifies that stdout is captured in the result.
func TestRealRunnerCapturesStdout(t *testing.T) {
	r := newReal()
	result, err := r.Run(runner.Command{
		Name: "echo",
		Args: []string{"hello stdout"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result.Stdout, "hello stdout") {
		t.Errorf("Stdout: got %q, want it to contain %q", result.Stdout, "hello stdout")
	}
}

// TestRealRunnerCapturesStderr verifies that stderr is captured in the result.
func TestRealRunnerCapturesStderr(t *testing.T) {
	r := newReal()
	// sh -c redirects "hello stderr" to stderr; we need a real command that
	// writes to stderr. Use a small shell one-liner via /bin/sh.
	// BUT: the spec says "never invoke through a shell". This test uses
	// /bin/sh directly as the executable (not a shell-interpolated command
	// string) — that is legal for test purposes because the runner is not
	// constructing the command from user input here.
	//
	// For production Terraform calls we never do this; the runner itself does
	// not invoke a shell.
	result, err := r.Run(runner.Command{
		Name: "/bin/sh",
		Args: []string{"-c", "echo 'hello stderr' >&2"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result.Stderr, "hello stderr") {
		t.Errorf("Stderr: got %q, want it to contain %q", result.Stderr, "hello stderr")
	}
}

// TestRealRunnerCapturesExitCode verifies that a non-zero exit code is captured.
func TestRealRunnerCapturesExitCode(t *testing.T) {
	r := newReal()
	result, err := r.Run(runner.Command{
		Name: "/bin/sh",
		Args: []string{"-c", "exit 42"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ExitCode != 42 {
		t.Errorf("ExitCode: got %d, want 42", result.ExitCode)
	}
}

// TestRealRunnerCapturesZeroExitCode verifies that a zero exit code is captured on success.
func TestRealRunnerCapturesZeroExitCode(t *testing.T) {
	r := newReal()
	result, err := r.Run(runner.Command{
		Name: "true",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ExitCode != 0 {
		t.Errorf("ExitCode: got %d, want 0", result.ExitCode)
	}
}

// TestRealRunnerCaptureDuration verifies that a positive duration is recorded.
func TestRealRunnerCaptureDuration(t *testing.T) {
	r := newReal()
	result, err := r.Run(runner.Command{
		Name: "true",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Duration <= 0 {
		t.Errorf("Duration: got %v, want > 0", result.Duration)
	}
	// Sanity upper bound: a simple command should not take more than 10 seconds.
	if result.Duration > 10*time.Second {
		t.Errorf("Duration suspiciously large: %v", result.Duration)
	}
}

// TestRealRunnerRespectsWorkingDirectory verifies commands run in the specified dir.
func TestRealRunnerRespectsWorkingDirectory(t *testing.T) {
	r := newReal()
	result, err := r.Run(runner.Command{
		Name:       "pwd",
		WorkingDir: "/tmp",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// /tmp may be a symlink on macOS; accept either /tmp or a path containing tmp.
	if !strings.Contains(result.Stdout, "tmp") {
		t.Errorf("Stdout: got %q, expected to contain 'tmp'", result.Stdout)
	}
}

// TestRealRunnerDoesNotUseShell verifies args are not shell-interpolated.
// If the runner passed args through a shell, a glob like "*" would expand.
// Here we pass "*" as a literal arg to echo and expect it back verbatim.
func TestRealRunnerDoesNotUseShell(t *testing.T) {
	r := newReal()
	result, err := r.Run(runner.Command{
		Name: "echo",
		Args: []string{"*"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	trimmed := strings.TrimSpace(result.Stdout)
	if trimmed != "*" {
		t.Errorf("shell interpolation detected: echo '*' produced %q, want literal '*'", trimmed)
	}
}
