#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Verify a page containing head, Go, div, script, and style sections
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

mkdir -p dreego/routes

cat > dreego/routes/get.dreego << 'DREEGO'
<head><title>All</title></head>
<go>x := "ok"</go>
<div><p>{x}</p></div>
<script>console.log("js")</script>
<style>p { color: red; }</style>
DREEGO

$DREEGO_BIN generate
go build -o /dev/null .
echo ok
