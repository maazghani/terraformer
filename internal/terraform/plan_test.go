package terraform

import (
	"testing"

	"github.com/maazghani/terraformer/internal/runner"
)

func TestService_Plan_DefaultArgsAndDetailedExitcodeNoChanges(t *testing.T) {
	fake := runner.NewFakeRunner()
	fake.Register("terraform", runner.Result{ExitCode: 0}, nil)
	svc := NewService(fake, "/repo")

	resp := svc.Plan(PlanRequest{
		Out:              ".terraformer/plan.tfplan",
		DetailedExitCode: true,
		Refresh:          true,
	})

	if !resp.OK {
		t.Fatalf("expected OK=true, got false")
	}
	if resp.PlanStatus != "no_changes" {
		t.Errorf("plan_status: got %q want no_changes", resp.PlanStatus)
	}
	if resp.DesiredStateStatus != "not_checked" {
		t.Errorf("desired_state_status: got %q want not_checked", resp.DesiredStateStatus)
	}
	if err := fake.AssertCalled("terraform", []string{
		"plan", "-input=false", "-detailed-exitcode", "-out=.terraformer/plan.tfplan",
	}); err != nil {
		t.Fatalf("unexpected runner call: %v", err)
	}
	calls := fake.Calls()
	if len(calls) != 1 || calls[0].WorkingDir != "/repo" {
		t.Fatalf("expected 1 call at /repo, got %+v", calls)
	}
}

func TestService_Plan_ExitCode2MapsToChangesPresent(t *testing.T) {
	fake := runner.NewFakeRunner()
	fake.Register("terraform", runner.Result{ExitCode: 2}, nil)
	svc := NewService(fake, "/repo")

	resp := svc.Plan(PlanRequest{
		Out:              ".terraformer/plan.tfplan",
		DetailedExitCode: true,
		Refresh:          true,
	})

	if !resp.OK {
		t.Fatalf("expected OK=true on exit 2 with detailed_exitcode (changes present is success)")
	}
	if resp.PlanStatus != "changes_present" {
		t.Errorf("plan_status: got %q want changes_present", resp.PlanStatus)
	}
	if resp.DesiredStateStatus != "not_checked" {
		t.Errorf("plan success must NOT imply desired_state_status; got %q", resp.DesiredStateStatus)
	}
}

func TestService_Plan_ExitCode1MapsToFailure(t *testing.T) {
	fake := runner.NewFakeRunner()
	fake.Register("terraform", runner.Result{ExitCode: 1, Stderr: "boom"}, nil)
	svc := NewService(fake, "/repo")

	resp := svc.Plan(PlanRequest{
		Out:              ".terraformer/plan.tfplan",
		DetailedExitCode: true,
		Refresh:          true,
	})

	if resp.OK {
		t.Fatalf("expected OK=false on exit 1")
	}
	if resp.PlanStatus != "failure" {
		t.Errorf("plan_status: got %q want failure", resp.PlanStatus)
	}
}

func TestService_Plan_RefreshFalseAddsFlag(t *testing.T) {
	fake := runner.NewFakeRunner()
	fake.Register("terraform", runner.Result{ExitCode: 0}, nil)
	svc := NewService(fake, "/repo")

	svc.Plan(PlanRequest{
		Out:              ".terraformer/plan.tfplan",
		DetailedExitCode: true,
		Refresh:          false,
	})

	if err := fake.AssertCalled("terraform", []string{
		"plan", "-input=false", "-refresh=false", "-detailed-exitcode", "-out=.terraformer/plan.tfplan",
	}); err != nil {
		t.Fatalf("unexpected runner call: %v", err)
	}
}

func TestService_Plan_NoDetailedExitCodeNoOut(t *testing.T) {
	fake := runner.NewFakeRunner()
	fake.Register("terraform", runner.Result{ExitCode: 0}, nil)
	svc := NewService(fake, "/repo")

	resp := svc.Plan(PlanRequest{Refresh: true})

	if !resp.OK {
		t.Fatalf("expected OK=true")
	}
	if err := fake.AssertCalled("terraform", []string{"plan", "-input=false"}); err != nil {
		t.Fatalf("unexpected runner call: %v", err)
	}
	if resp.DesiredStateStatus != "not_checked" {
		t.Errorf("desired_state_status: got %q want not_checked", resp.DesiredStateStatus)
	}
}
