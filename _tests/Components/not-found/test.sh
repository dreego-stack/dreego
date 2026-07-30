#!/bin/sh
# Using standard: _tests/how-to-test-sh.md

realrepo="$(cd "$(dirname "$0")"/../../.. && pwd)"
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

mkdir -p dreego/routes dreego/components
cat > dreego/routes/get.dreego << 'DREEGO'
<div><@Missing/></div>
DREEGO

go run $realrepo/cmd/dreego generate
if go build -o /dev/null . 2>/dev/null; then echo "expected build failure but succeeded"; exit 1; fi
echo ok
