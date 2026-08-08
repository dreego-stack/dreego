#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that dreego generate --check fails when no generated files exist
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

mkdir -p dreego/routes

cat > dreego/routes/get.dreego << 'DREEGO'
<div><p>hello</p></div>
DREEGO

$DREEGO_BIN generate --check 2>&1 && { echo "FAIL: expected non-zero exit when no generated files exist"; exit 1; }

$DREEGO_BIN generate --check 2>&1 | grep -q "no generated files found" || { echo "FAIL: expected 'no generated files found' message"; exit 1; }

echo ok
