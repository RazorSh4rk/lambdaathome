#!/bin/bash
# Gets a single function by name.

set -e
source "$(dirname "$0")/config.sh"

NAME="${1:-$FUNCTION_NAME}"

echo "Getting function '$NAME'..."
RESPONSE=$(curl -s -w "\n%{http_code}" \
    -H "Authorization: $AUTH" \
    "$BASE_URL/function/get/$NAME")

CODE=$(echo "$RESPONSE" | tail -1)
BODY=$(echo "$RESPONSE" | sed '$d')

check_response "$CODE" 200 "Get function"
echo "$BODY"
