#!/bin/sh
set -e
cd "$(dirname "$0")"
rm -rf dreego/* 2>/dev/null || true
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
import (_ "t/dreego/gen"; core "codeberg.org/dreego/dreego/dreego-core")
func main() { core.Listen(":0") }
GO
go mod init t >/dev/null 2>&1
go mod edit -replace codeberg.org/dreego/dreego=../../..
go mod edit -require codeberg.org/dreego/dreego@v0.0.0
dreego generate 2>&1
grep -q "core.Register(\"POST\"" dreego/gen/routes.go && { echo "FAIL: POST handler registered without g-action"; exit 1; }
go build -o /dev/null .
echo ok
