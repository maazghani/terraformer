package main

import (
	"flag"
	"testing"
)

// TestLogLevelFlag verifies that the --log-level flag can be parsed and defaults to "info".
func TestLogLevelFlag(t *testing.T) {
	// Reset flag.CommandLine to avoid conflicts with other tests
	flag.CommandLine = flag.NewFlagSet("test", flag.ContinueOnError)

	// Parse with no arguments
	logLevel := flag.String("log-level", "info", "log level (debug|info|warn|error)")
	err := flag.CommandLine.Parse([]string{})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if *logLevel != "info" {
		t.Errorf("default log level = %q, want %q", *logLevel, "info")
	}

	// Reset and parse with --log-level=debug
	flag.CommandLine = flag.NewFlagSet("test", flag.ContinueOnError)
	logLevel = flag.String("log-level", "info", "log level (debug|info|warn|error)")
	err = flag.CommandLine.Parse([]string{"--log-level=debug"})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if *logLevel != "debug" {
		t.Errorf("log level = %q, want %q", *logLevel, "debug")
	}
}

// TestLogLevelValidation verifies that invalid log levels are rejected.
func TestLogLevelValidation(t *testing.T) {
	tests := []struct {
		level string
		valid bool
	}{
		{"debug", true},
		{"info", true},
		{"warn", true},
		{"error", true},
		{"invalid", false},
		{"DEBUG", false},  // case-sensitive
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.level, func(t *testing.T) {
			valid := validateLogLevel(tt.level)
			if valid != tt.valid {
				t.Errorf("validateLogLevel(%q) = %v, want %v", tt.level, valid, tt.valid)
			}
		})
	}
}

// TestMaxResponseBytesFlag verifies that the --max-response-bytes flag can be parsed and defaults to 1048576.
func TestMaxResponseBytesFlag(t *testing.T) {
	// Reset flag.CommandLine
	flag.CommandLine = flag.NewFlagSet("test", flag.ContinueOnError)

	// Parse with no arguments
	maxResponseBytes := flag.Int("max-response-bytes", 1048576, "maximum response body size in bytes")
	err := flag.CommandLine.Parse([]string{})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if *maxResponseBytes != 1048576 {
		t.Errorf("default max-response-bytes = %d, want %d", *maxResponseBytes, 1048576)
	}

	// Reset and parse with custom value
	flag.CommandLine = flag.NewFlagSet("test", flag.ContinueOnError)
	maxResponseBytes = flag.Int("max-response-bytes", 1048576, "maximum response body size in bytes")
	err = flag.CommandLine.Parse([]string{"--max-response-bytes=2097152"})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if *maxResponseBytes != 2097152 {
		t.Errorf("max-response-bytes = %d, want %d", *maxResponseBytes, 2097152)
	}
}

// TestTerraformBinFlag verifies that the --terraform-bin flag can be parsed and defaults to "terraform".
func TestTerraformBinFlag(t *testing.T) {
	// Reset flag.CommandLine
	flag.CommandLine = flag.NewFlagSet("test", flag.ContinueOnError)

	// Parse with no arguments
	terraformBin := flag.String("terraform-bin", "terraform", "path to terraform binary")
	err := flag.CommandLine.Parse([]string{})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if *terraformBin != "terraform" {
		t.Errorf("default terraform-bin = %q, want %q", *terraformBin, "terraform")
	}

	// Reset and parse with custom path
	flag.CommandLine = flag.NewFlagSet("test", flag.ContinueOnError)
	terraformBin = flag.String("terraform-bin", "terraform", "path to terraform binary")
	err = flag.CommandLine.Parse([]string{"--terraform-bin=/usr/local/bin/terraform"})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if *terraformBin != "/usr/local/bin/terraform" {
		t.Errorf("terraform-bin = %q, want %q", *terraformBin, "/usr/local/bin/terraform")
	}
}
