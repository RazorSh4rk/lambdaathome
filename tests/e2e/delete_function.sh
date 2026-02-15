#!/bin/bash
# Deletes a function by name.
# Kills container, removes container, removes image, cleans DB.

set -e
source "$(dirname "$0")/config.sh"

NAME="${1:-$FUNCTION_NAME}"

echo "Deleting function '$NAME'..."
RESPONSE=$(curl -s -w "\n%{http_code}" \
    -X DELETE \
    -H "Authorization: $AUTH" \
    "$BASE_URL/function/delete/$NAME")

CODE=$(echo "$RESPONSE" | tail -1)
BODY=$(echo "$RESPONSE" | sed '$d')

check_response "$CODE" 200 "Delete function"
echo "$BODY"
