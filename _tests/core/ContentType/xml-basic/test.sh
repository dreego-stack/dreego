#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that a basic XML route generates correct code
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
<go type="xml">
    user := struct{XMLName struct{} `xml:"user"`; Name string `xml:"name"`}{Name: "Lukas"}
    c.XML(200, user)
</go>
DREEGO
cat > main.go << 'GO'
package main
import (_ "t/dreego/gen"; core "codeberg.org/dreego/dreego/core")
func main() { core.Listen(":0") }
GO
$DREEGO_BIN generate 2>&1
grep -q "application/xml" dreego/gen/routes.go || { echo "FAIL: no application/xml content-type"; exit 1; }
grep -q "c.XML" dreego/gen/routes.go || { echo "FAIL: c.XML not in generated code"; exit 1; }
go build -o /dev/null .
echo ok
