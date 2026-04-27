package terraform

import (
	"reflect"
	"testing"

	"github.com/maazghani/terraformer/internal/runner"
)

func TestService_Init_DefaultArgs(t *testing.T) {
	fake := runner.NewFakeRunner()
	fake.Register("terraform", runner.Result{ExitCode: 0}, nil)

	svc := NewService(fake, "/repo")

	resp := svc.Init(InitRequest{})

	if !resp.OK {
		t.Fatalf("expected OK=true, got false; stderr=%q", resp.Stderr)
	}

	if err := fake.AssertCalled("terraform", []string{"init", "-input=false"}); err != nil {
		t.Fatalf("unexpected runner call: %v", err)
	}

	calls := fake.Calls()
	if len(calls) != 1 {
		t.Fatalf("expected exactly 1 call, got %d", len(calls))
	}
	if calls[0].WorkingDir != "/repo" {
		t.Errorf("expected WorkingDir=%q, got %q", "/repo", calls[0].WorkingDir)
	}
	if calls[0].Name != "terraform" {
		t.Errorf("expected command Name=terraform, got %q", calls[0].Name)
	}
	if !reflect.DeepEqual(resp.Command.Args, []string{"init", "-input=false"}) {
		t.Errorf("unexpected response Command.Args: %v", resp.Command.Args)
	}
	if resp.Command.Name != "terraform" {
		t.Errorf("unexpected response Command.Name: %q", resp.Command.Name)
	}
	if resp.Command.WorkingDir != "/repo" {
		t.Errorf("unexpected response Command.WorkingDir: %q", resp.Command.WorkingDir)
	}
}
