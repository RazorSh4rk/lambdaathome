#!/bin/bash
# Full cleanup: deletes the test function and runtime.
# Safe to run even if some resources don't exist.

source "$(dirname "$0")/config.sh"

echo "=== Cleanup ==="

echo "Deleting function '$FUNCTION_NAME'..."
curl -s -o /dev/null -w "  HTTP %{http_code}\n" \
    -X DELETE \
    -H "Authorization: $AUTH" \
    "$BASE_URL/function/delete/$FUNCTION_NAME"

echo "Deleting runtime '$RUNTIME_NAME'..."
curl -s -o /dev/null -w "  HTTP %{http_code}\n" \
    -X DELETE \
    -H "Authorization: $AUTH" \
    "$BASE_URL/runtime/delete/$RUNTIME_NAME"

green "Cleanup done."
