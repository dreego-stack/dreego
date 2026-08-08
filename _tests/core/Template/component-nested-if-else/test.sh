#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: A component with nested {#if} inside a non-final {#else} renders correctly.
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

mkdir -p dreego/components dreego/routes

cat > dreego/components/Grade.dreego << 'DREEGO'
Component Grade (score int)
<div class="grade">
{#if score >= 90}
A
{#else}
{#if score >= 80}
B
{#else}
C
{/if}
D
{/if}
</div>
DREEGO

cat > dreego/routes/get.dreego << 'DREEGO'
<go>score := 85</go>
<div>
<@Grade score={score}/>
</div>
DREEGO

$DREEGO_BIN generate
grep -q 'if score >= 80' dreego/gen/components.go || {
    echo "FAIL: nested {#if} in component {#else} branch was silently dropped"
    exit 1
}
go build -o /dev/null .
echo ok
