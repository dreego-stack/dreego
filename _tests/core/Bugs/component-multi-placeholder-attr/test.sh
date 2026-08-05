#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: <Card href="{a}-{b}"/> resolves both placeholders and compiles
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

# A component call passes a quoted attribute value containing TWO placeholders:
# href="{a}-{b}". attrVal/extractAttrValues must split this into concatenation of
# a and b. Currently the whole inner "a}-{b" is coded as a single broken Go
# expression, so the generated code does not compile.
cat > dreego/components/Card.dreego << 'DREEGO'
Component Card (url string)
<div><a href="{url}">go</a></div>
DREEGO

cat > dreego/routes/get.dreego << 'DREEGO'
<go>left := "x"; right := "y"</go>
<div><@Card href="{left}-{right}"/></div>
DREEGO

$DREEGO_BIN generate

# The generated route must resolve BOTH placeholders in the attribute. The raw
# placeholder sequence "left}-{right" must not survive as one broken expression.
if grep -q 'left}-{right' dreego/gen/routes.go; then
    echo "FAIL: {a}-{b} was coded as single broken expression a}-{b"
    grep -n 'Card(' dreego/gen/routes.go || true
    exit 1
fi

# The generated project must compile.
go build -o /dev/null . || { echo "FAIL: go build failed"; exit 1; }

echo ok
