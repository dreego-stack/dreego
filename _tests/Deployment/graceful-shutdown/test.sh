#!/bin/sh
set -e
cd "$(dirname "$0")"
rm -rf dreego/* 2>/dev/null || true
mkdir -p dreego/routes
cat > dreego/routes/get.dreego << 'DREEGO'
<div><h1>hello</h1></div>
DREEGO
cat > main.go << 'GO'
package main
import (_ "t/dreego/gen"; core "codeberg.org/dreego/dreego/core")
func main() { core.SetLogging(false); core.Listen(":8080") }
GO
go mod init t >/dev/null 2>&1
go mod edit -replace codeberg.org/dreego/dreego=../../..
go mod edit -require codeberg.org/dreego/dreego@v0.0.0
sed -i 's|_ "gen"|_ "t/dreego/gen"|' main.go
dreego generate 2>&1
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
