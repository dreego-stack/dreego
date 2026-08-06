#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that a basic g-action form handler compiles
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
cat > dreego/routes/get-login.dreego << 'DREEGO'
<go>
    type LoginForm struct {
        Email string
    }
    func Login(c dreego.Context, form LoginForm) error {
        return nil
    }
</go>
<form g-action="Login" method="post">
    <input name="email" type="email">
    <button type="submit">Login</button>
</form>
DREEGO
cat > main.go << 'GO'
package main
import (_ "t/dreego/gen"; dreego "codeberg.org/dreego/dreego/core")
func main() { dreego.Listen(":0") }
GO
$DREEGO_BIN generate 2>&1
go build -o /dev/null .
echo ok
