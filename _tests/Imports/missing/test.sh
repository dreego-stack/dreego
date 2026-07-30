#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
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

mkdir -p dreego/routes
cat > dreego/routes/get.dreego << 'DREEGO'
import "dreego/components/Nope"
<div><p>hi</p></div>
DREEGO
go run codeberg.org/dreego/dreego/cmd/dreego generate
go build -o /dev/null .
echo ok
