#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
set -e
realrepo="$(cd "$(dirname "$0")"/../../.. && pwd)"
workdir="$(mktemp -d)"
trap "rm -rf $workdir" EXIT

cd "$workdir"

port=$(awk 'BEGIN{srand();print int(rand()*50000)+10000}')

cat > go.mod << EOF
module t
go 1.22
require codeberg.org/dreego/dreego v0.0.0
replace codeberg.org/dreego/dreego => $realrepo
EOF

mkdir -p dreego/routes
BT='`'
cat > dreego/routes/post-login.dreego << DREEGO
<go>
    type LoginForm struct {
        Email string ${BT}validate:"required,email"${BT}
    }
    func Login(c core.Context, form LoginForm) error {
        return c.Redirect("/dashboard", 303)
    }
</go>
<form g-action="Login" method="post">
    <input name="email" type="email">
    <button type="submit">Login</button>
</form>
DREEGO
cat > main.go << 'GO'
package main
import (_ "t/dreego/gen"; core "codeberg.org/dreego/dreego/core")
func main() { core.SetCSRF(false); core.Listen(":8080") }
GO
sed -i "s/8080/$port/" main.go
go run codeberg.org/dreego/dreego/cmd/dreego generate 2>&1
go build -o /tmp/srv .
/tmp/srv &
PID=$!
trap "kill $PID 2>/dev/null" EXIT
sleep 1
CODE=$(curl -s -o /dev/null -w "%{http_code}" -H "Content-Type: application/x-www-form-urlencoded" -d "email=test@dreego.dev" http://localhost:$port/)
[ "$CODE" = "303" ] || { echo "FAIL: expected 303 redirect, got $CODE"; exit 1; }
echo ok
