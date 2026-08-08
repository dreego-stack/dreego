#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: 404 page is self-contained (no layout wrap, no scope div before doctype), own head/style/script render
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
require github.com/dreego-stack/dreego v0.0.0
replace github.com/dreego-stack/dreego => $realrepo
EOF

mkdir -p dreego/layouts dreego/routes

cat > dreego/layouts/default.dreego << 'DREEGO'
<head>
    <title>Layout Site</title>
    {#head}
</head>
<div><main>{#slot}</main></div>
DREEGO

cat > dreego/routes/404.dreego << 'DREEGO'
<!doctype html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>Not Found</title>
    <link rel="stylesheet" href="/err.css">
</head>
<body>
    <div><p>Not Found</p></div>
</body>
</html>
<style>
p { color: red; }
</style>
<script>
console.log("err");
</script>
DREEGO

cat > main.go << GO
package main
import (_ "t/dreego/gen"; dreego "github.com/dreego-stack/dreego/core")
func main() { dreego.Listen(":$port") }
GO

$DREEGO_BIN generate 2>&1
go build -o "$workdir/srv" .
"$workdir/srv" &
PID=$!
trap "kill $PID 2>/dev/null; rm -rf $workdir" EXIT
for i in $(seq 1 30); do curl -s -o /dev/null http://localhost:$port/ && break; sleep 0.1; done

BODY=$(curl -s -w '\n%{http_code}' http://localhost:$port/nope)
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
echo "$HTML" | grep -q '<title>Not Found</title>' || { echo "FAIL: 404 head content missing: $HTML"; exit 1; }
echo "$HTML" | grep -q 'Layout Site' && { echo "FAIL: 404 page must not be wrapped in the layout (layout title leaked): $HTML"; exit 1; }
echo "$HTML" | grep -q 'href="/err.css"' || { echo "FAIL: 404 head CSS link missing: $HTML"; exit 1; }
echo "$HTML" | grep -q '<style>' || { echo "FAIL: 404 style section missing: $HTML"; exit 1; }
echo "$HTML" | grep -q 'color: red' || { echo "FAIL: 404 scoped style content missing: $HTML"; exit 1; }
echo "$HTML" | grep -q '<script>' || { echo "FAIL: 404 script section missing: $HTML"; exit 1; }
echo "$HTML" | grep -q 'console.log("err")' || { echo "FAIL: 404 script content missing: $HTML"; exit 1; }

echo ok
