#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: scoped <style> preserves @keyframes name and body, does not scope inside it
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
cat > dreego/components/Spin.dreego << 'DREEGO'
Component Spin ()
<div><div class="spinner"><p>scoped</p></div></div>
<style>
@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
.spinner { animation: spin 1s linear infinite; }
</style>
DREEGO
cat > dreego/routes/get.dreego << 'DREEGO'
<div><@Spin/></div>
DREEGO

$DREEGO_BIN generate
go build -o /dev/null .

generated="dreego/gen/components.go"
# @keyframes name and its body must survive
if ! grep -q '@keyframes spin' "$generated"; then
    echo "FAIL: @keyframes name missing"
    exit 1
fi
if ! grep -q 'transform: rotate(360deg)' "$generated"; then
    echo "FAIL: @keyframes body dropped by scopeCSS"
    exit 1
fi
# Keyframe from/to steps are not valid CSS selectors to scope: they must be unscoped
if grep -q '\[data-scope=[^ ]* *from' "$generated"; then
    echo "FAIL: @keyframes 'from' step wrongly scoped"
    exit 1
fi
if grep -q '\[data-scope=[^ ]* *to' "$generated"; then
    echo "FAIL: @keyframes 'to' step wrongly scoped"
    exit 1
fi
# The regular selector must still be scoped
if ! grep -q '\[data-scope=' "$generated"; then
    echo "FAIL: scoped CSS prefix missing"
    exit 1
fi
echo ok
