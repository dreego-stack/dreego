#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that session set/get operations compile
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

import (
	_ "t/dreego/gen"

	dreego "github.com/dreego-stack/dreego/core"
)

func main() {
	dreego.Listen(":8080")
}
GO

mkdir -p dreego/routes

cat > dreego/routes/get.dreego << 'DREEGO'
<go>
    c.SetSessionVal("key", "val")
    v := c.SessionVal("key")
    _ = v
</go>
<div><p>session set/get</p></div>
DREEGO

$DREEGO_BIN generate
go build -o /dev/null .
echo ok
