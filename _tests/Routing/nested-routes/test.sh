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

mkdir -p dreego/routes/about dreego/routes/users/about

cat > dreego/routes/about/get.dreego << 'DREEGO'
<div><p>about page</p></div>
DREEGO

cat > dreego/routes/users/about/get.dreego << 'DREEGO'
<div><p>users about page</p></div>
DREEGO

go run codeberg.org/dreego/dreego/cmd/dreego generate
go build -o /dev/null .
echo ok
