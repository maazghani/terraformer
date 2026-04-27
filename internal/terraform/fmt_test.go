package terraform

import (
	"errors"
	"testing"

	"github.com/maazghani/terraformer/internal/runner"
)

func TestService_Fmt_DefaultArgs(t *testing.T) {
	fake := runner.NewFakeRunner()
	fake.Register("terraform", runner.Result{ExitCode: 0}, nil)
	svc := NewService(fake, "/repo")

	resp := svc.Fmt(FmtRequest{})

	if !resp.OK {
		t.Fatalf("expected OK=true, got false; stderr=%q", resp.Stderr)
	}
	if err := fake.AssertCalled("terraform", []string{"fmt"}); err != nil {
		t.Fatalf("unexpected runner call: %v", err)
	}
	if calls := fake.Calls(); len(calls) != 1 || calls[0].WorkingDir != "/repo" {
		t.Fatalf("expected 1 call at /repo, got %+v", calls)
	}
}

func TestService_Fmt_Check(t *testing.T) {
	fake := runner.NewFakeRunner()
	fake.Register("terraform", runner.Result{ExitCode: 0}, nil)
	svc := NewService(fake, "/repo")

	svc.Fmt(FmtRequest{Check: true})

	if err := fake.AssertCalled("terraform", []string{"fmt", "-check"}); err != nil {
		t.Fatalf("unexpected runner call: %v", err)
	}
}

func TestService_Fmt_Recursive(t *testing.T) {
	fake := runner.NewFakeRunner()
	fake.Register("terraform", runner.Result{ExitCode: 0}, nil)
	svc := NewService(fake, "/repo")

	svc.Fmt(FmtRequest{Recursive: true})

	if err := fake.AssertCalled("terraform", []string{"fmt", "-recursive"}); err != nil {
		t.Fatalf("unexpected runner call: %v", err)
	}
}

func TestService_Fmt_CheckAndRecursive(t *testing.T) {
	fake := runner.NewFakeRunner()
	fake.Register("terraform", runner.Result{ExitCode: 0}, nil)
	svc := NewService(fake, "/repo")

	svc.Fmt(FmtRequest{Check: true, Recursive: true})

	if err := fake.AssertCalled("terraform", []string{"fmt", "-check", "-recursive"}); err != nil {
		t.Fatalf("unexpected runner call: %v", err)
	}
}

func TestService_Fmt_RunnerFailureMapsToNotOK(t *testing.T) {
	fake := runner.NewFakeRunner()
	fake.Register("terraform", runner.Result{ExitCode: 3, Stderr: "boom"}, errors.New("exec failed"))
	svc := NewService(fake, "/repo")

	resp := svc.Fmt(FmtRequest{})

	if resp.OK {
		t.Fatalf("expected OK=false on runner failure")
	}
	if resp.ExitCode != 3 {
		t.Errorf("expected ExitCode=3, got %d", resp.ExitCode)
	}
	if resp.Stderr != "boom" {
		t.Errorf("expected Stderr=boom, got %q", resp.Stderr)
	}
}
