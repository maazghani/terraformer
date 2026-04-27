// Command terraformer-mcp is the CLI entry point for the terraformer MCP server.
// It parses configuration, validates the repo root, constructs dependencies,
// and starts the HTTP/JSON server on the configured port.
//
// Usage:
//
//	terraformer-mcp --repo-root=/path/to/repo [--port=9001]
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/maazghani/terraformer/internal/httpserver"
	"github.com/maazghani/terraformer/internal/patch"
	"github.com/maazghani/terraformer/internal/repo"
	"github.com/maazghani/terraformer/internal/runner"
	"github.com/maazghani/terraformer/internal/terraform"
)

// validateLogLevel returns true if level is a valid log level.
func validateLogLevel(level string) bool {
	switch level {
	case "debug", "info", "warn", "error":
		return true
	default:
		return false
	}
}

func main() {
	repoRootFlag := flag.String("repo-root", "", "absolute path to the repository root (required)")
	portFlag := flag.Int("port", 9001, "TCP port to listen on")
	logLevelFlag := flag.String("log-level", "info", "log level (debug|info|warn|error)")
	flag.Parse()

	if *repoRootFlag == "" {
		fmt.Fprintln(os.Stderr, "error: --repo-root is required")
		os.Exit(1)
	}

	if !validateLogLevel(*logLevelFlag) {
		fmt.Fprintf(os.Stderr, "error: invalid log level %q (must be debug|info|warn|error)\n", *logLevelFlag)
		os.Exit(1)
	}

	repoSvc, err := repo.New(*repoRootFlag)
	if err != nil {
		log.Fatalf("repo.New: %v", err)
	}

	localRunner := runner.NewLocalRunner()

	tfSvc, err := terraform.NewService(localRunner, *repoRootFlag)
	if err != nil {
		log.Fatalf("terraform.NewService: %v", err)
	}

	patchSvc, err := patch.New(*repoRootFlag)
	if err != nil {
		log.Fatalf("patch.New: %v", err)
	}

	cfg := httpserver.Config{
		RepoRoot: *repoRootFlag,
		Port:     *portFlag,
		LogLevel: *logLevelFlag,
	}

	srv := httpserver.New(cfg, repoSvc, tfSvc, patchSvc, os.Stdout)

	addr := fmt.Sprintf(":%d", *portFlag)
	log.Printf(`{"level":"info","message":"terraformer starting","addr":%q,"repo_root":%q}`, addr, *repoRootFlag)

	if err := srv.ListenAndServe(addr); err != nil {
		log.Fatalf(`{"level":"fatal","message":"server error","error":%q}`, err.Error())
	}
}
