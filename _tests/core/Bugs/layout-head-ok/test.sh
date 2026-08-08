#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that layout and route head content merges correctly in generated code
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
import (_ "t/dreego/gen"; dreego "github.com/dreego-stack/dreego/core")
func main() { dreego.Listen(":0") }
GO

mkdir -p dreego/layouts dreego/routes

cat > dreego/layouts/default.dreego << 'DREEGO'
<head>
    <title>Layout Title</title>
    <script src="https://cdn.tailwindcss.com"></script>
    {#head}
</head>
<div><main>{#slot}</main></div>
DREEGO

cat > dreego/routes/get.dreego << 'DREEGO'
<head><meta name="desc" content="route meta"></head>
<div><h1>Page</h1></div>
DREEGO

$DREEGO_BIN generate 2>&1

grep -q "cdn.tailwindcss.com" dreego/gen/routes.go || { echo "FAIL: tailwind CDN not in generated code"; exit 1; }
grep -q "Layout Title" dreego/gen/routes.go || { echo "FAIL: title not in generated code"; exit 1; }
grep -q "route meta" dreego/gen/routes.go || { echo "FAIL: route meta not in generated code"; exit 1; }

go build -o /dev/null .
echo ok
