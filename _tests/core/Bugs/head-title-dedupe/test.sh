#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Layout title + route title merge dedupes to exactly one title (route wins)
set -e

realrepo="$(cd "$(dirname "$0")"/../../../.. && pwd)"
workdir="$(mktemp -d)"
trap "rm -rf $workdir" EXIT INT TERM

cd "$workdir"

if [ -z "$DREEGO_BIN" ]; then
    DREEGO_BIN="$workdir/.dreego-bin"
    (cd "$realrepo" && go build -o "$DREEGO_BIN" ./cli/dreego) || { echo "FAIL: could not build dreego CLI"; exit 1; }
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
    <title>Site</title>
    <meta name="description" content="site desc">
    {#head}
</head>
<div><main>{#slot}</main></div>
DREEGO

cat > dreego/routes/get.dreego << 'DREEGO'
<head><title>Page</title><meta name="description" content="route desc"></head>
<div><h1>Page</h1></div>
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

BODY=$(curl -s http://localhost:$port/)

TITLE_COUNT=$(printf '%s' "$BODY" | grep -c '<title>')
if [ "$TITLE_COUNT" -ne 1 ]; then
    echo "FAIL: expected exactly 1 <title>, got $TITLE_COUNT: $BODY"
    exit 1
fi

echo "$BODY" | grep -q '<title>Page</title>' || { echo "FAIL: route title missing in body: $BODY"; exit 1; }
echo "$BODY" | grep -q 'Site' && { echo "FAIL: layout title still present in body: $BODY"; exit 1; }

DESC_COUNT=$(printf '%s' "$BODY" | grep -c 'name="description"')
if [ "$DESC_COUNT" -ne 1 ]; then
    echo "FAIL: expected exactly 1 meta description, got $DESC_COUNT: $BODY"
    exit 1
fi

echo "$BODY" | grep -q 'content="route desc"' || { echo "FAIL: route meta description missing in body: $BODY"; exit 1; }
echo "$BODY" | grep -q 'site desc' && { echo "FAIL: layout meta description still present in body: $BODY"; exit 1; }

echo ok
