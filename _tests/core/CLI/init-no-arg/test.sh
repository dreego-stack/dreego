#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that dreego init without an argument fails non-zero with usage
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

$DREEGO_BIN init 2>&1 && { echo "FAIL: expected non-zero exit when no path given"; exit 1; }

$DREEGO_BIN init 2>&1 | grep -q "usage: dreego init" || { echo "FAIL: expected usage message"; exit 1; }

echo ok
