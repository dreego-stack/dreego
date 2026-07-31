#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that {#else} outside {#if} produces a compile error
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

cat > main.go << 'GO'
package main
import _ "t/dreego/gen"
func main() {}
GO

mkdir -p dreego/routes

cat > dreego/routes/get.dreego << 'DREEGO'
<div>{#else}</div>
DREEGO

if $DREEGO_BIN generate 2>/dev/null; then echo "expected failure but succeeded"; exit 1; fi
echo ok
