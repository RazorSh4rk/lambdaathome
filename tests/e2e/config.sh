#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

BASE_URL="${BASE_URL:-http://localhost:8080}"
AUTH="$(cat "$PROJECT_ROOT/passfile")"

RUNTIME_NAME="e2e-runtime"
FUNCTION_NAME="e2e-testfn"
FUNCTION_PORT="9002"

red()   { printf "\033[31m%s\033[0m\n" "$1"; }
green() { printf "\033[32m%s\033[0m\n" "$1"; }

check_response() {
    local code="$1"
    local expected="$2"
    local label="$3"
    if [ "$code" -eq "$expected" ]; then
        green "PASS: $label (HTTP $code)"
    else
        red  "FAIL: $label (expected $expected, got $code)"
        return 1
    fi
}
