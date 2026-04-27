# Scripts Directory

This directory contains utility scripts for terraformer development and demonstration.

## demo.sh

End-to-end demonstration script that showcases all terraformer v0 capabilities.

### What It Does

1. **Builds** the `terraformer-mcp` binary
2. **Starts** the server against the `testdata/fixtures/valid-basic` fixture
3. **Exercises** all 9 MCP tools via curl:
   - `list_repo_files` — List Terraform files
   - `read_repo_file` — Read file contents
   - `terraform_fmt` — Format check
   - `terraform_validate` — Validate configuration
   - `terraform_plan` — Generate execution plan
   - `terraform_show_json` — Parse plan to JSON
   - `apply_patch` — Write a demo file
   - `check_desired_state` — Check desired state (minimal in v0)
   - `terraform_init` — Initialize Terraform
4. **Shuts down** the server cleanly

### Prerequisites

- **Go toolchain** (for building)
- **curl** (for making API calls)
- **jq** (for pretty-printing JSON responses)
- **Terraform** (optional; demo uses fake runner by default)

### Running the Demo

From the repository root:

```bash
./scripts/demo.sh
```

The script will:
- Build the binary if needed
- Find a free port (starting from 9001)
- Start the server
- Execute all tool demonstrations
- Display colorized output
- Clean up automatically on exit

### Output

The demo prints:
- Colored status messages (INFO/WARN/ERROR)
- Pretty-printed JSON responses for each tool
- Success/failure indicators

Server logs are saved to `server.log` in the repository root.

### Exit Codes

- `0` — All tools tested successfully
- `1` — Build failure or server startup failure

### Customization

Edit the script to:
- Change the fixture directory
- Modify tool request parameters
- Add additional tool invocations
- Change log levels or output formats

## Future Scripts

Additional scripts may be added for:
- Integration testing
- Performance benchmarking
- Client SDK generation
