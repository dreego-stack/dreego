#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that named slots with expressions compile
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

mkdir -p dreego/components dreego/routes

cat > dreego/components/Greet.dreego << 'DREEGO'
Component Greet ()
<div><p>{#slot header}{/slot} {#slot}</p></div>
DREEGO

cat > dreego/routes/get.dreego << 'DREEGO'
<go>name := "World"</go>
<div><@Greet>{#slot header}<strong>{name}</strong>{/slot}!</@Greet></div>
DREEGO

$DREEGO_BIN generate
go build -o /dev/null .
echo ok
