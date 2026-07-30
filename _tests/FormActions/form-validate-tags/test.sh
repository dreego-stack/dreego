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
BT='`'
cat > dreego/routes/get-form.dreego << DREEGO
<go>
    type ValForm struct {
        Email string ${BT}validate:"required,email"${BT}
    }
    func DoVal(c core.Context, form ValForm) error {
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
import (_ "t/dreego/gen"; core "codeberg.org/dreego/dreego/core")
func main() { core.Listen(":0") }
GO
go run codeberg.org/dreego/dreego/cmd/dreego generate 2>&1
grep -q "core.ValidateForm" dreego/gen/routes.go || { echo "FAIL: ValidateForm not in generated code"; exit 1; }
go build -o /dev/null .
echo ok
