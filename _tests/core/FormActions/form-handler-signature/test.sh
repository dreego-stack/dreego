#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that a wrong form handler signature causes a build failure
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
cat > dreego/routes/get-fail.dreego << 'DREEGO'
<go>
    type BadForm struct {
        X string
    }
    func bad(c dreego.Context, form BadForm) string {
        return "wrong"
    }
</go>
<form g-action="bad" method="post">
    <input name="x">
    <button>OK</button>
</form>
DREEGO
cat > main.go << 'GO'
package main
import (_ "t/dreego/gen"; dreego "codeberg.org/dreego/dreego/core")
func main() { dreego.Listen(":0") }
GO
$DREEGO_BIN generate 2>&1
go build -o /dev/null . 2>&1 && { echo "FAIL: should not build with wrong return type"; exit 1; }
echo ok
