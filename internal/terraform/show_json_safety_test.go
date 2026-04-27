package terraform

import (
	"testing"

	"github.com/maazghani/terraformer/internal/runner"
)

func TestService_ShowJSON_RejectsAbsolutePlanPath(t *testing.T) {
	fake := runner.NewFakeRunner()
	fake.Register("terraform", runner.Result{ExitCode: 0}, nil)
	svc, err := NewService(fake, t.TempDir())
	if err != nil {
		t.Fatalf("NewService: %v", err)
	}

	resp := svc.ShowJSON(ShowJSONRequest{PlanPath: "/etc/passwd"})

	if resp.OK {
		t.Fatalf("expected OK=false for absolute plan_path")
	}
	if len(fake.Calls()) != 0 {
		t.Errorf("runner must not be invoked when plan_path is unsafe; got %d calls", len(fake.Calls()))
	}
}

func TestService_ShowJSON_RejectsTraversalPlanPath(t *testing.T) {
	fake := runner.NewFakeRunner()
	fake.Register("terraform", runner.Result{ExitCode: 0}, nil)
	svc, err := NewService(fake, t.TempDir())
	if err != nil {
		t.Fatalf("NewService: %v", err)
	}

	resp := svc.ShowJSON(ShowJSONRequest{PlanPath: "../../etc/passwd"})

	if resp.OK {
		t.Fatalf("expected OK=false for traversal plan_path")
	}
	if len(fake.Calls()) != 0 {
		t.Errorf("runner must not be invoked when plan_path escapes repo; got %d calls", len(fake.Calls()))
	}
}

// TestService_ShowJSON_RejectsDashPrefixedPlanPath verifies that a plan path
// beginning with '-' is rejected before being passed to the runner. Without
// this check a path like "-out=evil" could be interpreted as a Terraform flag
// even when "-- path" is used, depending on terraform version.
func TestService_ShowJSON_RejectsDashPrefixedPlanPath(t *testing.T) {
	fake := runner.NewFakeRunner()
	fake.Register("terraform", runner.Result{ExitCode: 0}, nil)
	svc, err := NewService(fake, t.TempDir())
	if err != nil {
		t.Fatalf("NewService: %v", err)
	}

	resp := svc.ShowJSON(ShowJSONRequest{PlanPath: "-out=evil.tfplan"})

	if resp.OK {
		t.Fatalf("expected OK=false for dash-prefixed plan_path")
	}
	if len(fake.Calls()) != 0 {
		t.Errorf("runner must not be invoked when plan_path starts with '-'; got %d calls", len(fake.Calls()))
	}
	if len(resp.Warnings) == 0 {
		t.Errorf("expected a warning explaining rejection, got none")
	}
}
