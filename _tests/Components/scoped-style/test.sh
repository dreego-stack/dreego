#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
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
cat > dreego/components/Box.dreego << 'DREEGO'
Component Box ()
<div><div class="box"><p>scoped</p></div></div>
<style>.box{border:1px solid red}</style>
DREEGO
cat > dreego/routes/get.dreego << 'DREEGO'
<head><title>T</title></head>
<div><@Box/><p class="box">unscoped</p></div>
<style>.box{color:blue}</style>
DREEGO
go run codeberg.org/dreego/dreego/cmd/dreego generate
go build -o /dev/null .
echo ok
