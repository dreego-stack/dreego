#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that form validation is skipped without validate struct tags
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
cat > dreego/routes/get-form.dreego << 'DREEGO'
<go>
    type NoValForm struct {
        Email string
    }
    func NoVal(c core.Context, form NoValForm) error {
        return nil
    }
</go>
<form g-action="NoVal" method="post">
    <input name="email">
    <button>OK</button>
</form>
DREEGO
cat > main.go << 'GO'
package main
import (_ "t/dreego/gen"; core "codeberg.org/dreego/dreego/core")
func main() { core.Listen(":0") }
GO
$DREEGO_BIN generate 2>&1
grep -q "core.ValidateForm" dreego/gen/routes.go && { echo "FAIL: ValidateForm found but should not exist"; exit 1; }
go build -o /dev/null .
echo ok
