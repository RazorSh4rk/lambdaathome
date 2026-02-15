#!/bin/bash
# Creates a test function using the e2e-runtime.
# Requires: create_runtime.sh has been run first.

set -e
source "$(dirname "$0")/config.sh"

TMPDIR=$(mktemp -d)
trap "rm -rf $TMPDIR" EXIT

mkdir -p "$TMPDIR/code"
cat > "$TMPDIR/code/server.py" <<'PYTHON'
from http.server import HTTPServer, BaseHTTPRequestHandler

class Handler(BaseHTTPRequestHandler):
    def do_GET(self):
        self.send_response(200)
        self.end_headers()
        self.wfile.write(b'hello from e2e test v1')

HTTPServer(('', 8080), Handler).serve_forever()
PYTHON

(cd "$TMPDIR/code" && zip -q "$TMPDIR/code.zip" server.py)

echo "Uploading function '$FUNCTION_NAME'..."
RESPONSE=$(curl -s -w "\n%{http_code}" \
    -X POST \
    -H "Authorization: $AUTH" \
    -F "name=$FUNCTION_NAME" \
    -F "tag=$FUNCTION_NAME:latest" \
    -F "runtime=$RUNTIME_NAME" \
    -F "port=$FUNCTION_PORT" \
    -F "volume=" \
    -F "file=@$TMPDIR/code.zip" \
    "$BASE_URL/function/upload")

CODE=$(echo "$RESPONSE" | tail -1)
BODY=$(echo "$RESPONSE" | sed '$d')

check_response "$CODE" 200 "Create function"
echo "$BODY"

echo ""
echo "Build started in background. Use get_function.sh to check when it's ready."
