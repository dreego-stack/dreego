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

mkdir -p dreego/routes dreego/components/button
printf 'Component FlatButton (label string)\n\n<div><button>{label}</button></div>' > dreego/components/button.dreego
printf 'Component NestedButton (label string)\n\n<div><button class="nested">{label}</button></div>' > dreego/components/button/button.dreego

cat > dreego/routes/get.dreego << 'DREEGO'
<div><@FlatButton label="Click"/><@NestedButton label="Go"/></div>
DREEGO

go run codeberg.org/dreego/dreego/cmd/dreego generate
go build -o /dev/null .
echo ok
