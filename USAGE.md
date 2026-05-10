# Usage Guide: Integrating terraformer with AI Coding Assistants

This guide shows how to integrate the terraformer MCP server with various AI coding assistants. Each assistant connects to the HTTP/JSON server running on `http://localhost:9001`.

## Prerequisites

1. Build terraformer:
   ```bash
   make build
   ```

2. Start the server pointing at your Terraform repository:
   ```bash
   ./terraformer-mcp --repo-root=/absolute/path/to/your/terraform/repo
   ```

The server must be running before the AI assistant can use the tools.

---

## Integration by Client

<details>
<summary><strong>GitHub Copilot (VS Code)</strong></summary>

### Configuration

GitHub Copilot in VS Code (agent mode) can connect to MCP servers via configuration. Add the following to your workspace or user `settings.json`:

```json
{
  "github.copilot.mcp.servers": {
    "terraformer": {
      "type": "http",
      "url": "http://localhost:9001"
    }
  }
}
```

Alternatively, if using a dedicated `mcp.json` file:

```json
{
  "servers": {
    "terraformer": {
      "type": "http",
      "url": "http://localhost:9001"
    }
  }
}
```

### Smoke Test

Ask Copilot:
> "List all Terraform files in this repository using the terraformer MCP server."

Copilot should invoke the `list_repo_files` tool and return results.

### Tool Reference

For detailed tool contracts, see [PLAN/spec/02-mcp-tool-contracts.md](PLAN/spec/02-mcp-tool-contracts.md).

</details>

<details>
<summary><strong>Cursor</strong></summary>

### Configuration

Cursor uses an MCP configuration file at `~/.cursor/mcp.json`:

```json
{
  "mcpServers": {
    "terraformer": {
      "type": "http",
      "url": "http://localhost:9001",
      "transport": "http"
    }
  }
}
```

After updating the config, restart Cursor.

### Smoke Test

In a Cursor chat, ask:
> "Use the terraformer server to validate the Terraform configuration in this repo."

Cursor should invoke `terraform_validate` and show diagnostics.

### Tool Reference

For detailed tool contracts, see [PLAN/spec/02-mcp-tool-contracts.md](PLAN/spec/02-mcp-tool-contracts.md).

</details>

<details>
<summary><strong>Claude Code / Claude Desktop</strong></summary>

### Configuration

Claude Desktop and Claude Code use a config file at:
- **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Linux**: `~/.config/Claude/claude_desktop_config.json`
- **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

Add the following configuration:

```json
{
  "mcpServers": {
    "terraformer": {
      "command": "echo",
      "args": ["HTTP server at http://localhost:9001"],
      "env": {},
      "http": {
        "url": "http://localhost:9001"
      }
    }
  }
}
```

Restart Claude Desktop after saving.

### Smoke Test

In Claude Desktop, ask:
> "Read the main.tf file from my Terraform repo using the terraformer MCP server."

Claude should invoke `read_repo_file` for `main.tf`.

### Tool Reference

For detailed tool contracts, see [PLAN/spec/02-mcp-tool-contracts.md](PLAN/spec/02-mcp-tool-contracts.md).

</details>

<details>
<summary><strong>Codex CLI</strong></summary>

### Protocol Support

The terraformer server now implements the **MCP Streamable HTTP transport** protocol,
which Codex uses for the startup handshake. The server dispatches JSON-RPC 2.0 requests
at `POST /`, handling `initialize`, `tools/list`, and `tools/call` as required by the
MCP specification (protocol version `2024-11-05`).

If you previously saw:
```
⚠ MCP client for `terraform` failed to start: MCP startup failed: handshaking with
MCP server failed: Unexpected content type: Some("text/plain; charset=utf-8; body: 404
page not found\n"), when send initialize request
```
this is now fixed. The server correctly handles the JSON-RPC `initialize` request at
`POST /` and returns `Content-Type: application/json`.

### Configuration

Codex CLI uses a TOML config file at `~/.codex/config.toml`:

```toml
[[mcp_servers]]
name = "terraformer"
type = "http"
url = "http://localhost:9001"
```

Restart the Codex CLI after updating the config.

### Smoke Test

Run:
```bash
codex ask "List all Terraform files using the terraformer server."
```

Codex should invoke `list_repo_files` and display results.

### Verification

You can manually verify the MCP handshake with curl:
```bash
# Step 1: initialize
curl -s -X POST http://localhost:9001/ \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}' | jq .

# Step 2: tools/list (all 9 tools)
curl -s -X POST http://localhost:9001/ \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":2,"method":"tools/list"}' | jq .result.tools[].name

# Step 3: tools/call
curl -s -X POST http://localhost:9001/ \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"list_repo_files","arguments":{}}}' | jq .
```

The existing `/tools/*` REST endpoints remain available for non-MCP clients.

### Tool Reference

