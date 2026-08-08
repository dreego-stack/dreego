#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that the compression middleware compiles
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

mkdir -p dreego/routes
cat > dreego/routes/get.dreego << 'DREEGO'
<div><p>ok</p></div>
DREEGO
cat > main.go << 'GO'
package main
import (
	_ "t/dreego/gen"
	dreego "github.com/dreego-stack/dreego/core"
)
func main() { dreego.Listen(":0") }
GO
$DREEGO_BIN generate
go build -o /dev/null .
echo ok
