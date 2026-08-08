#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Verify a head block works in a standalone page without a layout
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
<head>
    <title>Test</title>
    <script src="https://cdn.test/script.js"></script>
</head>
<div><p>hello</p></div>
DREEGO

$DREEGO_BIN generate

grep -q 'script.js' dreego/gen/routes.go
go build -o /dev/null .
echo ok
