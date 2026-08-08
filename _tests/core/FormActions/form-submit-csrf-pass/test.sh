#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that a form submit with valid CSRF token succeeds
set -e
realrepo="$(cd "$(dirname "$0")"/../../../.. && pwd)"
workdir="$(mktemp -d)"
trap "rm -rf $workdir" EXIT

cd "$workdir"
apk add --no-cache curl >/dev/null 2>&1 || true

port="${DREEGO_PORT:-$(( ( $(od -An -N2 -i /dev/urandom | tr -d ' ') % 50000 ) + 10000 ))}"

cat > go.mod << EOF
module t
go 1.22
require github.com/dreego-stack/dreego v0.0.0
replace github.com/dreego-stack/dreego => $realrepo
EOF

mkdir -p dreego/routes
BT='`'
cat > dreego/routes/post-login.dreego << DREEGO
<go>
    type LoginForm struct {
        Email string ${BT}validate:"required"${BT}
    }
    func Login(c dreego.Context, form LoginForm) error {
        return c.Redirect("/ok", 303)
    }
</go>
<form g-action="Login" method="post">
    <input name="email" type="email">
    <button type="submit">Login</button>
</form>
DREEGO
cat > main.go << GO
package main
import (_ "t/dreego/gen"; dreego "github.com/dreego-stack/dreego/core")
func main() {
    dreego.SetSessionStore(dreego.NewCookieStore([]byte("test-secret")));
    dreego.Listen(":$port")
}
GO
$DREEGO_BIN generate 2>&1
go build -o $workdir/srv .
$workdir/srv &
PID=$!
trap "kill $PID 2>/dev/null; rm -rf $workdir" EXIT
for i in $(seq 1 60); do curl -s -o /dev/null http://localhost:$port/health && break; sleep 0.5; done
COOKIE_JAR="$workdir/cookies"
curl -s -c "$COOKIE_JAR" http://localhost:$port/health > /dev/null
CSRF=$(grep csrf_token "$COOKIE_JAR" | awk '{print $NF}')
CODE=$(curl -s -o /dev/null -w "%{http_code}" -b "$COOKIE_JAR" -H "Content-Type: application/x-www-form-urlencoded" -d "email=test@dreego.dev&csrf_token=$CSRF" http://localhost:$port/)
[ "$CODE" = "303" ] || { echo "FAIL: expected 303 with valid CSRF, got $CODE"; exit 1; }
echo ok
