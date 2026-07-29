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
import (_ "t/dreego/gen"; core "codeberg.org/dreego/dreego/dreego-core")
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
trap "kill $PID 2>/dev/null" EXIT
sleep 1
# send a custom request ID and verify it's echoed back
CUSTOM="abc123def456789a"
ID=$(curl -s -D- -H "X-Request-ID: $CUSTOM" http://localhost:8080/health 2>/dev/null | grep -i "x-request-id" | tr -d '\r' | awk '{print $2}')
[ "$ID" = "$CUSTOM" ] || { echo "FAIL: expected $CUSTOM, got $ID"; exit 1; }
echo ok
