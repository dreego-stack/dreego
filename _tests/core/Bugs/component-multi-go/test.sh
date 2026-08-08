#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that GenerateComponent emits all <go> sections, not just Go[0] (B14)
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

mkdir -p dreego/components dreego/routes

cat > dreego/components/Greet.dreego << 'DREEGO'
Component Greet (name string)
<go>greeting := "hello"</go>
<go>msg := greeting + " world"</go>
<div>{msg} {name}</div>
DREEGO

cat > dreego/routes/get.dreego << 'DREEGO'
<div><@Greet name="Ada"/></div>
DREEGO

$DREEGO_BIN generate
generated="dreego/gen/components.go"
if ! grep -q 'greeting := "hello"' "$generated"; then
    echo "FAIL: first <go> section missing in generated component (B14)"
    exit 1
fi
if ! grep -q 'msg := greeting + " world"' "$generated"; then
    echo "FAIL: second <go> section missing in generated component (B14)"
    exit 1
fi
go build -o /dev/null .
echo ok
