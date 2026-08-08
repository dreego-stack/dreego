#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that a component with no props compiles
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

mkdir -p dreego/components dreego/routes

cat > dreego/components/Empty.dreego << 'DREEGO'
Component Empty ()
<div><p>no props</p></div>
DREEGO

cat > dreego/routes/get.dreego << 'DREEGO'
<div><@Empty/></div>
DREEGO

$DREEGO_BIN generate
go build -o /dev/null .
echo ok
