package terraform

import (
	"testing"

	"github.com/maazghani/terraformer/internal/runner"
)

func TestService_Validate_DefaultArgsJSON(t *testing.T) {
	fake := runner.NewFakeRunner()
	fake.Register("terraform", runner.Result{ExitCode: 0, Stdout: `{"valid":true,"diagnostics":[]}`}, nil)
	root := t.TempDir()
	svc, err := NewService(fake, root)
	if err != nil {
		t.Fatalf("NewService: %v", err)
	}

	resp := svc.Validate(ValidateRequest{JSON: true})

	if !resp.OK {
		t.Fatalf("expected OK=true, got false; stderr=%q", resp.Stderr)
	}
	if err := fake.AssertCalled("terraform", []string{"validate", "-json"}); err != nil {
		t.Fatalf("unexpected runner call: %v", err)
	}
	calls := fake.Calls()
	if len(calls) != 1 || calls[0].WorkingDir != root {
		t.Fatalf("expected 1 call at %q, got %+v", root, calls)
	}
}

func TestService_Validate_NonJSON(t *testing.T) {
	fake := runner.NewFakeRunner()
	fake.Register("terraform", runner.Result{ExitCode: 0}, nil)
	svc, err := NewService(fake, t.TempDir())
	if err != nil {
		t.Fatalf("NewService: %v", err)
	}

	svc.Validate(ValidateRequest{JSON: false})

	if err := fake.AssertCalled("terraform", []string{"validate"}); err != nil {
		t.Fatalf("unexpected runner call: %v", err)
	}
}
