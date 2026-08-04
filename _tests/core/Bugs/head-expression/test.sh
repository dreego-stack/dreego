#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that expressions in a route <head> section are resolved at runtime
set -e

realrepo="$(cd "$(dirname "$0")"/../../../.. && pwd)"
workdir="$(mktemp -d)"
trap "rm -rf $workdir" EXIT INT TERM

cd "$workdir"

apk add --no-cache curl >/dev/null 2>&1 || true

port="${DREEGO_PORT:-$(( ( $(od -An -N2 -i /dev/urandom | tr -d ' ') % 50000 ) + 10000 ))}"

cat > go.mod << EOF
module t
go 1.22
require codeberg.org/dreego/dreego/core v0.0.0
replace codeberg.org/dreego/dreego/core => $realrepo/core
EOF

mkdir -p dreego/routes

cat > dreego/routes/get.dreego << 'DREEGO'
<go>doc := struct{ Title string }{Title: "My Docs Title"}</go>
<head><title>{doc.Title}</title></head>
<div><h1>{doc.Title}</h1></div>
DREEGO

cat > main.go << GO
package main
import (_ "t/dreego/gen"; core "codeberg.org/dreego/dreego/core")
func main() { core.Listen(":$port") }
GO

$DREEGO_BIN generate 2>&1
go build -o "$workdir/srv" .
"$workdir/srv" &
PID=$!
trap "kill $PID 2>/dev/null; rm -rf $workdir" EXIT
for i in $(seq 1 30); do curl -s -o /dev/null http://localhost:$port/ && break; sleep 0.1; done

BODY=$(curl -s http://localhost:$port/)

echo "$BODY" | grep -q '<title>My Docs Title</title>' || { echo "FAIL: head expression not resolved, got: $BODY"; exit 1; }
echo "$BODY" | grep -q '{doc.Title}' && { echo "FAIL: unresolved {doc.Title} found in HTML: $BODY"; exit 1; }

echo ok