#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
set -e

realrepo="$(cd "$(dirname "$0")"/../../.. && pwd)"
workdir="$(mktemp -d)"
trap "rm -rf $workdir" EXIT

cd "$workdir"

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

mkdir -p dreego/routes

cat > dreego/routes/get.dreego << 'DREEGO'
<div><h1>hello</h1></div>
DREEGO

go run codeberg.org/dreego/dreego/cmd/dreego generate 2>&1
go build -o /tmp/srv .
/tmp/srv &
PID=$!
trap "kill $PID 2>/dev/null" EXIT
sleep 1
# send a custom request ID and verify it's echoed back
CUSTOM="abc123def456789a"
ID=$(curl -s -D- -H "X-Request-ID: $CUSTOM" http://localhost:8080/health 2>/dev/null | grep -i "x-request-id" | tr -d '\r' | awk '{print $2}')
[ "$ID" = "$CUSTOM" ] || { echo "FAIL: expected $CUSTOM, got $ID"; exit 1; }
echo ok
