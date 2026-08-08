#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that an each-loop with the $loop variable compiles
set -e

realrepo="$(cd "$(dirname "$0")"/../../../.. && pwd)"
workdir="$(mktemp -d)"
trap "rm -rf $workdir" EXIT

cd "$workdir"

cat > go.mod << EOF
module t
go 1.22
require github.com/dreego-stack/dreego/core v0.0.0
replace github.com/dreego-stack/dreego/core => $realrepo/core
EOF

cat > main.go << 'GO'
package main
import _ "t/dreego/gen"
func main() {}
GO

mkdir -p dreego/routes

cat > dreego/routes/get.dreego << 'DREEGO'
<go>items := []string{"a", "b", "c"}</go>
<div>{#each items as item}<p>{$loop.Index}: {item}</p>{/each}</div>
DREEGO

$DREEGO_BIN generate
go build -o /dev/null .
echo ok
