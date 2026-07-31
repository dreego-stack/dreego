#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests deep directory nesting for routes
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

mkdir -p dreego/routes/a/b/c/d

cat > dreego/routes/a/b/c/d/get.dreego << 'DREEGO'
<div><p>deep</p></div>
DREEGO

$DREEGO_BIN generate
go build -o /dev/null .
echo ok
