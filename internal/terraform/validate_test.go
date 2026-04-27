package terraform

import (
	"testing"

	"github.com/maazghani/terraformer/internal/runner"
)

func TestService_Validate_DefaultArgsJSON(t *testing.T) {
	fake := runner.NewFakeRunner()
	fake.Register("terraform", runner.Result{ExitCode: 0, Stdout: `{"valid":true,"diagnostics":[]}`}, nil)
	svc := NewService(fake, "/repo")

	resp := svc.Validate(ValidateRequest{JSON: true})

	if !resp.OK {
		t.Fatalf("expected OK=true, got false; stderr=%q", resp.Stderr)
	}
	if err := fake.AssertCalled("terraform", []string{"validate", "-json"}); err != nil {
		t.Fatalf("unexpected runner call: %v", err)
	}
	calls := fake.Calls()
	if len(calls) != 1 || calls[0].WorkingDir != "/repo" {
		t.Fatalf("expected 1 call at /repo, got %+v", calls)
	}
}

func TestService_Validate_NonJSON(t *testing.T) {
	fake := runner.NewFakeRunner()
	fake.Register("terraform", runner.Result{ExitCode: 0}, nil)
	svc := NewService(fake, "/repo")

	svc.Validate(ValidateRequest{JSON: false})

	if err := fake.AssertCalled("terraform", []string{"validate"}); err != nil {
		t.Fatalf("unexpected runner call: %v", err)
	}
}
