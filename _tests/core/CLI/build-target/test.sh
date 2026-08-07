#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that dreego build --target linux/amd64 cross-compiles a binary to build/bin/<name>-linux-amd64
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

mkdir -p dreego/routes

cat > dreego/routes/get.dreego << 'DREEGO'
<div><p>hello</p></div>
DREEGO

cat > main.go << GO
package main

import (
	_ "t/dreego/gen"
	dreego "codeberg.org/dreego/dreego/core"
)

func main() {
	dreego.Listen(":8080")
}
GO

$DREEGO_BIN build --target linux/amd64

name="$(basename "$workdir")"
bin="build/bin/${name}-linux-amd64"
[ -f "$bin" ] || { echo "FAIL: binary not found at $bin"; exit 1; }
[ -x "$bin" ] || { echo "FAIL: binary not executable at $bin"; exit 1; }

echo ok
