#!/bin/sh
set -e
cd "$(dirname "$0")"
rm -rf dreego/* 2>/dev/null || true
mkdir -p dreego/routes
cat > dreego/routes/get.dreego << 'DREEGO'
<div><h1>hello</h1></div>
DREEGO
cat > main.go << 'GO'
package main
import (_ "t/dreego/gen"; core "codeberg.org/dreego/dreego/core")
func main() { core.Listen(":0") }
GO
go mod init t >/dev/null 2>&1
go mod edit -replace codeberg.org/dreego/dreego=../../..
go mod edit -require codeberg.org/dreego/dreego@v0.0.0
sed -i 's|_ "gen"|_ "t/dreego/gen"|' main.go
dreego generate 2>&1
# cross-compile for linux/amd64 (native in Docker, but exercises the flag)
dreego build --target linux/amd64 2>&1
BIN=$(ls build/bin/*-linux-amd64 2>/dev/null)
[ -n "$BIN" ] || { echo "FAIL: binary not found"; exit 1; }
# verify it's executable
[ -x "$BIN" ] || { echo "FAIL: not executable"; exit 1; }
echo ok
