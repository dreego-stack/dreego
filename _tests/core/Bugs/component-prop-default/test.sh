#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: <@Card/> without props uses the signature defaults, generate + build succeed
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

cat > dreego/components/Card.dreego << 'DREEGO'
Component Card (title string = "Default Title")
<div><h1>{title}</h1></div>
DREEGO

cat > dreego/routes/get.dreego << 'DREEGO'
<div><@Card/></div>
DREEGO

$DREEGO_BIN generate

# The generated component must apply the default, so the call site may omit the prop.
if ! grep -q 'Default Title' dreego/gen/components.go; then
    echo "FAIL: default value \"Default Title\" missing from generated component code"
    grep -n 'func Card' dreego/gen/components.go || true
    exit 1
fi

# The generated project must compile: <@Card/> without props must not produce
# "not enough arguments in call to Card".
go build -o /dev/null . || { echo "FAIL: go build failed (missing prop default)"; exit 1; }

echo ok
