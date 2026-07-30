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

mkdir -p dreego/routes
cat > dreego/routes/get.dreego << 'DREEGO'
<div><h1>hello</h1></div>
DREEGO
cat > main.go << 'GO'
package main
import (_ "t/dreego/gen"; core "codeberg.org/dreego/dreego/core")
func main() { core.SetLogging(false); core.Listen(":8080") }
GO
go run codeberg.org/dreego/dreego/cmd/dreego generate 2>&1
go build -o /tmp/srv .
/tmp/srv &
PID=$!
sleep 1
# verify server is running
RESP=$(curl -s http://localhost:8080/health)
[ "$RESP" = "ok" ] || { echo "FAIL: health check failed, got: $RESP"; exit 1; }
# send SIGTERM and verify graceful shutdown
kill -TERM $PID 2>/dev/null
wait $PID 2>/dev/null
EXIT_CODE=$?
# http.Server.Shutdown returns nil on graceful close
[ $EXIT_CODE -eq 0 ] || { echo "FAIL: server exit code $EXIT_CODE"; exit 1; }
echo ok
