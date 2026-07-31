#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that the server shuts down gracefully on SIGTERM
set -e
realrepo="$(cd "$(dirname "$0")"/../../../.. && pwd)"
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

mkdir -p dreego/routes
cat > dreego/routes/get.dreego << 'DREEGO'
<div><h1>hello</h1></div>
DREEGO
cat > main.go << 'GO'
package main
import (_ "t/dreego/gen"; core "codeberg.org/dreego/dreego/core")
func main() { core.SetLogging(false); core.Listen(":8080") }
GO
sed -i "s/8080/$port/" main.go
go run $realrepo/cmd/dreego generate 2>&1
go build -o $workdir/srv .
$workdir/srv &
PID=$!
for i in $(seq 1 30); do curl -s -o /dev/null http://localhost:$port/ && break; sleep 0.1; done
# verify server is running
RESP=$(curl -s http://localhost:$port/health)
[ "$RESP" = "ok" ] || { echo "FAIL: health check failed, got: $RESP"; exit 1; }
# send SIGTERM and verify graceful shutdown
kill -TERM $PID 2>/dev/null
wait $PID 2>/dev/null
EXIT_CODE=$?
# http.Server.Shutdown returns nil on graceful close
[ $EXIT_CODE -eq 0 ] || { echo "FAIL: server exit code $EXIT_CODE"; exit 1; }
echo ok
