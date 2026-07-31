#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that named slots in components compile
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

printf 'Component Card (title string)\n<div><article>{#slot header}{/slot}<h2>{title}</h2><div>{#slot}</div></article></div>\n' > dreego/components/Card.dreego
printf '<div><@Card title="Hi">{#slot header}<strong>HEADER</strong>{/slot}<p>body</p></@Card></div>\n' > dreego/routes/get.dreego

go run $realrepo/cmd/dreego generate
go build -o /dev/null .
echo ok
