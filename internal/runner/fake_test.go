package runner_test

import (
	"errors"
	"testing"

	"github.com/maazghani/terraformer/internal/runner"
)

// TestFakeRunnerReturnsConfiguredResult verifies the fake runner returns the
// result that was registered for a given command name.
func TestFakeRunnerReturnsConfiguredResult(t *testing.T) {
	fake := runner.NewFakeRunner()
	fake.Register("terraform", runner.Result{Stdout: "fake output", ExitCode: 0}, nil)

	result, err := fake.Run(runner.Command{Name: "terraform", Args: []string{"validate"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Stdout != "fake output" {
		t.Errorf("Stdout: got %q, want %q", result.Stdout, "fake output")
	}
	if result.ExitCode != 0 {
		t.Errorf("ExitCode: got %d, want 0", result.ExitCode)
	}
}

// TestFakeRunnerReturnsConfiguredError verifies the fake can simulate failures.
func TestFakeRunnerReturnsConfiguredError(t *testing.T) {
	fake := runner.NewFakeRunner()
	fake.Register("terraform", runner.Result{}, errors.New("exec failed"))

	_, err := fake.Run(runner.Command{Name: "terraform"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "exec failed" {
		t.Errorf("error message: got %q, want %q", err.Error(), "exec failed")
	}
}

// TestFakeRunnerRecordsCalledCommands verifies every Run call is recorded for
// later assertion.
func TestFakeRunnerRecordsCalledCommands(t *testing.T) {
	fake := runner.NewFakeRunner()
	fake.Register("terraform", runner.Result{ExitCode: 0}, nil)

	cmd1 := runner.Command{Name: "terraform", Args: []string{"init"}, WorkingDir: "/repo"}
	cmd2 := runner.Command{Name: "terraform", Args: []string{"validate"}, WorkingDir: "/repo"}
	fake.Run(cmd1) //nolint:errcheck
	fake.Run(cmd2) //nolint:errcheck

	calls := fake.Calls()
	if len(calls) != 2 {
		t.Fatalf("calls: got %d, want 2", len(calls))
	}
	if calls[0].Args[0] != "init" {
		t.Errorf("calls[0].Args[0]: got %q, want %q", calls[0].Args[0], "init")
	}
	if calls[1].Args[0] != "validate" {
		t.Errorf("calls[1].Args[0]: got %q, want %q", calls[1].Args[0], "validate")
	}
}

// TestFakeRunnerAssertExactArgs verifies AssertCalled passes when the exact
// command was called with the expected arguments.
func TestFakeRunnerAssertExactArgs(t *testing.T) {
	fake := runner.NewFakeRunner()
	fake.Register("terraform", runner.Result{ExitCode: 0}, nil)
	fake.Run(runner.Command{Name: "terraform", Args: []string{"fmt", "-recursive"}}) //nolint:errcheck

	if err := fake.AssertCalled("terraform", []string{"fmt", "-recursive"}); err != nil {
		t.Errorf("AssertCalled returned unexpected error: %v", err)
	}
}

// TestFakeRunnerAssertExactArgsFails verifies AssertCalled fails when the
// command was not called with those arguments.
func TestFakeRunnerAssertExactArgsFails(t *testing.T) {
	fake := runner.NewFakeRunner()
	fake.Register("terraform", runner.Result{ExitCode: 0}, nil)
	fake.Run(runner.Command{Name: "terraform", Args: []string{"init"}}) //nolint:errcheck

	if err := fake.AssertCalled("terraform", []string{"validate"}); err == nil {
		t.Error("AssertCalled should have returned an error for mismatched args")
	}
}

// TestFakeRunnerAssertWorkingDir verifies that the working directory is
// captured and can be inspected via Calls().
func TestFakeRunnerAssertWorkingDir(t *testing.T) {
	fake := runner.NewFakeRunner()
	fake.Register("terraform", runner.Result{ExitCode: 0}, nil)
	fake.Run(runner.Command{Name: "terraform", Args: []string{"plan"}, WorkingDir: "/the/repo"}) //nolint:errcheck

	calls := fake.Calls()
	if len(calls) != 1 {
		t.Fatalf("calls: got %d, want 1", len(calls))
	}
	if calls[0].WorkingDir != "/the/repo" {
		t.Errorf("WorkingDir: got %q, want %q", calls[0].WorkingDir, "/the/repo")
	}
}

// TestFakeRunnerUnregisteredCommandErrors verifies that calling a command that
// has not been registered returns an error rather than silently succeeding.
func TestFakeRunnerUnregisteredCommandErrors(t *testing.T) {
	fake := runner.NewFakeRunner()
	_, err := fake.Run(runner.Command{Name: "notregistered"})
	if err == nil {
		t.Error("expected error for unregistered command, got nil")
	}
}
