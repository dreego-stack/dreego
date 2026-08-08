#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that static file serving compiles
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

mkdir -p dreego/routes dreego/static

cat > dreego/routes/get.dreego << 'DREEGO'
<div><p>hello</p></div>
DREEGO

printf 'body{color:red}' > dreego/static/style.css

$DREEGO_BIN generate
grep -q 'RegisterStatic("/style.css"' dreego/gen/dree.go
go build -o /dev/null .
echo ok
