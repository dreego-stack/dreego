#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that a plain POST form without g-action skips POST registration
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
cat > dreego/routes/get-plain.dreego << 'DREEGO'
<go>
    email := c.FormValue("email")
    c.Set("email", email)
</go>
<form method="post">
    <input name="email">
    <button>Submit</button>
</form>
DREEGO
cat > main.go << 'GO'
package main
import (_ "t/dreego/gen"; core "codeberg.org/dreego/dreego/core")
func main() { core.Listen(":0") }
GO
$DREEGO_BIN generate 2>&1
grep -q "core.Register(\"POST\"" dreego/gen/routes.go && { echo "FAIL: POST handler registered without g-action"; exit 1; }
go build -o /dev/null .
echo ok
