package terraform

import (
	"reflect"
	"sort"
	"testing"

	"github.com/maazghani/terraformer/internal/runner"
)

// allowedExportedMethods is the v0 allowlist of *Service exported methods that
// invoke Terraform. If you add a new exported method that calls runner.Run,
// add it here AND make sure its first arg is part of the v0 subcommand
// allowlist enforced below.
var allowedExportedMethods = []string{"Fmt", "Init", "Plan", "ShowJSON", "Validate"}

var allowedSubcommands = map[string]bool{
	"init":     true,
	"fmt":      true,
	"validate": true,
	"plan":     true,
	"show":     true,
}

func TestService_ExportedMethodsAreAllowlisted(t *testing.T) {
	svc, err := NewService(runner.NewFakeRunner(), t.TempDir())
	if err != nil {
		t.Fatalf("NewService: %v", err)
	}
	typ := reflect.TypeOf(svc)

	var got []string
	for i := 0; i < typ.NumMethod(); i++ {
		got = append(got, typ.Method(i).Name)
	}
	sort.Strings(got)

	want := append([]string(nil), allowedExportedMethods...)
	sort.Strings(want)

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Service exported methods drifted from v0 allowlist.\n got:  %v\n want: %v\n"+
			"If you added a new method, update allowedExportedMethods AND ensure its\n"+
			"subcommand is in allowedSubcommands. apply/destroy must NEVER appear here.",
			got, want)
	}
}

func TestService_NeverInvokesForbiddenSubcommand(t *testing.T) {
	cases := []struct {
		name string
		fn   func(s *Service)
	}{
		{"Init", func(s *Service) { s.Init(InitRequest{}) }},
		{"Fmt", func(s *Service) { s.Fmt(FmtRequest{}) }},
		{"Validate", func(s *Service) { s.Validate(ValidateRequest{JSON: true}) }},
		{"Plan", func(s *Service) {
			s.Plan(PlanRequest{Out: ".terraformer/plan.tfplan", DetailedExitCode: true, Refresh: true})
		}},
		{"ShowJSON", func(s *Service) { s.ShowJSON(ShowJSONRequest{PlanPath: ".terraformer/plan.tfplan"}) }},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			fake := runner.NewFakeRunner()
			fake.Register("terraform", runner.Result{ExitCode: 0}, nil)
			svc, err := NewService(fake, t.TempDir())
			if err != nil {
				t.Fatalf("NewService: %v", err)
			}

			c.fn(svc)

			for _, call := range fake.Calls() {
				if call.Name != "terraform" {
					t.Errorf("non-terraform executable invoked: %q (no shell or arbitrary binary allowed)", call.Name)
				}
				if len(call.Args) == 0 {
					t.Errorf("empty args for call %+v", call)
					continue
				}
				sub := call.Args[0]
				if sub == "apply" || sub == "destroy" {
					t.Errorf("forbidden subcommand %q invoked via %s", sub, c.name)
				}
				if !allowedSubcommands[sub] {
					t.Errorf("subcommand %q is not in v0 allowlist (via %s)", sub, c.name)
				}
			}
		})
	}
}

func TestService_NeverInvokesShell(t *testing.T) {
	// Drive every exported method and assert no command name resembles a shell.
	fake := runner.NewFakeRunner()
	fake.Register("terraform", runner.Result{ExitCode: 0}, nil)
	svc, err := NewService(fake, t.TempDir())
	if err != nil {
		t.Fatalf("NewService: %v", err)
	}

	svc.Init(InitRequest{})
	svc.Fmt(FmtRequest{})
	svc.Validate(ValidateRequest{})
	svc.Plan(PlanRequest{Refresh: true})
	svc.ShowJSON(ShowJSONRequest{PlanPath: ".terraformer/plan.tfplan"})

	forbidden := map[string]bool{
		"sh": true, "bash": true, "zsh": true, "/bin/sh": true,
		"/bin/bash": true, "cmd": true, "cmd.exe": true, "powershell": true,
	}
	for _, call := range fake.Calls() {
		if forbidden[call.Name] {
			t.Errorf("forbidden shell-like executable invoked: %q", call.Name)
		}
		if call.Name != "terraform" {
			t.Errorf("only terraform may be invoked from terraform.Service; got %q", call.Name)
		}
	}
}
