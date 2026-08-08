#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: scaffold either contains a working layout or no layouts/ directory at all
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
require github.com/dreego-stack/dreego v0.0.0
replace github.com/dreego-stack/dreego => $realrepo
EOF

$DREEGO_BIN new testapp 2>&1

[ -d testapp ] || { echo "FAIL: testapp directory not created"; exit 1; }

if [ ! -d testapp/dreego/layouts ]; then
    echo "ok"
    exit 0
fi

if [ -z "$(find testapp/dreego/layouts -type f -name '*.dreego' 2>/dev/null)" ]; then
    echo "FAIL: layouts/ directory exists but contains no .dreego layout file"
    exit 1
fi

cd testapp
$DREEGO_BIN generate 2>&1

if ! grep -q '<html>' dreego/gen/routes.go; then
    echo "FAIL: layout exists but generated route does not produce a complete HTML document (no <html> found)"
    exit 1
fi

echo "ok"
