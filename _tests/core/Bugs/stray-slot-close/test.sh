#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that a stray {/slot} close tag produces a compile error

realrepo="$(cd "$(dirname "$0")"/../../../.. && pwd)"
workdir="$(mktemp -d)"
trap "rm -rf $workdir" EXIT

cd "$workdir"

cat > go.mod << EOF
module t
go 1.22
require codeberg.org/dreego/dreego v0.0.0
replace codeberg.org/dreego/dreego => $realrepo
EOF

mkdir -p dreego/components dreego/routes

cat > dreego/components/Card.dreego << 'DREEGO'
Component Card ()
<div><article>{#slot header}</article></div>
DREEGO

cat > dreego/routes/get.dreego << 'DREEGO'
<div><@Card>{/slot}</@Card></div>
DREEGO

if go run $realrepo/cmd/dreego generate 2>/dev/null; then echo "expected failure but succeeded"; exit 1; fi
echo ok
