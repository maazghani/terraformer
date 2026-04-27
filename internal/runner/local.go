package runner

import (
	"bytes"
	"os/exec"
	"time"
)

// LocalRunner is a real process runner. It executes commands directly without
// a shell, captures stdout, stderr, exit code, and duration.
type LocalRunner struct{}

// NewLocalRunner returns a LocalRunner ready for use.
func NewLocalRunner() *LocalRunner {
	return &LocalRunner{}
}

// Run executes cmd.Name with cmd.Args in cmd.WorkingDir.
// It never invokes a shell. Stdout and stderr are fully captured.
// If the command exits with a non-zero code, Run returns a nil error but
// sets Result.ExitCode accordingly. A non-nil error is only returned for
// failures that prevent the process from starting at all.
func (r *LocalRunner) Run(cmd Command) (Result, error) {
	c := exec.Command(cmd.Name, cmd.Args...) //nolint:gosec // args are typed, not shell-interpolated
	if cmd.WorkingDir != "" {
		c.Dir = cmd.WorkingDir
	}
	if len(cmd.Env) > 0 {
		c.Env = cmd.Env
	}

	var stdout, stderr bytes.Buffer
	c.Stdout = &stdout
	c.Stderr = &stderr

	start := time.Now()
	runErr := c.Run()
	duration := time.Since(start)

	exitCode := 0
	if runErr != nil {
		if exitErr, ok := runErr.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
			runErr = nil // non-zero exit is not an infrastructure error
		} else {
			return Result{}, runErr
		}
	}

	return Result{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: exitCode,
		Duration: duration,
	}, nil
}
