#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: scoped <style> keeps complex declarations (radial-gradient with commas)
set -e

realrepo="$(cd "$(dirname "$0")"/../../../.. && pwd)"
workdir="$(mktemp -d)"
trap "rm -rf $workdir" EXIT

cd "$workdir"

if [ -z "$DREEGO_BIN" ]; then
    DREEGO_BIN="$workdir/.dreego-bin"
    (cd "$realrepo" && go build -o "$DREEGO_BIN" ./cli/dreego) || { echo "FAIL: could not build dreego CLI"; exit 1; }
    export DREEGO_BIN
fi

cat > go.mod << EOF
module t
go 1.22
require github.com/dreego-stack/dreego v0.0.0
replace github.com/dreego-stack/dreego => $realrepo
EOF

cat > main.go << 'GO'
package main
import _ "t/dreego/gen"
func main() {}
GO

mkdir -p dreego/components dreego/routes
cat > dreego/components/Grid.dreego << 'DREEGO'
Component Grid ()
<div><div class="bg-grid"><p>scoped</p></div></div>
<style>.bg-grid { background-image: radial-gradient(circle, #ccfbf1 1px, transparent 1px); }</style>
DREEGO
cat > dreego/routes/get.dreego << 'DREEGO'
<div><@Grid/></div>
DREEGO

$DREEGO_BIN generate
go build -o /dev/null .

generated="dreego/gen/components.go"
if ! grep -q 'radial-gradient' "$generated"; then
    echo "FAIL: radial-gradient dropped by scopeCSS"
    exit 1
fi
if ! grep -q '#ccfbf1' "$generated"; then
    echo "FAIL: first gradient color dropped by scopeCSS"
    exit 1
fi
if ! grep -q 'transparent 1px' "$generated"; then
    echo "FAIL: second gradient color dropped by scopeCSS"
    exit 1
fi
if ! grep -q '\[data-scope=' "$generated"; then
    echo "FAIL: scoped CSS prefix missing"
    exit 1
fi
echo ok
