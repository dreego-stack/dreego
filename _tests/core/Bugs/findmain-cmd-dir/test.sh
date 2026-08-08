#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that findMain detects cmd/main.go (B6)
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

mkdir -p cmd dreego/routes

cat > cmd/main.go << 'GO'
package main
import _ "t/dreego/gen"
func main() {}
GO

cat > dreego/routes/get.dreego << 'DREEGO'
<div>hello</div>
DREEGO

$DREEGO_BIN build
if [ ! -f "cmd/build/bin/cmd" ]; then
    echo "FAIL: findMain did not build cmd/main.go (B6)"
    exit 1
fi
echo ok
