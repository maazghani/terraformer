package terraform

import (
	"reflect"
	"testing"

	"github.com/maazghani/terraformer/internal/runner"
)

func TestService_Init_DefaultArgs(t *testing.T) {
	fake := runner.NewFakeRunner()
	fake.Register("terraform", runner.Result{ExitCode: 0}, nil)

	root := t.TempDir()
	svc, err := NewService(fake, root)
	if err != nil {
		t.Fatalf("NewService: %v", err)
	}

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
	if calls[0].WorkingDir != root {
		t.Errorf("expected WorkingDir=%q, got %q", root, calls[0].WorkingDir)
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
	if resp.Command.WorkingDir != "." {
		t.Errorf("unexpected response Command.WorkingDir: %q", resp.Command.WorkingDir)
	}
}

// TestService_Init_UpgradeAddsFlag verifies that InitRequest.Upgrade=true
// causes terraform init to be invoked with the -upgrade flag.
func TestService_Init_UpgradeAddsFlag(t *testing.T) {
	fake := runner.NewFakeRunner()
	fake.Register("terraform", runner.Result{ExitCode: 0}, nil)
	svc, err := NewService(fake, t.TempDir())
	if err != nil {
		t.Fatalf("NewService: %v", err)
	}

	svc.Init(InitRequest{Upgrade: true})

	if err := fake.AssertCalled("terraform", []string{"init", "-input=false", "-upgrade"}); err != nil {
		t.Fatalf("unexpected runner call: %v", err)
	}
}

// TestService_Init_UpgradeFalseDoesNotAddFlag verifies that
// InitRequest.Upgrade=false (the default) does not add -upgrade.
func TestService_Init_UpgradeFalseDoesNotAddFlag(t *testing.T) {
	fake := runner.NewFakeRunner()
	fake.Register("terraform", runner.Result{ExitCode: 0}, nil)
	svc, err := NewService(fake, t.TempDir())
	if err != nil {
		t.Fatalf("NewService: %v", err)
	}

	svc.Init(InitRequest{Upgrade: false})

	if err := fake.AssertCalled("terraform", []string{"init", "-input=false"}); err != nil {
		t.Fatalf("unexpected runner call: %v", err)
	}
}

// TestService_Init_BackendFalseAddsFlag verifies that InitRequest.Backend=false
// causes terraform init to be invoked with -backend=false to skip backend init.
func TestService_Init_BackendFalseAddsFlag(t *testing.T) {
	fake := runner.NewFakeRunner()
	fake.Register("terraform", runner.Result{ExitCode: 0}, nil)
	svc, err := NewService(fake, t.TempDir())
	if err != nil {
		t.Fatalf("NewService: %v", err)
	}

	backendFalse := false
	svc.Init(InitRequest{Backend: &backendFalse})

	if err := fake.AssertCalled("terraform", []string{"init", "-input=false", "-backend=false"}); err != nil {
		t.Fatalf("unexpected runner call: %v", err)
	}
}

// TestService_Init_BackendNilDoesNotAddFlag verifies that omitting Backend
// (nil pointer) does not add any -backend flag, relying on Terraform's default.
func TestService_Init_BackendNilDoesNotAddFlag(t *testing.T) {
	fake := runner.NewFakeRunner()
	fake.Register("terraform", runner.Result{ExitCode: 0}, nil)
	svc, err := NewService(fake, t.TempDir())
	if err != nil {
		t.Fatalf("NewService: %v", err)
	}

	svc.Init(InitRequest{Backend: nil})

	if err := fake.AssertCalled("terraform", []string{"init", "-input=false"}); err != nil {
		t.Fatalf("unexpected runner call: %v", err)
	}
}

// TestService_Init_UpgradeAndBackendFalse verifies both flags can be combined.
func TestService_Init_UpgradeAndBackendFalse(t *testing.T) {
	fake := runner.NewFakeRunner()
	fake.Register("terraform", runner.Result{ExitCode: 0}, nil)
	svc, err := NewService(fake, t.TempDir())
	if err != nil {
		t.Fatalf("NewService: %v", err)
	}

	backendFalse := false
	svc.Init(InitRequest{Upgrade: true, Backend: &backendFalse})

	if err := fake.AssertCalled("terraform", []string{
		"init", "-input=false", "-upgrade", "-backend=false",
	}); err != nil {
		t.Fatalf("unexpected runner call: %v", err)
	}
}
