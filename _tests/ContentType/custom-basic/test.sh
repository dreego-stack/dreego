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
<go type="custom">
    msg := []byte("hello world")
    c.Write(200, "text/plain", msg)
</go>
DREEGO
cat > main.go << 'GO'
package main
import (_ "t/dreego/gen"; core "codeberg.org/dreego/dreego/core")
func main() { core.Listen(":0") }
GO
go run $realrepo/cmd/dreego generate 2>&1
grep -q "text/plain" dreego/gen/routes.go || { echo "FAIL: no text/plain content-type"; exit 1; }
grep -q "c.Write" dreego/gen/routes.go || { echo "FAIL: c.Write not in generated code"; exit 1; }
go build -o /dev/null .
echo ok
