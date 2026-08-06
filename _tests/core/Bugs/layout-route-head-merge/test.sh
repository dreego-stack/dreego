#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Route head is merged into layout {#head}
set -e

realrepo="$(cd "$(dirname "$0")"/../../../.. && pwd)"
workdir="$(mktemp -d)"
trap "rm -rf $workdir" EXIT INT TERM

cd "$workdir"

if [ -z "$DREEGO_BIN" ]; then
    DREEGO_BIN="$workdir/.dreego-bin"
    (cd "$realrepo" && go build -o "$DREEGO_BIN" ./cmd/dreego) || { echo "FAIL: could not build dreego CLI"; exit 1; }
    export DREEGO_BIN
fi

apk add --no-cache curl >/dev/null 2>&1 || true

port="${DREEGO_PORT:-$(( ( $(od -An -N2 -i /dev/urandom | tr -d ' ') % 50000 ) + 10000 ))}"

cat > go.mod << EOF
module t
go 1.22
require codeberg.org/dreego/dreego/core v0.0.0
replace codeberg.org/dreego/dreego/core => $realrepo/core
EOF

mkdir -p dreego/layouts dreego/routes

cat > dreego/layouts/default.dreego << 'DREEGO'
<html>
<head>
<meta charset="utf-8">
{#head}
</head>
<body><main>{#slot}</main></body>
</html>
DREEGO

cat > dreego/routes/get.dreego << 'DREEGO'
<head><title>Merged Title</title><script src="route-script.js"></script></head>
<div><p>page content</p></div>
DREEGO

cat > main.go << GO
package main
import (_ "t/dreego/gen"; dreego "codeberg.org/dreego/dreego/core")
func main() { dreego.Listen(":$port") }
GO

$DREEGO_BIN generate 2>&1
go build -o "$workdir/srv" .
"$workdir/srv" &
PID=$!
trap "kill $PID 2>/dev/null; rm -rf $workdir" EXIT
for i in $(seq 1 30); do curl -s -o /dev/null http://localhost:$port/ && break; sleep 0.1; done

BODY=$(curl -s http://localhost:$port/)

echo "$BODY" | grep -q '<meta charset="utf-8">' || { echo "FAIL: layout head meta missing in body: $BODY"; exit 1; }
echo "$BODY" | grep -q '<title>Merged Title</title>' || { echo "FAIL: route head title missing in body: $BODY"; exit 1; }
echo "$BODY" | grep -q 'route-script.js' || { echo "FAIL: route head script missing in body: $BODY"; exit 1; }

echo ok
