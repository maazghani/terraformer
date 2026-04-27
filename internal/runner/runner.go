// Package runner defines the command runner abstraction used throughout terraformer.
// All command execution must go through a Runner implementation. No package may
// execute shell commands directly.
package runner

import "time"

// Command describes a command to be executed. WorkingDir must be set to the
// configured repo root or a subdirectory within it. Env contains additional
// environment variables (KEY=VALUE format). The command must not be invoked
// through a shell.
type Command struct {
	// Name is the executable name (e.g., "terraform"). No shell metacharacters.
	Name string
	// Args are the command-line arguments. Never passed through a shell.
	Args []string
	// WorkingDir is the directory in which the command runs.
	WorkingDir string
	// Env contains supplemental environment variables in KEY=VALUE format.
	// When empty the process inherits the parent environment.
	Env []string
}

// Result holds the captured output of a completed command.
type Result struct {
	// Stdout is the full standard output.
	Stdout string
	// Stderr is the full standard error.
	Stderr string
	// ExitCode is the process exit code. 0 indicates success.
	ExitCode int
	// Duration is the wall-clock time from process start to finish.
	Duration time.Duration
}

// Runner executes a Command and returns its Result.
// Implementations must not invoke a shell, must honor WorkingDir,
// and must capture stdout, stderr, exit code, and duration.
type Runner interface {
	Run(cmd Command) (Result, error)
}
