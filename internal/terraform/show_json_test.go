package terraform

import (
	"testing"

	"github.com/maazghani/terraformer/internal/runner"
)

const showJSONOutput = `{
  "format_version": "1.2",
  "resource_changes": [
    {"address":"a","change":{"actions":["create"]}},
    {"address":"b","change":{"actions":["update"]}},
    {"address":"c","change":{"actions":["delete","create"]}},
    {"address":"d","change":{"actions":["delete"]}},
    {"address":"e","change":{"actions":["no-op"]}}
  ]
}`

func TestService_ShowJSON_DefaultArgs(t *testing.T) {
	fake := runner.NewFakeRunner()
	fake.Register("terraform", runner.Result{ExitCode: 0, Stdout: showJSONOutput}, nil)
	root := t.TempDir()
	svc, err := NewService(fake, root)
	if err != nil {
		t.Fatalf("NewService: %v", err)
	}

	resp := svc.ShowJSON(ShowJSONRequest{PlanPath: ".terraformer/plan.tfplan"})

	if !resp.OK {
		t.Fatalf("expected OK=true, got false; stderr=%q", resp.Stderr)
	}
	if err := fake.AssertCalled("terraform", []string{
		"show", "-json", "--", ".terraformer/plan.tfplan",
	}); err != nil {
		t.Fatalf("unexpected runner call: %v", err)
	}
	calls := fake.Calls()
	if len(calls) != 1 || calls[0].WorkingDir != root {
		t.Fatalf("expected 1 call at %q, got %+v", root, calls)
	}
	s := resp.PlanSummary
	if s.Create != 1 || s.Update != 1 || s.Delete != 1 || s.Replace != 1 || s.NoOp != 1 {
		t.Errorf("unexpected plan summary: %+v", s)
	}
}

func TestService_ShowJSON_MalformedJSONDoesNotPanic(t *testing.T) {
	fake := runner.NewFakeRunner()
	fake.Register("terraform", runner.Result{ExitCode: 0, Stdout: "not json"}, nil)
	svc, err := NewService(fake, t.TempDir())
	if err != nil {
		t.Fatalf("NewService: %v", err)
	}

	resp := svc.ShowJSON(ShowJSONRequest{PlanPath: ".terraformer/plan.tfplan"})

	if !resp.OK {
		t.Fatalf("expected OK=true even when summary cannot be parsed (command itself succeeded)")
	}
	s := resp.PlanSummary
	if s.Create != 0 || s.Update != 0 || s.Delete != 0 || s.Replace != 0 || s.NoOp != 0 {
		t.Errorf("expected zero summary on malformed json, got %+v", s)
	}
}

func TestService_ShowJSON_RequiresPlanPath(t *testing.T) {
	fake := runner.NewFakeRunner()
	fake.Register("terraform", runner.Result{ExitCode: 0}, nil)
	svc, err := NewService(fake, t.TempDir())
	if err != nil {
		t.Fatalf("NewService: %v", err)
	}

	resp := svc.ShowJSON(ShowJSONRequest{PlanPath: ""})

	if resp.OK {
		t.Fatalf("expected OK=false when plan_path is empty")
	}
	if len(fake.Calls()) != 0 {
		t.Errorf("runner must not be invoked when plan_path is empty; got %d calls", len(fake.Calls()))
	}
}
