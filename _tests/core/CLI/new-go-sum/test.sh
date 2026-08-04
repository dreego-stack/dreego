#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: fresh dreego new scaffold builds without go mod tidy
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

$DREEGO_BIN new testapp 2>&1

[ -d testapp ] || { echo "FAIL: testapp directory not created"; exit 1; }
[ -f testapp/go.sum ] || { echo "FAIL: missing go.sum"; exit 1; }

cd testapp
$DREEGO_BIN generate
GOWORK=off go build .

echo "ok"
