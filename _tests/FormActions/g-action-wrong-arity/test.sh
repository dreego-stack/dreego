#!/bin/sh
set -e
cd "$(dirname "$0")"
rm -rf dreego/* 2>/dev/null || true
mkdir -p dreego/routes
cat > dreego/routes/get-fail.dreego << 'DREEGO'
<go>
    type BadForm struct {
        X string
    }
    func Bad(c core.Context) error {
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
import (_ "t/dreego/gen"; core "codeberg.org/dreego/dreego/core")
func main() { core.Listen(":0") }
GO
go mod init t >/dev/null 2>&1
go mod edit -replace codeberg.org/dreego/dreego=../../..
go mod edit -require codeberg.org/dreego/dreego@v0.0.0
dreego generate 2>&1
grep -q "core.Register(\"POST\"" dreego/gen/routes.go && { echo "FAIL: POST handler registered for bad arity"; exit 1; }
go build -o /dev/null .
echo ok
