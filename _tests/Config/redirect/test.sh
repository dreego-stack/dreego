#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
set -e

realrepo="$(cd "$(dirname "$0")"/../../.. && pwd)"
workdir="$(mktemp -d)"
trap "rm -rf $workdir" EXIT INT TERM

cd "$workdir"

apk add --no-cache curl >/dev/null 2>&1 || true

port=$(awk 'BEGIN{srand();print int(rand()*50000)+10000}')

cat > go.mod << EOF
module t
go 1.22
require codeberg.org/dreego/dreego v0.0.0
replace codeberg.org/dreego/dreego => $realrepo
EOF

mkdir -p dreego/routes

cat > dreego/config.json << 'JSON'
{
    "redirects": [
        { "from": "/old", "to": "/new", "status": 301 }
    ]
}
JSON

cat > dreego/routes/get.dreego << 'DREEGO'
<div><p>new page</p></div>
DREEGO

cat > main.go << 'GO'
package main
import (_ "t/dreego/gen"; core "codeberg.org/dreego/dreego/core")
func main() { core.Listen(":8080") }
GO
sed -i "s/8080/$port/" main.go

go run codeberg.org/dreego/dreego/cmd/dreego generate 2>&1
go build -o "$workdir/srv" .
"$workdir/srv" &
PID=$!
trap "kill $PID 2>/dev/null; rm -rf $workdir" EXIT
sleep 1

STATUS=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:$port/old)
[ "$STATUS" = "301" ] || { echo "FAIL: expected 301 redirect, got $STATUS"; exit 1; }

LOCATION=$(curl -s -o /dev/null -w "%{redirect_url}" http://localhost:$port/old)
echo "$LOCATION" | grep -q "/new$" || { echo "FAIL: expected redirect to /new, got $LOCATION"; exit 1; }

echo ok
