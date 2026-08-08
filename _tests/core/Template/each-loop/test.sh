#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that an each-loop template compiles
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

cat > main.go << 'GO'
package main
import _ "t/dreego/gen"
func main() {}
GO

mkdir -p dreego/routes

cat > dreego/routes/get.dreego << 'DREEGO'
<go>items := []string{"a", "b"}</go>
<ul>{#each items as item}<li>{item}</li>{/each}</ul>
DREEGO

$DREEGO_BIN generate
grep -q '{#each' dreego/gen/routes.go && { echo "FAIL: {#each} not transpiled in route"; exit 1; }
go build -o /dev/null .
echo ok
