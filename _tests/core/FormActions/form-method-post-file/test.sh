#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that a form with g-action generates both GET and POST handlers
set -e
realrepo="$(cd "$(dirname "$0")"/../../../.. && pwd)"
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
cat > dreego/routes/post-login.dreego << 'DREEGO'
<go>
    type LoginForm struct {
        Email string
    }
    func Login(c core.Context, form LoginForm) error {
        return nil
    }
</go>
<form g-action="Login" method="post">
    <input name="email">
    <button>Login</button>
</form>
DREEGO
cat > main.go << 'GO'
package main
import (_ "t/dreego/gen"; core "codeberg.org/dreego/dreego/core")
func main() { core.Listen(":0") }
GO
go run $realrepo/cmd/dreego generate 2>&1
grep -q "HandleIndexGet" dreego/gen/routes.go || { echo "FAIL: GET handler missing"; exit 1; }
grep -q "HandleIndexPost" dreego/gen/routes.go || { echo "FAIL: POST handler missing"; exit 1; }
go build -o /dev/null .
echo ok
