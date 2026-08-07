#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: $loop.X in a {#if} cond inside {#each} is substituted to loop.X, generate + build succeed
set -e

realrepo="$(cd "$(dirname "$0")"/../../../.. && pwd)"
workdir="$(mktemp -d)"
trap "rm -rf $workdir" EXIT

cd "$workdir"

if [ -z "$DREEGO_BIN" ]; then
    DREEGO_BIN="$workdir/.dreego-bin"
    (cd "$realrepo" && go build -o "$DREEGO_BIN" ./cmd/dreego) || { echo "FAIL: could not build dreego CLI"; exit 1; }
    export DREEGO_BIN
fi

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

mkdir -p dreego/routes

cat > dreego/routes/get.dreego << 'DREEGO'
<go>items := []string{"a", "b", "c"}</go>
<div>
{#each items as item}<div>{#if !$loop.Last}, {/if}{item}</div>{/each}
</div>
DREEGO

$DREEGO_BIN generate

# The generated Go must reference loop.Last (substituted), never raw $loop.
# Routes land in gen/routes.go (gen/dree.go only carries settings/static).
if ! grep -q 'loop.Last' dreego/gen/routes.go; then
    echo "FAIL: generated code missing loop.Last (substitution did not happen)"
    grep -n 'loop' dreego/gen/routes.go || true
    exit 1
fi
if grep -q '\$loop' dreego/gen/routes.go; then
    echo "FAIL: raw \$loop. left in generated code"
    grep -n '\$loop' dreego/gen/routes.go || true
    exit 1
fi

# The generated project must compile: "if !loop.Last {" is valid Go.
go build -o /dev/null . || { echo "FAIL: go build failed (\$loop in cond not substituted)"; exit 1; }

echo ok
