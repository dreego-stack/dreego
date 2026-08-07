#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: 404 error page with doctype renders doctype first (no scope div before), returns 404 status
set -e

realrepo="$(cd "$(dirname "$0")"/../../../.. && pwd)"
workdir="$(mktemp -d)"
trap "rm -rf $workdir" EXIT INT TERM

cd "$workdir"

if [ -z "$DREEGO_BIN" ]; then
    DREEGO_BIN="$workdir/.dreego-bin"
    (cd "$realrepo" && go build -o "$DREEGO_BIN" ./cmd/dreego) || { echo "FAIL: could not build dreego CLI"; exit 1; }
    export DREEGO_BIN
fi

apk add --no-cache curl >/dev/null 2>&1 || true

port="${DREEGO_PORT:-$(( ( $(od -An -N2 -i /dev/urandom | tr -d ' ') % 50000 ) + 10000 ))}"

cat > go.mod << EOF
module t
go 1.22
require codeberg.org/dreego/dreego/core v0.0.0
replace codeberg.org/dreego/dreego/core => $realrepo/core
EOF

mkdir -p dreego/routes

cat > dreego/routes/404.dreego << 'DREEGO'
<!doctype html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>Not Found</title>
</head>
<body>
    <div><p>Not Found</p></div>
</body>
</html>
DREEGO

cat > main.go << GO
package main
import (_ "t/dreego/gen"; dreego "codeberg.org/dreego/dreego/core")
func main() { dreego.Listen(":$port") }
GO

$DREEGO_BIN generate 2>&1
go build -o "$workdir/srv" .
"$workdir/srv" &
PID=$!
trap "kill $PID 2>/dev/null; rm -rf $workdir" EXIT
for i in $(seq 1 30); do curl -s -o /dev/null http://localhost:$port/ && break; sleep 0.1; done

BODY=$(curl -s -w '\n%{http_code}' http://localhost:$port/nonexistent)
CODE=$(printf '%s' "$BODY" | tail -n 1)
HTML=$(printf '%s' "$BODY" | sed '$d')

if [ "$CODE" -ne 404 ]; then
    echo "FAIL: expected HTTP 404, got $CODE"
    exit 1
fi

case "$HTML" in
    "<!doctype html>"*) ;;
    *)
        echo "FAIL: body must start with <!doctype html>, got: $(printf '%s' "$HTML" | head -c 80)"
        exit 1
        ;;
esac

echo "$HTML" | grep -q 'data-scope' && { echo "FAIL: scope div must not appear in 404 body: $HTML"; exit 1; }
echo "$HTML" | grep -q '<title>Not Found</title>' || { echo "FAIL: head content missing in 404 body: $HTML"; exit 1; }
echo "$HTML" | grep -q '<p>Not Found</p>' || { echo "FAIL: body content missing in 404 body: $HTML"; exit 1; }

echo ok
