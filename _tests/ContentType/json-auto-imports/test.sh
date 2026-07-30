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
import (_ "t/dreego/gen"; core "codeberg.org/dreego/dreego/core")
func main() { core.Listen(":0") }
GO
go run codeberg.org/dreego/dreego/cmd/dreego generate 2>&1
grep -q "c.JSON" dreego/gen/routes.go || { echo "FAIL: c.JSON not in generated code"; exit 1; }
grep -q "c.Bind" dreego/gen/routes.go || { echo "FAIL: c.Bind not in generated code"; exit 1; }
grep -q "application/json" dreego/gen/routes.go || { echo "FAIL: no application/json content-type"; exit 1; }
go build -o /dev/null .
echo ok
