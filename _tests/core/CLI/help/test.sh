#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that dreego --help and dreego -h show help output
set -e

realrepo="$(cd "$(dirname "$0")"/../../../.. && pwd)"
workdir="$(mktemp -d)"
trap "rm -rf $workdir" EXIT

cd "$workdir"

cat > go.mod << EOF
module t
go 1.22
require codeberg.org/dreego/dreego/core v0.0.0
replace codeberg.org/dreego/dreego/core => $realrepo/core
EOF

$DREEGO_BIN --help 2>&1 | grep -q "usage:" || { echo "FAIL: --help did not show usage"; exit 1; }
$DREEGO_BIN -h 2>&1 | grep -q "usage:" || { echo "FAIL: -h did not show usage"; exit 1; }

echo ok
