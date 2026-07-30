#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that auto-generated X-Request-ID is returned
set -e

realrepo="$(cd "$(dirname "$0")"/../../.. && pwd)"
workdir="$(mktemp -d)"
trap "rm -rf $workdir" EXIT

cd "$workdir"

port=$(awk 'BEGIN{srand();print int(rand()*50000)+10000}')

cat > go.mod << EOF
module t
go 1.22
require codeberg.org/dreego/dreego v0.0.0
replace codeberg.org/dreego/dreego => $realrepo
EOF

cat > main.go << 'GO'
package main
import (_ "t/dreego/gen"; core "codeberg.org/dreego/dreego/core")
func main() { core.SetLogging(false); core.Listen(":8080") }
GO
sed -i "s/8080/$port/" main.go

mkdir -p dreego/routes

cat > dreego/routes/get.dreego << 'DREEGO'
<div><h1>hello</h1></div>
DREEGO

go run $realrepo/cmd/dreego generate 2>&1
go build -o /tmp/srv .
/tmp/srv &
PID=$!
trap "kill $PID 2>/dev/null" EXIT
sleep 1
ID=$(curl -s -D- http://localhost:$port/health 2>/dev/null | grep -i "x-request-id" | tr -d '\r' | awk '{print $2}')
[ -z "$ID" ] && { echo "FAIL: no X-Request-ID header"; exit 1; }
# should be 16 hex chars
echo "$ID" | grep -qE '^[a-f0-9]{16}$' || { echo "FAIL: invalid request ID format: $ID"; exit 1; }
echo ok
