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
    msg := "hello"
</go>
<div><h1>{msg}</h1></div>
DREEGO
cat > main.go << 'GO'
package main
import (_ "t/dreego/gen"; core "codeberg.org/dreego/dreego/core")
func main() { core.Listen(":0") }
GO
go run $realrepo/cmd/dreego generate 2>&1
grep -q "text/html" dreego/gen/routes.go || { echo "FAIL: no text/html content-type"; exit 1; }
grep -q "b.WriteString" dreego/gen/routes.go || { echo "FAIL: template rendering missing"; exit 1; }
go build -o /dev/null .
echo ok
