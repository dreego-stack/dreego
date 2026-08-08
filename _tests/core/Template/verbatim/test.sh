#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Verify the verbatim block preserves raw JavaScript content in generated code
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
<div><p>before</p>{#verbatim}<script>var x = {a: 1};</script>{/verbatim}<p>after</p></div>
DREEGO

$DREEGO_BIN generate
go build -o /dev/null .
grep -q 'var x = {a: 1}' dreego/gen/routes.go
echo ok
