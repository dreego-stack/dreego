#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that {#if} and {#each} inside a Component compile (B1)
set -e

realrepo="$(cd "$(dirname "$0")"/../../../.. && pwd)"
workdir="$(mktemp -d)"
trap "rm -rf $workdir" EXIT

cd "$workdir"

cat > go.mod << EOF
module t
go 1.22
require codeberg.org/dreego/dreego/core v0.0.0
replace codeberg.org/dreego/dreego/core => $realrepo/core
EOF

cat > main.go << 'GO'
package main
import _ "t/dreego/gen"
func main() {}
GO

mkdir -p dreego/components dreego/routes

cat > dreego/components/List.dreego << 'DREEGO'
Component List (items []string)
<ul>
{#each items as item}
    <li>{item}</li>
{/each}
</ul>
DREEGO

cat > dreego/routes/get.dreego << 'DREEGO'
<go>items := []string{"a", "b"}</go>
<div><@List items={items}/></div>
DREEGO

$DREEGO_BIN generate
grep -q '{#each' dreego/gen/components.go && { echo "FAIL: {#each} not transpiled in component (B1)"; exit 1; }
grep -q '{#if' dreego/gen/components.go && { echo "FAIL: {#if} not transpiled in component (B1)"; exit 1; }
go build -o /dev/null .
echo ok
