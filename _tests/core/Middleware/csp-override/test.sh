#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that CSP override via SetCSP applies at runtime
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
require codeberg.org/dreego/dreego/core v0.0.0
replace codeberg.org/dreego/dreego/core => $realrepo/core
EOF

cat > main.go << GO
package main
import (_ "t/dreego/gen"; core "codeberg.org/dreego/dreego/core")
func main() { core.SetCSP("default-src 'none'"); core.Listen(":$port") }
GO

mkdir -p dreego/routes

cat > dreego/routes/get.dreego << 'DREEGO'
<div><p>csp override</p></div>
DREEGO

$DREEGO_BIN generate
go build -o $workdir/srv .
$workdir/srv &
PID=$!
trap "kill $PID 2>/dev/null; rm -rf $workdir" EXIT
for i in $(seq 1 30); do curl -s -o /dev/null http://localhost:$port/ && break; sleep 0.1; done
HEADERS=$(curl -s -I http://localhost:$port/)
echo "$HEADERS" | grep -qi "default-src 'none'" || { echo "FAIL: custom CSP not applied"; exit 1; }
echo ok