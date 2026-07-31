#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that dreego with an unknown command exits non-zero
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

$DREEGO_BIN bogus 2>&1 && { echo "FAIL: expected non-zero exit"; exit 1; }

echo ok
