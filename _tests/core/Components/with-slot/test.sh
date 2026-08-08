#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that a component with a default slot compiles
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

mkdir -p dreego/components dreego/routes
cat > dreego/components/Card.dreego << 'DREEGO'
Component Card (title string)
<div><article><h2>{title}</h2><div>{#slot}</div></article></div>
DREEGO
cat > dreego/routes/get.dreego << 'DREEGO'
<div><@Card title="Welcome"><p>body text</p></@Card></div>
DREEGO
$DREEGO_BIN generate
go build -o /dev/null .
echo ok
