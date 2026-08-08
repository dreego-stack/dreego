#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests the /ready endpoint with readiness flag
set -e

realrepo="$(cd "$(dirname "$0")"/../../../.. && pwd)"
workdir="$(mktemp -d)"
trap "rm -rf $workdir" EXIT

cd "$workdir"
apk add --no-cache curl >/dev/null 2>&1 || true

port="${DREEGO_PORT:-$(( ( $(od -An -N2 -i /dev/urandom | tr -d ' ') % 50000 ) + 10000 ))}"

cat > go.mod << EOF
module t
go 1.22
require github.com/dreego-stack/dreego/core v0.0.0
replace github.com/dreego-stack/dreego/core => $realrepo/core
EOF

cat > main.go << GO
package main
import (_ "t/dreego/gen"; dreego "github.com/dreego-stack/dreego/core")
func main() { dreego.SetReady(false); dreego.Listen(":$port") }
GO

mkdir -p dreego/routes

cat > dreego/routes/get.dreego << 'DREEGO'
<div><p>root</p></div>
DREEGO

$DREEGO_BIN generate 2>&1
go build -o $workdir/srv .
$workdir/srv &
PID=$!
trap "kill $PID 2>/dev/null; rm -rf $workdir" EXIT
for i in $(seq 1 30); do curl -s -o /dev/null http://localhost:$port/ && break; sleep 0.1; done
RESP=$(curl -s http://localhost:$port/ready)
[ "$RESP" = "not ready" ] || { echo "FAIL: /ready not returning not ready, got: $RESP"; exit 1; }
echo ok
