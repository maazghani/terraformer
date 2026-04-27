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
