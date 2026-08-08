#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that dreego generate --check reports stale after file modification
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

cat > main.go << 'GO'
package main
import _ "t/dreego/gen"
func main() {}
GO

mkdir -p dreego/routes

cat > dreego/routes/get.dreego << 'DREEGO'
<head><title>T</title></head>
<div><p>check me</p></div>
DREEGO

$DREEGO_BIN generate
$DREEGO_BIN generate --check 2>&1 | grep -q "up-to-date" || { echo "initial check failed"; exit 1; }
touch dreego/routes/get.dreego
if $DREEGO_BIN generate --check 2>/dev/null; then echo "expected stale but got up-to-date"; exit 1; fi
echo ok
