package terraform

import (
	"strings"
	"testing"

	"github.com/maazghani/terraformer/internal/runner"
)

func TestService_Plan_RejectsAbsoluteOutPath(t *testing.T) {
	fake := runner.NewFakeRunner()
	fake.Register("terraform", runner.Result{ExitCode: 0}, nil)
	svc, err := NewService(fake, t.TempDir())
	if err != nil {
		t.Fatalf("NewService: %v", err)
	}

	resp := svc.Plan(PlanRequest{Out: "/tmp/evil.tfplan", Refresh: true})

	if resp.OK {
		t.Fatalf("expected OK=false for absolute out path")
	}
	if len(fake.Calls()) != 0 {
		t.Errorf("runner must not be invoked when out path is unsafe; got %d calls", len(fake.Calls()))
	}
	if !hasWarningContaining(resp.Warnings, "out") && resp.Stderr == "" {
		t.Errorf("expected an error/warning describing rejection; warnings=%v stderr=%q", resp.Warnings, resp.Stderr)
	}
}

func TestService_Plan_RejectsTraversalOutPath(t *testing.T) {
	fake := runner.NewFakeRunner()
	fake.Register("terraform", runner.Result{ExitCode: 0}, nil)
	svc, err := NewService(fake, t.TempDir())
	if err != nil {
		t.Fatalf("NewService: %v", err)
	}

	resp := svc.Plan(PlanRequest{Out: "../etc/passwd", Refresh: true})

	if resp.OK {
		t.Fatalf("expected OK=false for traversal out path")
	}
	if len(fake.Calls()) != 0 {
		t.Errorf("runner must not be invoked when out path escapes repo; got %d calls", len(fake.Calls()))
	}
}

func hasWarningContaining(warnings []string, sub string) bool {
	for _, w := range warnings {
		if strings.Contains(strings.ToLower(w), strings.ToLower(sub)) {
			return true
		}
	}
	return false
}
