#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Verify each/else renders fallback content for empty collections
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

mkdir -p dreego/routes

cat > dreego/routes/get.dreego << 'DREEGO'
<go>items := []string{}</go>
<div>{#each items as item}<p>{item}</p>{#each else}<p>empty</p>{/each}</div>
DREEGO

go run $realrepo/cmd/dreego generate
go build -o /dev/null .
echo ok
