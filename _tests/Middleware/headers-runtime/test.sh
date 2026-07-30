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
func main() { core.Listen(":8080") }
GO

mkdir -p dreego/routes

cat > dreego/routes/get.dreego << 'DREEGO'
<div><p>hello</p></div>
DREEGO

go run codeberg.org/dreego/dreego/cmd/dreego generate 2>&1
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
