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
cat > dreego/routes/post.dreego << 'DREEGO'
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
import (_ "t/dreego/gen"; core "codeberg.org/dreego/dreego/core")
func main() {
    core.SetCSRF(false)
    core.Listen(":8080")
}
GO
go run codeberg.org/dreego/dreego/cmd/dreego generate 2>&1
go build -o /tmp/srv .
/tmp/srv &
PID=$!
trap "kill $PID 2>/dev/null" EXIT
sleep 1
RESP=$(curl -s -o /dev/null -w "%{http_code}" -H "Content-Type: application/x-www-form-urlencoded" -d "email=hello@test.com" http://localhost:8080/)
[ "$RESP" = "200" ] || { echo "FAIL: expected 200 for plain POST, got $RESP"; exit 1; }
echo ok
