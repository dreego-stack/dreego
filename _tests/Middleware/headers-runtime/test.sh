#!/bin/sh
set -e
cd "$(dirname "$0")"
rm -rf dreego/* 2>/dev/null || true
mkdir -p dreego/routes
cat > dreego/routes/get.dreego << 'DREEGO'
<div><p>hello</p></div>
DREEGO
cat > main.go << 'GO'
package main
import (_ "t/dreego/gen"; core "codeberg.org/dreego/dreego/core")
func main() { core.Listen(":8080") }
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
HEADERS=$(curl -s -I http://localhost:8080/)
echo "$HEADERS" | grep -q "X-Content-Type-Options: nosniff" || { echo "FAIL: X-Content-Type-Options missing"; exit 1; }
echo "$HEADERS" | grep -q "X-Frame-Options: DENY" || { echo "FAIL: X-Frame-Options missing"; exit 1; }
echo "$HEADERS" | grep -q "Referrer-Policy: strict-origin-when-cross-origin" || { echo "FAIL: Referrer-Policy missing"; exit 1; }
echo ok
