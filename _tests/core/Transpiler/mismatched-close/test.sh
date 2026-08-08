#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Verify mismatched close tags produce a compile error
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
<go>x:=1</go>
<div>text</go>
DREEGO

if $DREEGO_BIN generate 2>/dev/null; then echo "expected failure but succeeded"; exit 1; fi
echo ok
