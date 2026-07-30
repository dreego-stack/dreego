#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that a custom X-Request-ID header is echoed back
set -e

realrepo="$(cd "$(dirname "$0")"/../../.. && pwd)"
workdir="$(mktemp -d)"
trap "rm -rf $workdir" EXIT

cd "$workdir"
apk add --no-cache curl >/dev/null 2>&1 || true

port=$(od -An -N2 -i /dev/urandom | tr -d ' ')
port=$((port % 50000 + 10000))

cat > go.mod << EOF
module t
go 1.22
require codeberg.org/dreego/dreego v0.0.0
replace codeberg.org/dreego/dreego => $realrepo
EOF

cat > main.go << 'GO'
package main
import (_ "t/dreego/gen"; core "codeberg.org/dreego/dreego/core")
func main() { core.SetLogging(false); core.Listen(":8080") }
GO
sed -i "s/8080/$port/" main.go

mkdir -p dreego/routes

cat > dreego/routes/get.dreego << 'DREEGO'
<div><h1>hello</h1></div>
DREEGO

go run $realrepo/cmd/dreego generate 2>&1
go build -o $workdir/srv .
$workdir/srv &
PID=$!
trap "kill $PID 2>/dev/null; rm -rf $workdir" EXIT
for i in $(seq 1 30); do curl -s -o /dev/null http://localhost:$port/ && break; sleep 0.1; done
# send a custom request ID and verify it's echoed back
CUSTOM="abc123def456789a"
ID=$(curl -s -D- -H "X-Request-ID: $CUSTOM" http://localhost:$port/health 2>/dev/null | grep -i "x-request-id" | tr -d '\r' | awk '{print $2}')
[ "$ID" = "$CUSTOM" ] || { echo "FAIL: expected $CUSTOM, got $ID"; exit 1; }
echo ok
