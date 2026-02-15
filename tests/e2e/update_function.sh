#!/bin/bash
# Re-uploads the same function name with different code and port.
# This triggers the upsert path: old container is torn down, new one deployed.
# Requires: create_function.sh has been run first.

set -e
source "$(dirname "$0")/config.sh"

UPDATED_PORT="9003"

TMPDIR=$(mktemp -d)
trap "rm -rf $TMPDIR" EXIT

mkdir -p "$TMPDIR/code"
cat > "$TMPDIR/code/server.py" <<'PYTHON'
from http.server import HTTPServer, BaseHTTPRequestHandler

class Handler(BaseHTTPRequestHandler):
    def do_GET(self):
        self.send_response(200)
        self.end_headers()
        self.wfile.write(b'hello from e2e test v2 (updated)')

HTTPServer(('', 8080), Handler).serve_forever()
PYTHON

(cd "$TMPDIR/code" && zip -q "$TMPDIR/code.zip" server.py)

echo "Updating function '$FUNCTION_NAME' (port $FUNCTION_PORT -> $UPDATED_PORT)..."
RESPONSE=$(curl -s -w "\n%{http_code}" \
    -X POST \
    -H "Authorization: $AUTH" \
    -F "name=$FUNCTION_NAME" \
    -F "tag=$FUNCTION_NAME:v2" \
    -F "runtime=$RUNTIME_NAME" \
    -F "port=$UPDATED_PORT" \
    -F "volume=" \
    -F "file=@$TMPDIR/code.zip" \
    "$BASE_URL/function/upload")

CODE=$(echo "$RESPONSE" | tail -1)
BODY=$(echo "$RESPONSE" | sed '$d')

check_response "$CODE" 200 "Update function"
echo "$BODY"

echo ""
echo "Old container should be torn down, new build started."
echo "Use get_function.sh to verify the port changed to $UPDATED_PORT."
