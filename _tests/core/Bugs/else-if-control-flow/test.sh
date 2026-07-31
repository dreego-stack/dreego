#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that {#else if} control flow compiles (not just {#else})
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
<go>score := 85</go>
<div>
{#if score >= 90}
A
{#else if score >= 80}
B
{#else}
C
{/if}
</div>
DREEGO

go run $realrepo/cmd/dreego generate
grep -q '{#else if' dreego/gen/routes.go && { echo "FAIL: {#else if} not transpiled"; exit 1; }
go build -o /dev/null .
echo ok
