#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that form validation tags generate validation code
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
BT='`'
cat > dreego/routes/get-form.dreego << DREEGO
<go>
    type ValForm struct {
        Email string ${BT}validate:"required,email"${BT}
    }
    func DoVal(c dreego.Context, form ValForm) error {
        return nil
    }
</go>
<form g-action="DoVal" method="post">
    <input name="email">
    <button>OK</button>
</form>
DREEGO
cat > main.go << 'GO'
package main
import (_ "t/dreego/gen"; dreego "github.com/dreego-stack/dreego/core")
func main() { dreego.Listen(":0") }
GO
$DREEGO_BIN generate 2>&1
grep -q "dreego.ValidateForm" dreego/gen/routes.go || { echo "FAIL: ValidateForm not in generated code"; exit 1; }
go build -o /dev/null .
echo ok
