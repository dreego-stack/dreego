#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that dreego generate --check reports up-to-date
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

$DREEGO_BIN init .
sed -i "s|_ \"gen\"|_ \"t/dreego/gen\"|" main.go
$DREEGO_BIN generate
$DREEGO_BIN generate --check 2>&1 | grep -q "up-to-date" || { echo "check failed"; exit 1; }
echo ok
