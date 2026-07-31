#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that JSON and HTML routes sharing variables compile
set -e

realrepo="$(cd "$(dirname "$0")"/../../../.. && pwd)"
workdir="$(mktemp -d)"
trap "rm -rf $workdir" EXIT

cd "$workdir"

cat > go.mod << EOF
module t
go 1.22
require codeberg.org/dreego/dreego/core v0.0.0
replace codeberg.org/dreego/dreego/core => $realrepo/core
EOF

mkdir -p dreego/routes
cat > dreego/routes/get.dreego << 'DREEGO'
<go>
    msg := "Lukas"
</go>

<go type="json">
    c.JSON(200, map[string]string{"name": msg})
</go>

<div>
    <h1>{msg}</h1>
</div>
DREEGO
cat > main.go << 'GO'
package main
import (_ "t/dreego/gen"; core "codeberg.org/dreego/dreego/core")
func main() { core.Listen(":0") }
GO
$DREEGO_BIN generate 2>&1
grep -q "application/json" dreego/gen/routes.go || { echo "FAIL: no application/json content-type"; exit 1; }
grep -q "c.JSON" dreego/gen/routes.go || { echo "FAIL: c.JSON not in generated code"; exit 1; }
grep -q "Lukas" dreego/gen/routes.go || { echo "FAIL: shared var not in output"; exit 1; }
go build -o /dev/null .
echo ok
