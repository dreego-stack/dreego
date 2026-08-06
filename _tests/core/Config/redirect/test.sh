#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that redirects configured in config.json work at runtime
set -e

realrepo="$(cd "$(dirname "$0")"/../../../.. && pwd)"
workdir="$(mktemp -d)"
trap "rm -rf $workdir" EXIT INT TERM

cd "$workdir"

apk add --no-cache curl >/dev/null 2>&1 || true

port="${DREEGO_PORT:-$(( ( $(od -An -N2 -i /dev/urandom | tr -d ' ') % 50000 ) + 10000 ))}"

cat > go.mod << EOF
module t
go 1.22
require codeberg.org/dreego/dreego/core v0.0.0
replace codeberg.org/dreego/dreego/core => $realrepo/core
EOF

mkdir -p dreego/routes

cat > dreego/config.json << 'JSON'
{
    "redirects": [
        { "from": "/old", "to": "/new", "status": 301 }
    ]
}
JSON

cat > dreego/routes/get.dreego << 'DREEGO'
<div><p>new page</p></div>
DREEGO

cat > main.go << GO
package main
import (_ "t/dreego/gen"; dreego "codeberg.org/dreego/dreego/core")
func main() { dreego.Listen(":$port") }
GO

$DREEGO_BIN generate 2>&1
go build -o "$workdir/srv" .
"$workdir/srv" &
PID=$!
trap "kill $PID 2>/dev/null; rm -rf $workdir" EXIT
for i in $(seq 1 30); do curl -s -o /dev/null http://localhost:$port/ && break; sleep 0.1; done

STATUS=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:$port/old)
[ "$STATUS" = "301" ] || { echo "FAIL: expected 301 redirect, got $STATUS"; exit 1; }

LOCATION=$(curl -s -o /dev/null -w "%{redirect_url}" http://localhost:$port/old)
echo "$LOCATION" | grep -q "/new$" || { echo "FAIL: expected redirect to /new, got $LOCATION"; exit 1; }

echo ok
