#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that a <div> inside a component slot compiles correctly
set -e

realrepo="$(cd "$(dirname "$0")"/../../.. && pwd)"
workdir="$(mktemp -d)"
trap "rm -rf $workdir" EXIT

cd "$workdir"

cat > go.mod << EOF
module t
go 1.22
require codeberg.org/dreego/dreego v0.0.0
replace codeberg.org/dreego/dreego => $realrepo
EOF

cat > main.go << 'GO'
package main
import _ "t/dreego/gen"
func main() {}
GO

mkdir -p dreego/components dreego/routes

cat > dreego/components/Card.dreego << 'DREEGO'
Component Card ()
<div><article>{#slot}</article></div>
DREEGO

cat > dreego/routes/get.dreego << 'DREEGO'
<div><@Card><div class=\"inner\">hi</div></@Card></div>
DREEGO

go run $realrepo/cmd/dreego generate
go build -o /dev/null .
echo ok
