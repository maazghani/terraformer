#!/usr/bin/env bash
# demo.sh - End-to-end demonstration of terraformer MCP server
#
# This script:
# 1. Builds the terraformer-mcp binary
# 2. Starts it against the valid-basic test fixture
# 3. Exercises all 9 v0 tools via curl
# 4. Shuts down the server cleanly
#
# Prerequisites:
# - Go toolchain (for building)
# - curl (for API calls)
# - Terraform (optional; only needed if running with real terraform)

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}[INFO]${NC} $*"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $*"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $*"
}

# Find repo root
REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$REPO_ROOT"

log_info "Repository root: $REPO_ROOT"

# Build the binary
log_info "Building terraformer-mcp..."
make build

if [[ ! -x ./terraformer-mcp ]]; then
    log_error "Build failed: ./terraformer-mcp not found or not executable"
    exit 1
fi

# Fixture path
FIXTURE_PATH="$REPO_ROOT/testdata/fixtures/valid-basic"
if [[ ! -d "$FIXTURE_PATH" ]]; then
    log_error "Fixture directory not found: $FIXTURE_PATH"
    exit 1
fi

log_info "Using fixture: $FIXTURE_PATH"

# Find a free port
PORT=9001
while lsof -Pi :$PORT -sTCP:LISTEN -t >/dev/null 2>&1 ; do
    log_warn "Port $PORT in use, trying next port..."
    PORT=$((PORT + 1))
done

log_info "Starting server on port $PORT..."

# Start the server in the background
./terraformer-mcp --repo-root="$FIXTURE_PATH" --port=$PORT --log-level=info > server.log 2>&1 &
SERVER_PID=$!

log_info "Server started with PID $SERVER_PID"

# Function to cleanup on exit
cleanup() {
    if [[ -n "${SERVER_PID:-}" ]] && kill -0 "$SERVER_PID" 2>/dev/null; then
        log_info "Shutting down server (PID $SERVER_PID)..."
        kill "$SERVER_PID" 2>/dev/null || true
        wait "$SERVER_PID" 2>/dev/null || true
    fi
    if [[ -f server.log ]]; then
        log_info "Server logs available in server.log"
    fi
}

trap cleanup EXIT INT TERM

# Wait for server to be ready
log_info "Waiting for server to be ready..."
for i in {1..30}; do
    if curl -s -f "http://localhost:$PORT/tools/list_repo_files" \
        -X POST \
        -H "Content-Type: application/json" \
        -d '{"path": ".", "max_files": 1}' > /dev/null 2>&1; then
        log_info "Server is ready!"
        break
    fi
    if [[ $i -eq 30 ]]; then
        log_error "Server failed to start within 30 seconds"
        cat server.log
        exit 1
    fi
    sleep 1
done

BASE_URL="http://localhost:$PORT"

echo ""
log_info "==================== DEMO START ===================="
echo ""

# Tool 1: list_repo_files
log_info "1. Testing list_repo_files..."
curl -s "$BASE_URL/tools/list_repo_files" \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{
        "path": ".",
        "include_globs": ["*.tf"],
        "exclude_globs": [".terraform/**"],
        "max_files": 100
    }' | jq -C '.'

echo ""

# Tool 2: read_repo_file
log_info "2. Testing read_repo_file..."
curl -s "$BASE_URL/tools/read_repo_file" \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{
        "path": "main.tf",
        "max_bytes": 65536
    }' | jq -C '.'

echo ""

# Tool 3: terraform_fmt (check mode)
log_info "3. Testing terraform_fmt..."
curl -s "$BASE_URL/tools/terraform_fmt" \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{
        "check": true,
        "recursive": false
    }' | jq -C '.ok, .command, .exit_code'

echo ""

# Tool 4: terraform_validate
log_info "4. Testing terraform_validate..."
curl -s "$BASE_URL/tools/terraform_validate" \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{
        "json": true
    }' | jq -C '.ok, .command, .exit_code, .diagnostics'

echo ""

# Tool 5: terraform_plan
log_info "5. Testing terraform_plan..."
curl -s "$BASE_URL/tools/terraform_plan" \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{
        "detailed_exitcode": true,
        "out": ".terraformer/demo.tfplan",
        "refresh": false
    }' | jq -C '.ok, .plan_status, .desired_state_status, .command'

echo ""

# Tool 6: terraform_show_json
log_info "6. Testing terraform_show_json..."
curl -s "$BASE_URL/tools/terraform_show_json" \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{
        "plan_path": ".terraformer/demo.tfplan"
    }' | jq -C '.ok, .command, .plan_summary'

echo ""

# Tool 7: apply_patch (write a new file)
log_info "7. Testing apply_patch..."
curl -s "$BASE_URL/tools/apply_patch" \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{
        "files": [
            {
                "path": "demo_output.tf",
                "operation": "write",
                "content": "# Demo file created by demo.sh\n\noutput \"demo\" {\n  value = \"Hello from terraformer!\"\n}\n"
            }
        ]
    }' | jq -C '.'

echo ""

# Tool 8: check_desired_state
log_info "8. Testing check_desired_state..."
curl -s "$BASE_URL/tools/check_desired_state" \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{
        "desired_state": {
            "resources": []
        },
        "plan_json_path": ".terraformer/plan.json"
    }' | jq -C '.'

echo ""

# Tool 9: terraform_init
log_info "9. Testing terraform_init..."
curl -s "$BASE_URL/tools/terraform_init" \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{
        "upgrade": false,
        "backend": false
    }' | jq -C '.ok, .command, .exit_code'

echo ""
log_info "==================== DEMO COMPLETE ===================="
echo ""

log_info "All 9 tools tested successfully!"
log_info "Server will be shut down now."
