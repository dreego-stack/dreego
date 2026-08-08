#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that JSON routes with auto-imports compile
set -e

realrepo="$(cd "$(dirname "$0")"/../../../.. && pwd)"
workdir="$(mktemp -d)"
trap "rm -rf $workdir" EXIT

cd "$workdir"

cat > go.mod << EOF
module t
go 1.22
require github.com/dreego-stack/dreego/core v0.0.0
replace github.com/dreego-stack/dreego/core => $realrepo/core
EOF

mkdir -p dreego/routes
cat > dreego/routes/post.dreego << 'DREEGO'
<go type="json">
    var input map[string]any
    c.Bind(&input)
    input["echo"] = true
    c.JSON(200, input)
</go>
DREEGO
cat > main.go << 'GO'
package main
import (_ "t/dreego/gen"; dreego "github.com/dreego-stack/dreego/core")
func main() { dreego.Listen(":0") }
GO
$DREEGO_BIN generate 2>&1
grep -q "c.JSON" dreego/gen/routes.go || { echo "FAIL: c.JSON not in generated code"; exit 1; }
grep -q "c.Bind" dreego/gen/routes.go || { echo "FAIL: c.Bind not in generated code"; exit 1; }
grep -q "application/json" dreego/gen/routes.go || { echo "FAIL: no application/json content-type"; exit 1; }
go build -o /dev/null .
echo ok
