#!/bin/sh
set -e
cd "$(dirname "$0")"
rm -rf dreego/* 2>/dev/null || true
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
import (_ "t/dreego/gen"; core "codeberg.org/dreego/dreego/dreego-core")
func main() { core.Listen(":0") }
GO
go mod init t >/dev/null 2>&1
go mod edit -replace codeberg.org/dreego/dreego=../../..
go mod edit -require codeberg.org/dreego/dreego@v0.0.0
dreego generate 2>&1
grep -q "core.ValidateForm" dreego/gen/routes.go && { echo "FAIL: ValidateForm found but should not exist"; exit 1; }
go build -o /dev/null .
echo ok
