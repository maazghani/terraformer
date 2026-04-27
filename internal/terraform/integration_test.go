package terraform_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/maazghani/terraformer/internal/runner"
	"github.com/maazghani/terraformer/internal/terraform"
)

// skipIfTerraformNotEnabled skips the test unless TERRAFORMER_RUN_INTEGRATION=1.
func skipIfTerraformNotEnabled(t *testing.T) {
	t.Helper()
	if os.Getenv("TERRAFORMER_RUN_INTEGRATION") != "1" {
		t.Skip("Skipping integration test: TERRAFORMER_RUN_INTEGRATION is not set to 1")
	}
}

// isTerraformAvailable checks if terraform is available in PATH.
func isTerraformAvailable() bool {
	_, err := exec.LookPath("terraform")
	return err == nil
}

// skipIfTerraformNotAvailable skips the test if terraform is not installed.
func skipIfTerraformNotAvailable(t *testing.T) {
	t.Helper()
	if !isTerraformAvailable() {
		t.Skip("Skipping integration test: terraform not found in PATH")
	}
}

// copyFixture copies a fixture directory from testdata/fixtures into a
// temporary directory and returns the absolute path to the copy.
// The copy is used instead of the original to avoid mutating committed fixtures.
func copyFixture(t *testing.T, fixtureName string) string {
	t.Helper()

	// Find the testdata/fixtures directory
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	// Navigate up to find testdata/fixtures
	// In tests, we're in internal/terraform, so go up two levels
	fixtureSource := filepath.Join(wd, "..", "..", "testdata", "fixtures", fixtureName)
	if _, err := os.Stat(fixtureSource); os.IsNotExist(err) {
		t.Fatalf("Fixture %q does not exist at %s", fixtureName, fixtureSource)
	}

	// Create a temporary directory for the copy
	tempDir := t.TempDir()
	fixtureDest := filepath.Join(tempDir, fixtureName)

	// Copy the fixture
	if err := copyDir(fixtureSource, fixtureDest); err != nil {
		t.Fatalf("Failed to copy fixture %q: %v", fixtureName, err)
	}

	return fixtureDest
}

// copyDir recursively copies a directory tree.
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Compute destination path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		// Copy file
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(dstPath, data, info.Mode())
	})
}

// TestIntegrationFmtValidBasic tests terraform fmt on a valid fixture.
func TestIntegrationFmtValidBasic(t *testing.T) {
	skipIfTerraformNotEnabled(t)
	skipIfTerraformNotAvailable(t)

	fixturePath := copyFixture(t, "valid-basic")
	r := runner.NewLocalRunner()
	svc, err := terraform.NewService(r, fixturePath)
	if err != nil {
		t.Fatalf("NewService failed: %v", err)
	}

	req := terraform.FmtRequest{
		Check:     false,
		Recursive: false,
	}

	resp := svc.Fmt(req)
	if !resp.OK {
		t.Errorf("Fmt failed: exit_code=%d, stderr=%s", resp.ExitCode, resp.Stderr)
	}
}

// TestIntegrationValidateValidBasic tests terraform validate on a valid fixture.
func TestIntegrationValidateValidBasic(t *testing.T) {
	skipIfTerraformNotEnabled(t)
	skipIfTerraformNotAvailable(t)

	fixturePath := copyFixture(t, "valid-basic")
	r := runner.NewLocalRunner()
	svc, err := terraform.NewService(r, fixturePath)
	if err != nil {
		t.Fatalf("NewService failed: %v", err)
	}

	// Run init first
	initReq := terraform.InitRequest{
		Upgrade: false,
		Backend: nil,
	}
	initResp := svc.Init(initReq)
	if !initResp.OK {
		t.Fatalf("Init failed: exit_code=%d, stderr=%s", initResp.ExitCode, initResp.Stderr)
	}

	// Run validate
	req := terraform.ValidateRequest{
		JSON: true,
	}
	resp := svc.Validate(req)
	if !resp.OK {
		t.Errorf("Validate failed: exit_code=%d, diagnostics=%+v", resp.ExitCode, resp.Diagnostics)
	}
}

// TestIntegrationValidateInvalidHCL tests terraform validate on invalid HCL.
func TestIntegrationValidateInvalidHCL(t *testing.T) {
	skipIfTerraformNotEnabled(t)
	skipIfTerraformNotAvailable(t)

	fixturePath := copyFixture(t, "invalid-hcl")
	r := runner.NewLocalRunner()
	svc, err := terraform.NewService(r, fixturePath)
	if err != nil {
		t.Fatalf("NewService failed: %v", err)
	}

	// Run validate (no init needed for syntax errors)
	req := terraform.ValidateRequest{
		JSON: true,
	}
	resp := svc.Validate(req)
	if resp.OK {
		t.Errorf("Validate should fail on invalid HCL, but got OK")
	}
	if resp.ExitCode == 0 {
		t.Errorf("Validate should have non-zero exit code on invalid HCL, got 0")
	}
}

// TestIntegrationPlanValidBasic tests terraform plan on a valid fixture.
func TestIntegrationPlanValidBasic(t *testing.T) {
	skipIfTerraformNotEnabled(t)
	skipIfTerraformNotAvailable(t)

	fixturePath := copyFixture(t, "plan-basic")
	r := runner.NewLocalRunner()
	svc, err := terraform.NewService(r, fixturePath)
	if err != nil {
		t.Fatalf("NewService failed: %v", err)
	}

	// Run init first
	initReq := terraform.InitRequest{
		Upgrade: false,
		Backend: nil,
	}
	initResp := svc.Init(initReq)
	if !initResp.OK {
		t.Fatalf("Init failed: exit_code=%d, stderr=%s", initResp.ExitCode, initResp.Stderr)
	}

	// Run plan
	planPath := filepath.Join(fixturePath, "plan.tfplan")
	req := terraform.PlanRequest{
		Out:              planPath,
		DetailedExitCode: true,
		Refresh:          true,
	}
	resp := svc.Plan(req)

	// Plan should succeed (either with or without changes)
	// Exit code 0 = no changes, exit code 2 = changes present
	if resp.ExitCode != 0 && resp.ExitCode != 2 {
		t.Errorf("Plan failed: exit_code=%d, stderr=%s", resp.ExitCode, resp.Stderr)
	}
}
