#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: scoped <style> preserves comma selectors and nested parens
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
require codeberg.org/dreego/dreego/core v0.0.0
replace codeberg.org/dreego/dreego/core => $realrepo/core
EOF

cat > main.go << 'GO'
package main
import _ "t/dreego/gen"
func main() {}
GO

mkdir -p dreego/components dreego/routes
cat > dreego/components/Heading.dreego << 'DREEGO'
Component Heading ()
<div><h1>scoped</h1></div>
<style>h1, h2 { border-radius: calc(100% - 20px); color: rgb(1, 2, 3); }</style>
DREEGO
cat > dreego/routes/get.dreego << 'DREEGO'
<div><@Heading/></div>
DREEGO

$DREEGO_BIN generate
go build -o /dev/null .

generated="dreego/gen/components.go"
# Both comma selectors must survive, each scoped separately
if ! grep -q 'h1' "$generated"; then
    echo "FAIL: h1 selector missing"
    exit 1
fi
if ! grep -q 'h2' "$generated"; then
    echo "FAIL: h2 selector missing"
    exit 1
fi
# The declarations (with commas/parens) must be preserved
if ! grep -q 'calc(100% - 20px)' "$generated"; then
    echo "FAIL: calc() declaration dropped by scopeCSS"
    exit 1
fi
if ! grep -q 'rgb(1, 2, 3)' "$generated"; then
    echo "FAIL: rgb() declaration with commas dropped by scopeCSS"
    exit 1
fi
# Both selectors must be scoped
scoped_count=$(grep -c '\[data-scope=' "$generated" || true)
if [ "$scoped_count" -lt 2 ]; then
    echo "FAIL: expected at least 2 scoped selectors, got $scoped_count"
    exit 1
fi
echo ok
