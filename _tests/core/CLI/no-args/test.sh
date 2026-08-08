#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that dreego with no arguments shows usage help
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

$DREEGO_BIN 2>&1 | grep -q "usage:" || { echo "no help shown"; exit 1; }
echo ok
