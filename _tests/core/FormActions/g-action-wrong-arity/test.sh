#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that a wrong handler arity skips POST registration
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
cat > dreego/routes/get-fail.dreego << 'DREEGO'
<go>
    type BadForm struct {
        X string
    }
    func Bad(c dreego.Context) error {
        return nil
    }
</go>
<form g-action="Bad" method="post">
    <input name="x">
    <button>OK</button>
</form>
DREEGO
cat > main.go << 'GO'
package main
import (_ "t/dreego/gen"; dreego "github.com/dreego-stack/dreego/core")
func main() { dreego.Listen(":0") }
GO
$DREEGO_BIN generate 2>&1
grep -q "dreego.Register(\"POST\"" dreego/gen/routes.go && { echo "FAIL: POST handler registered for bad arity"; exit 1; }
go build -o /dev/null .
echo ok