For detailed tool contracts, see [PLAN/spec/02-mcp-tool-contracts.md](PLAN/spec/02-mcp-tool-contracts.md).

</details>

<details>
<summary><strong>Cline (VS Code Extension)</strong></summary>

### Configuration

Cline (formerly Claude Dev) is a VS Code extension. Configure MCP servers in VS Code settings:

1. Open VS Code Settings (JSON mode)
2. Add the following under Cline's MCP configuration:

```json
{
  "cline.mcpServers": {
    "terraformer": {
      "url": "http://localhost:9001",
      "type": "http"
    }
  }
}
```

Reload the VS Code window after saving settings.

### Smoke Test

In the Cline chat panel, ask:
> "Use terraformer to run terraform validate on this repo."

Cline should invoke `terraform_validate` via the MCP server.

### Tool Reference

For detailed tool contracts, see [PLAN/spec/02-mcp-tool-contracts.md](PLAN/spec/02-mcp-tool-contracts.md).

</details>

---

## Common Troubleshooting

### Server Not Responding

1. Verify the server is running (REST endpoint):
   ```bash
   curl http://localhost:9001/tools/list_repo_files \
     -X POST \
     -H "Content-Type: application/json" \
     -d '{"path": ".", "max_files": 10}'
   ```

2. Verify the MCP JSON-RPC endpoint (`POST /`) is working:
   ```bash
   curl http://localhost:9001/ \
     -X POST \
     -H "Content-Type: application/json" \
     -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{}}}'
   ```
   This should return JSON with `result.protocolVersion` set to `"2024-11-05"`.

3. Check the server logs for errors (it writes structured JSON logs to stdout/stderr).

4. Ensure the `--repo-root` path is absolute and points to an existing directory.

### Tool Not Found

If the AI assistant reports that a tool doesn't exist:

1. Verify the tool name matches the documented names (see [PLAN/spec/02-mcp-tool-contracts.md](PLAN/spec/02-mcp-tool-contracts.md))
2. Check that the server is running the latest build (`make build`)

### Path Safety Errors

If you see errors like "unsafe path" or "path traversal rejected":

- Ensure all file paths are **relative to the repo root**
- Do not use absolute paths or `../` in tool requests
- Check that symlinks don't escape the repository boundary

---

## Advanced Configuration

### Custom Port

If port 9001 is already in use, start the server on a different port:

```bash
./terraformer-mcp --repo-root=/path/to/repo --port=9002
```

Update your MCP client config to use `http://localhost:9002`.

### Custom Terraform Binary

If you need to use a specific Terraform binary:

```bash
./terraformer-mcp \
  --repo-root=/path/to/repo \
  --terraform-bin=/usr/local/bin/terraform
```

### Log Levels

Control log verbosity:

```bash
./terraformer-mcp --repo-root=/path/to/repo --log-level=debug
```

Valid levels: `debug`, `info`, `warn`, `error`

---

## Tool Reference

All tools accept and return structured JSON. See [PLAN/spec/02-mcp-tool-contracts.md](PLAN/spec/02-mcp-tool-contracts.md) for complete request/response schemas and examples.

### Available Tools

1. **terraform_init** — Initialize Terraform (download providers, modules)
2. **terraform_fmt** — Format `.tf` files
3. **terraform_validate** — Validate configuration syntax and references
4. **terraform_plan** — Generate execution plan (does not apply)
5. **terraform_show_json** — Parse binary plan file to JSON
6. **list_repo_files** — List files matching patterns
7. **read_repo_file** — Read file contents
8. **apply_patch** — Write file changes (does not run Terraform)
9. **check_desired_state** — Compare plan against desired state (minimal in v0)

### Forbidden Operations

The following are **not available** and will never be added to v0:

- `terraform apply` — Would modify infrastructure
- `terraform destroy` — Would destroy infrastructure
- `terraform import` — Could modify state
- `terraform workspace` — State management
- Arbitrary shell commands

---

## Example Workflows

### Workflow 1: Validate and Plan

```
1. list_repo_files → Find all .tf files
2. read_repo_file → Read main.tf and variables.tf
3. terraform_validate → Check syntax
4. terraform_plan → Generate plan
5. terraform_show_json → Parse plan details
```

### Workflow 2: Propose Changes

```
1. read_repo_file → Read existing configuration
2. apply_patch → Write proposed changes
3. terraform_fmt → Format the changes
4. terraform_validate → Ensure valid syntax
5. terraform_plan → Preview what would change
6. check_desired_state → Verify against expectations
```

The agent can iterate based on validation errors, but **cannot apply changes** to infrastructure.

---

## Support

For issues or questions:
- Check the [README.md](README.md) for general information
- Review tool contracts in [PLAN/spec/02-mcp-tool-contracts.md](PLAN/spec/02-mcp-tool-contracts.md)
- Inspect server logs (structured JSON to stdout/stderr)
