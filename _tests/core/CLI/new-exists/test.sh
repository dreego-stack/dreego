#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that dreego new fails non-zero when the target directory already exists
set -e

realrepo="$(cd "$(dirname "$0")"/../../../.. && pwd)"
workdir="$(mktemp -d)"
trap "rm -rf $workdir" EXIT

cd "$workdir"

cat > go.mod << EOF
module t
go 1.22
require github.com/dreego-stack/dreego v0.0.0
replace github.com/dreego-stack/dreego => $realrepo
EOF

mkdir -p myapp

$DREEGO_BIN new myapp 2>&1 && { echo "FAIL: expected non-zero exit when target exists"; exit 1; }

$DREEGO_BIN new myapp 2>&1 | grep -q "already exists" || { echo "FAIL: expected 'already exists' message"; exit 1; }

echo ok
