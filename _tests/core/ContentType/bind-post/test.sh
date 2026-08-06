#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that JSON POST body binding works at runtime
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

mkdir -p dreego/routes
cat > dreego/routes/post.dreego << 'DREEGO'
<go type="json">
    var input map[string]any
    c.Bind(&input)
    input["echo"] = true
    c.JSON(200, input)
</go>
DREEGO
cat > main.go << GO
package main
import (_ "t/dreego/gen"; dreego "codeberg.org/dreego/dreego/core")
func main() { dreego.Listen(":$port") }
GO
$DREEGO_BIN generate 2>&1
go build -o $workdir/srv .
$workdir/srv &
PID=$!
trap "kill $PID 2>/dev/null; rm -rf $workdir" EXIT
for i in $(seq 1 30); do curl -s -o /dev/null http://localhost:$port/ && break; sleep 0.1; done
RESP=$(curl -s -H "Accept: application/json" -H "Content-Type: application/json" -d '{"name":"Lukas"}' http://localhost:$port/)
echo "$RESP" | grep -q '"name":"Lukas"' || { echo "FAIL: Bind not working, got: $RESP"; exit 1; }
echo "$RESP" | grep -q '"echo":true' || { echo "FAIL: no echo field, got: $RESP"; exit 1; }
echo ok
