#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that a multi-segment route parameter compiles
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

mkdir -p dreego/routes/a

cat > dreego/routes/a/get.dreego << 'DREEGO'
<go>a:=c.Param("a")</go>
<div><p>{a}</p></div>
DREEGO

$DREEGO_BIN generate
go build -o /dev/null .
echo ok
