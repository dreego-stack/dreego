#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that a JSON route falls back to HTML without an Accept header
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
require github.com/dreego-stack/dreego v0.0.0
replace github.com/dreego-stack/dreego => $realrepo
EOF

mkdir -p dreego/routes
cat > dreego/routes/get.dreego << 'DREEGO'
<go>
    msg := "hello"
</go>
<go type="json">
    c.JSON(200, map[string]string{"msg": msg})
</go>
<div><h1>{msg}</h1></div>
DREEGO
cat > main.go << GO
package main
import (_ "t/dreego/gen"; dreego "github.com/dreego-stack/dreego/core")
func main() { dreego.Listen(":$port") }
GO
$DREEGO_BIN generate 2>&1
go build -o $workdir/srv .
$workdir/srv &
PID=$!
trap "kill $PID 2>/dev/null; rm -rf $workdir" EXIT
for i in $(seq 1 30); do curl -s -o /dev/null http://localhost:$port/ && break; sleep 0.1; done
RESP=$(curl -s http://localhost:$port/)
echo "$RESP" | grep -q '<h1>hello</h1>' || { echo "FAIL: no HTML fallback, got: $RESP"; exit 1; }
echo ok
