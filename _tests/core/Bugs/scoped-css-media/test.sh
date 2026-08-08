#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that scoped <style> preserves @media queries (B3)
set -e

realrepo="$(cd "$(dirname "$0")"/../../../.. && pwd)"
workdir="$(mktemp -d)"
trap "rm -rf $workdir" EXIT

cd "$workdir"

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
cat > dreego/components/Box.dreego << 'DREEGO'
Component Box ()
<div><div class="box"> scoped </div></div>
<style>.box { color: red; }
@media (min-width: 600px) {
  .box { color: blue; }
}
.box:hover { color: green; }</style>
DREEGO
cat > dreego/routes/get.dreego << 'DREEGO'
<div><@Box/></div>
DREEGO

$DREEGO_BIN generate
go build -o /dev/null .

# Verify @media is preserved and scoped inside the generated component
generated="dreego/gen/components.go"
if ! grep -q '@media (min-width: 600px)' "$generated"; then
    echo "FAIL: @media query dropped by scopeCSS (B3)"
    exit 1
fi
if ! grep -q '\[data-scope=' "$generated"; then
    echo "FAIL: scoped CSS prefix missing (B3)"
    exit 1
fi
echo ok
