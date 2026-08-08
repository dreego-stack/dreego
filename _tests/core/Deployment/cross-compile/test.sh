#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests cross-compilation with dreego build --target flag
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
<div><h1>hello</h1></div>
DREEGO
cat > main.go << 'GO'
package main
import (_ "t/dreego/gen"; dreego "github.com/dreego-stack/dreego/core")
func main() { dreego.Listen(":0") }
GO
$DREEGO_BIN generate 2>&1
# cross-compile for linux/amd64 (native in Docker, but exercises the flag)
$DREEGO_BIN build --target linux/amd64 2>&1
BIN=$(ls build/bin/*-linux-amd64 2>/dev/null)
[ -n "$BIN" ] || { echo "FAIL: binary not found"; exit 1; }
[ -x "$BIN" ] || { echo "FAIL: not executable"; exit 1; }
echo ok
