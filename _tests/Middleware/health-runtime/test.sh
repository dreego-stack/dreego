#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests the /health endpoint at runtime
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
func main() { core.Listen(":8080") }
GO
sed -i "s/8080/$port/" main.go

mkdir -p dreego/routes

cat > dreego/routes/get.dreego << 'DREEGO'
<div><p>root</p></div>
DREEGO

go run $realrepo/cmd/dreego generate 2>&1
go build -o /tmp/srv .
/tmp/srv &
PID=$!
trap "kill $PID 2>/dev/null" EXIT
sleep 1
RESP=$(curl -s http://localhost:$port/health)
[ "$RESP" = "ok" ] || { echo "FAIL: /health not ok, got: $RESP"; exit 1; }
echo ok
