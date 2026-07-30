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
cat > main.go << 'GO'
package main
import (_ "t/dreego/gen"; core "codeberg.org/dreego/dreego/core")
func main() { core.Listen(":8080") }
GO
go run codeberg.org/dreego/dreego/cmd/dreego generate 2>&1
go build -o /tmp/srv .
/tmp/srv &
PID=$!
trap "kill $PID 2>/dev/null" EXIT
sleep 1
J=$(curl -s -H "Accept: application/json" http://localhost:8080/)
echo "$J" | grep -q '"format":"json"' || { echo "FAIL: JSON broken, got: $J"; exit 1; }
X=$(curl -s -H "Accept: application/xml" http://localhost:8080/)
echo "$X" | grep -q '<user>' || { echo "FAIL: XML broken, got: $X"; exit 1; }
H=$(curl -s http://localhost:8080/)
echo "$H" | grep -q '<h1>Lukas</h1>' || { echo "FAIL: HTML broken, got: $H"; exit 1; }
echo ok
