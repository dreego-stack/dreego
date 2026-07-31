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
require codeberg.org/dreego/dreego v0.0.0
replace codeberg.org/dreego/dreego => $realrepo
EOF

cat > main.go << 'GO'
package main

import (
	_ "t/dreego/gen"

	core "codeberg.org/dreego/dreego/core"
)

func main() {
	core.Listen(":8080")
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

go run $realrepo/cmd/dreego generate
go build -o /dev/null .
echo ok
