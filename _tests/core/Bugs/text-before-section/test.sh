#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: <go> section after leading template text is parsed as Go code, never
#       silently lands as HTML text (Item 8, option a)
set -e

realrepo="$(cd "$(dirname "$0")"/../../../.. && pwd)"
workdir="$(mktemp -d)"
trap "rm -rf $workdir" EXIT

cd "$workdir"

if [ -z "$DREEGO_BIN" ]; then
    DREEGO_BIN="$workdir/.dreego-bin"
    (cd "$realrepo" && go build -o "$DREEGO_BIN" ./cmd/dreego) || { echo "FAIL: could not build dreego CLI"; exit 1; }
    export DREEGO_BIN
fi

cat > go.mod << EOF
module t
go 1.22
require github.com/dreego-stack/dreego/core v0.0.0
replace github.com/dreego-stack/dreego/core => $realrepo/core
EOF

cat > main.go << 'GO'
package main
import _ "t/dreego/gen"
func main() {}
GO

mkdir -p dreego/routes

# 1. Route starting with <!doctype html> followed by a <go> block: generate
#    and build must succeed, and the <go> body must be emitted as Go code,
#    not as template text. (A closing </html> after the <div> template is a
#    separate pre-existing parser behavior, pinned by TestParseExpectedSectionTag.)
cat > dreego/routes/get.dreego << 'DREEGO'
<!doctype html>
<html lang="en">
<go>msg := "hi"</go>
<div><p>{msg}</p></div>
DREEGO

$DREEGO_BIN generate

if grep -q 'html.EscapeString("msg := \\"hi\\""' dreego/gen/routes.go; then
    echo "FAIL: <go> body must not be emitted as template text"
    grep -n 'msg := "hi"' dreego/gen/routes.go || true
    exit 1
fi

if ! grep -q 'msg := "hi"' dreego/gen/routes.go; then
    echo "FAIL: <go> body must be emitted as Go code"
    exit 1
fi

go build -o /dev/null . || { echo "FAIL: go build failed"; exit 1; }

# 2. Control: <go> block at the file start still generates and compiles.
cat > dreego/routes/get.dreego << 'DREEGO'
<go>msg := "hi"</go>
<div><p>{msg}</p></div>
DREEGO

$DREEGO_BIN generate

if grep -q 'html.EscapeString("msg := \\"hi\\""' dreego/gen/routes.go; then
    echo "FAIL: <go> body must not be emitted as template text"
    grep -n 'msg := "hi"' dreego/gen/routes.go || true
    exit 1
fi

go build -o /dev/null . || { echo "FAIL: go build failed"; exit 1; }

echo ok
