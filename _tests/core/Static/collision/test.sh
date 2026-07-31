#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that route and static path collision produces an error
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

mkdir -p dreego/routes/about dreego/static

cat > dreego/routes/about/get.dreego << 'DREEGO'
<div><p>about</p></div>
DREEGO

printf 'text' > dreego/static/about

if go run $realrepo/cmd/dreego generate 2>/dev/null; then echo "expected collision error but succeeded"; exit 1; fi
echo ok
