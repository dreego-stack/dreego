#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that a basic component with props compiles
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

mkdir -p dreego/components dreego/routes

cat > dreego/components/Cmp.dreego << 'DREEGO'
Component Card (title string)
<div><article><h2>{title}</h2></article></div>
DREEGO

cat > dreego/routes/get.dreego << 'DREEGO'
<div><@Card title="Hello"/></div>
DREEGO

$DREEGO_BIN generate
go build -o /dev/null .
echo ok
