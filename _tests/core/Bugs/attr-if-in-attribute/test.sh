#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: {#if} inside an HTML attribute must fail fast with a clear error,
#       never produce broken Go. Wrapping the whole tag in {#if} works.
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
require github.com/dreego-stack/dreego v0.0.0
replace github.com/dreego-stack/dreego => $realrepo
EOF

cat > main.go << 'GO'
package main
import _ "t/dreego/gen"
func main() {}
GO

mkdir -p dreego/routes dreego/components

# 1. Route with {#if} inside an attribute: generate must fail with a clear error.
cat > dreego/routes/get.dreego << 'DREEGO'
<go>cond := true</go>
<div><a class="nav {#if cond}active{/if}">link</a></div>
DREEGO

if out=$($DREEGO_BIN generate 2>&1); then
    echo "FAIL: route with {#if} in attribute must fail generate"
    exit 1
fi
if ! echo "$out" | grep -q 'inside attribute value'; then
    echo "FAIL: route error must mention 'inside attribute value', got: $out"
    exit 1
fi

# 2. Route with {#if} in the <div> section opener attribute: generate must fail.
cat > dreego/routes/get.dreego << 'DREEGO'
<go>cond := true</go>
<div class="nav {#if cond}active{/if}"><p>link</p></div>
DREEGO

if out=$($DREEGO_BIN generate 2>&1); then
    echo "FAIL: route with {#if} in <div> opener attribute must fail generate"
    exit 1
fi
if ! echo "$out" | grep -q 'inside attribute value'; then
    echo "FAIL: div-opener error must mention 'inside attribute value', got: $out"
    exit 1
fi

# 3. Route with {#if} in a single-quoted attribute: generate must fail.
cat > dreego/routes/get.dreego << 'DREEGO'
<go>cond := true</go>
<div><a class='nav {#if cond}active{/if}'>link</a></div>
DREEGO

if out=$($DREEGO_BIN generate 2>&1); then
    echo "FAIL: route with {#if} in single-quoted attribute must fail generate"
    exit 1
fi
if ! echo "$out" | grep -q 'inside attribute value'; then
    echo "FAIL: single-quote error must mention 'inside attribute value', got: $out"
    exit 1
fi

# 4. Component with {#if} inside an attribute: generate must fail with a clear error.
cat > dreego/components/Nav.dreego << 'DREEGO'
Component Nav (cond bool)
<a class="nav {#if cond}active{/if}">x</a>
DREEGO

if out=$($DREEGO_BIN generate 2>&1); then
    echo "FAIL: component with {#if} in attribute must fail generate"
    exit 1
fi
if ! echo "$out" | grep -q 'inside attribute value'; then
    echo "FAIL: component error must mention 'inside attribute value', got: $out"
    exit 1
fi

# 5. Workaround: wrap the whole tag in {#if} — must generate and compile.
rm dreego/routes/get.dreego dreego/components/Nav.dreego

cat > dreego/routes/get.dreego << 'DREEGO'
<go>cond := true</go>
<div>{#if cond}<a class="nav active">link</a>{/if}</div>
DREEGO

cat > dreego/components/Nav.dreego << 'DREEGO'
Component Nav (cond bool)
{#if cond}<a class="nav active">x</a>{/if}
DREEGO

$DREEGO_BIN generate

if grep -q '#if' dreego/gen/dree.go dreego/gen/components.go; then
    echo "FAIL: {#if} leaked into generated Go"
    grep -n '#if' dreego/gen/dree.go dreego/gen/components.go || true
    exit 1
fi

go build -o /dev/null . || { echo "FAIL: go build failed"; exit 1; }

# 6. Workaround for the <div> opener: wrap the whole section tag in {#if}.
rm dreego/routes/get.dreego dreego/components/Nav.dreego

cat > dreego/routes/get.dreego << 'DREEGO'
<go>cond := true</go>
{#if cond}<div class="nav active"><p>link</p></div>{/if}
DREEGO

$DREEGO_BIN generate

if grep -q '#if' dreego/gen/dree.go; then
    echo "FAIL: {#if} leaked into generated Go (div workaround)"
    grep -n '#if' dreego/gen/dree.go || true
    exit 1
fi

go build -o /dev/null . || { echo "FAIL: go build failed (div workaround)"; exit 1; }

echo ok
