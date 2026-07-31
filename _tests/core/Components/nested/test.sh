#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that a component nested inside another component compiles
set -e

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

cat > main.go << 'GO'
package main
import _ "t/dreego/gen"
func main() {}
GO

mkdir -p dreego/components dreego/routes
cat > dreego/components/Inner.dreego << 'DREEGO'
Component Inner ()
<div><span>inner</span></div>
DREEGO
cat > dreego/components/Outer.dreego << 'DREEGO'
Component Outer ()
<div><article><@Inner/></article></div>
DREEGO
cat > dreego/routes/get.dreego << 'DREEGO'
<div><@Outer/></div>
DREEGO
go run $realrepo/cmd/dreego generate
go build -o /dev/null .
echo ok
