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
require github.com/dreego-stack/dreego/core v0.0.0
replace github.com/dreego-stack/dreego/core => $realrepo/core
EOF

mkdir -p dreego/components dreego/routes

cat > dreego/components/Card.dreego << 'DREEGO'
Component Card ()
<div><article>{#slot header}</article></div>
DREEGO

cat > dreego/routes/get.dreego << 'DREEGO'
<div><@Card>{/slot}</@Card></div>
DREEGO

if $DREEGO_BIN generate 2>/dev/null; then echo "expected failure but succeeded"; exit 1; fi
echo ok
