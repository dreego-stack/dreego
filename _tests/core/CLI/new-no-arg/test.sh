#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that dreego new without a name fails non-zero with usage
set -e

realrepo="$(cd "$(dirname "$0")"/../../../.. && pwd)"
workdir="$(mktemp -d)"
trap "rm -rf $workdir" EXIT

cd "$workdir"

cat > go.mod << EOF
module t
go 1.22
require github.com/dreego-stack/dreego/core v0.0.0
replace github.com/dreego-stack/dreego/core => $realrepo/core
EOF

$DREEGO_BIN new 2>&1 && { echo "FAIL: expected non-zero exit when no name given"; exit 1; }

$DREEGO_BIN new 2>&1 | grep -q "usage: dreego new" || { echo "FAIL: expected usage message"; exit 1; }

echo ok
