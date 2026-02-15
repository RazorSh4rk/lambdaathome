#!/bin/bash
# Lists all deployed functions.

set -e
source "$(dirname "$0")/config.sh"

echo "Listing all functions..."
RESPONSE=$(curl -s -w "\n%{http_code}" \
    -H "Authorization: $AUTH" \
    "$BASE_URL/function/list")

CODE=$(echo "$RESPONSE" | tail -1)
BODY=$(echo "$RESPONSE" | sed '$d')

check_response "$CODE" 200 "List functions"
echo "$BODY"
