#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that csrf_token cookie has SameSite set at runtime
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

cat > main.go << GO
package main
import (_ "t/dreego/gen"; dreego "github.com/dreego-stack/dreego/core")
func main() { dreego.SetSessionStore(dreego.NewCookieStore([]byte("secret-key-32-bytes!"))); dreego.Listen(":$port") }
GO

mkdir -p dreego/routes

cat > dreego/routes/get.dreego << 'DREEGO'
<div><p>csrf samesite</p></div>
DREEGO

$DREEGO_BIN generate
go build -o $workdir/srv .
$workdir/srv &
PID=$!
trap "kill $PID 2>/dev/null; rm -rf $workdir" EXIT
for i in $(seq 1 30); do curl -s -o /dev/null http://localhost:$port/ && break; sleep 0.1; done
COOKIES=$(curl -s -I http://localhost:$port/ | grep -i "set-cookie" || true)
echo "$COOKIES" | grep -qi "csrf_token" || { echo "FAIL: no csrf_token cookie"; exit 1; }
echo "$COOKIES" | grep -qi "SameSite" || { echo "FAIL: csrf cookie has no SameSite"; exit 1; }
echo ok