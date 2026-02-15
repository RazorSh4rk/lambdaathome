#!/bin/bash
# Creates a test runtime (Python HTTP server).
# Run this before create_function.sh.

set -e
source "$(dirname "$0")/config.sh"

TMPDIR=$(mktemp -d)
trap "rm -rf $TMPDIR" EXIT

cat > "$TMPDIR/Dockerfile" <<'DOCKERFILE'
FROM python:3-alpine
COPY . /app
WORKDIR /app
CMD ["python", "server.py"]
DOCKERFILE

echo "Uploading runtime '$RUNTIME_NAME'..."
CODE=$(curl -s -o /dev/null -w "%{http_code}" \
    -X POST \
    -H "Authorization: $AUTH" \
    -F "file=@$TMPDIR/Dockerfile" \
    "$BASE_URL/runtime/upload/$RUNTIME_NAME")

check_response "$CODE" 200 "Create runtime"
