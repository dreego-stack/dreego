#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that <head> from both layout and route appears in generated code
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

mkdir -p dreego/layouts dreego/routes

cat > dreego/layouts/default.dreego << 'DREEGO'
<head>
    <title>Layout Title</title>
</head>
<div>{#slot}</div>
DREEGO

cat > dreego/routes/get.dreego << 'DREEGO'
<head><script src="route-script.js"></script></head>
<div><p>hello</p></div>
DREEGO

$DREEGO_BIN generate 2>&1

grep -q "Layout Title" dreego/gen/routes.go || { echo "FAIL: layout head dropped — title not in generated code"; exit 1; }
grep -q "route-script.js" dreego/gen/routes.go || { echo "FAIL: route head dropped — script not in generated code"; exit 1; }

go build -o /dev/null .
echo ok
