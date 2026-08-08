#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that a multi-type route (JSON, XML, HTML) works at runtime
set -e

realrepo="$(cd "$(dirname "$0")"/../../../.. && pwd)"
workdir="$(mktemp -d)"
trap "rm -rf $workdir" EXIT

cd "$workdir"
apk add --no-cache curl >/dev/null 2>&1 || true

port="${DREEGO_PORT:-$(( ( $(od -An -N2 -i /dev/urandom | tr -d ' ') % 50000 ) + 10000 ))}"

cat > go.mod << EOF
module t
go 1.22
require github.com/dreego-stack/dreego/core v0.0.0
replace github.com/dreego-stack/dreego/core => $realrepo/core
EOF

mkdir -p dreego/routes
cat > dreego/routes/get.dreego << 'DREEGO'
<go>
    name := "Lukas"
</go>
<go type="json">
    c.JSON(200, map[string]string{"format": "json", "name": name})
</go>
<go type="xml">
    user := struct{XMLName struct{} `xml:"user"`; Name string `xml:"name"`}{Name: name}
    c.XML(200, user)
</go>
<div><h1>{name}</h1></div>
DREEGO
cat > main.go << GO
package main
import (_ "t/dreego/gen"; dreego "github.com/dreego-stack/dreego/core")
func main() { dreego.Listen(":$port") }
GO
$DREEGO_BIN generate 2>&1
go build -o $workdir/srv .
$workdir/srv &
PID=$!
trap "kill $PID 2>/dev/null; rm -rf $workdir" EXIT
for i in $(seq 1 30); do curl -s -o /dev/null http://localhost:$port/ && break; sleep 0.1; done
J=$(curl -s -H "Accept: application/json" http://localhost:$port/)
echo "$J" | grep -q '"format":"json"' || { echo "FAIL: JSON broken, got: $J"; exit 1; }
X=$(curl -s -H "Accept: application/xml" http://localhost:$port/)
echo "$X" | grep -q '<user>' || { echo "FAIL: XML broken, got: $X"; exit 1; }
H=$(curl -s http://localhost:$port/)
echo "$H" | grep -q '<h1>Lukas</h1>' || { echo "FAIL: HTML broken, got: $H"; exit 1; }
echo ok
