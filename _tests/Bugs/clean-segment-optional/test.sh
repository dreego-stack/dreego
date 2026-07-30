#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that cleanSegment strips all bracket pairs from optional segments (B15)
set -e

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

mkdir -p dreego/routes/[[opt]]

cat > dreego/routes/[[opt]]/get.dreego << 'DREEGO'
<div>optional</div>
DREEGO

go run $realrepo/cmd/dreego generate
generated="dreego/gen/routes.go"
if grep -q "\[opt\]" "$generated"; then
    echo "FAIL: optional segment still contains brackets (B15)"
    exit 1
fi
if grep -q 'core.Register("GET", "/{opt}",' "$generated"; then
    :
else
    echo "FAIL: optional segment not registered as /{opt}"
    exit 1
fi
go build -o /dev/null .
echo ok
