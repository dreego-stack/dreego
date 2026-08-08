#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that multiple named slots in a component compile
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

cat > main.go << 'GO'
package main
import _ "t/dreego/gen"
func main() {}
GO

mkdir -p dreego/components dreego/routes

cat > dreego/components/Page.dreego << 'DREEGO'
Component Page (title string)
<div><header>{#slot header}{/slot}</header><main>{#slot}</main><footer>{#slot footer}{/slot}</footer></div>
DREEGO

cat > dreego/routes/get.dreego << 'DREEGO'
<div><@Page title="Multi">
{#slot header}<nav>menu</nav>{/slot}
content
{#slot footer}<small>2026</small>{/slot}
</@Page></div>
DREEGO

$DREEGO_BIN generate
go build -o /dev/null .
echo ok
