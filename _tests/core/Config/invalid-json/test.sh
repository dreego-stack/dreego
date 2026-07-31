#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that the server runs despite invalid config.json
set -e

realrepo="$(cd "$(dirname "$0")"/../../../.. && pwd)"
workdir="$(mktemp -d)"
trap "rm -rf $workdir" EXIT INT TERM

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

cat > dreego/config.json << 'JSON'
{ broken json !!!
JSON

cat > dreego/routes/get.dreego << 'DREEGO'
<div><p>hello</p></div>
DREEGO

cat > main.go << 'GO'
package main
import (_ "t/dreego/gen"; core "codeberg.org/dreego/dreego/core")
func main() { core.Listen(":8080") }
GO
sed -i "s/8080/$port/" main.go

go run $realrepo/cmd/dreego generate 2>&1
go build -o "$workdir/srv" .
"$workdir/srv" &
PID=$!
trap "kill $PID 2>/dev/null; rm -rf $workdir" EXIT
for i in $(seq 1 30); do curl -s -o /dev/null http://localhost:$port/ && break; sleep 0.1; done

STATUS=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:$port/)
[ "$STATUS" = "200" ] || { echo "FAIL: expected 200, got $STATUS"; exit 1; }

BODY=$(curl -s http://localhost:$port/)
echo "$BODY" | grep -q 'hello' || { echo "FAIL: expected 'hello' in response"; exit 1; }

echo ok
