#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: < inside Go strings (quoted + backtick) survives generate + build
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
require github.com/dreego-stack/dreego/core v0.0.0
replace github.com/dreego-stack/dreego/core => $realrepo/core
EOF

cat > main.go << 'GO'
package main
import _ "t/dreego/gen"
func main() {}
GO

mkdir -p dreego/routes

cat > dreego/routes/get.dreego << 'DREEGO'
<go>
msg := "TO: <HASH>"
svg := `<svg viewBox="0 0 24 24"><path d="M12 2"/></svg>`
</go>
<div>
<p>{msg}</p>
<p>{svg}</p>
</div>
DREEGO

$DREEGO_BIN generate

# The Go section must survive verbatim in the generated handler: the quoted
# string keeps <HASH>, the backtick string keeps <svg>...</svg>.
# splitGoSections puts non-declaration code inline in the handler (routes.go).
if ! grep -q 'TO: <HASH>' dreego/gen/routes.go; then
    echo "FAIL: <HASH> lost from generated Go (quoted string destroyed)"
    grep -n 'HASH' dreego/gen/routes.go || true
    exit 1
fi
if ! grep -q '<svg viewBox="0 0 24 24">' dreego/gen/routes.go; then
    echo "FAIL: <svg> lost from generated Go (backtick string destroyed)"
    grep -n 'svg' dreego/gen/routes.go || true
    exit 1
fi
# The self-closing <path .../> must survive verbatim — the trailing / was
# silently dropped before the fix (silent data corruption).
if ! grep -Fq '<path d="M12 2"/>' dreego/gen/routes.go; then
    echo "FAIL: self-closing <path d=\"M12 2\"/> lost its trailing / in generated Go"
    grep -n 'path' dreego/gen/routes.go || true
    exit 1
fi
if grep -q 'TO:' dreego/gen/dree.go; then
    echo "FAIL: go body leaked as template text into dree.go"
    grep -n 'TO:' dreego/gen/dree.go || true
    exit 1
fi

# The generated project must compile: msg and svg are valid Go strings.
go build -o /dev/null . || { echo "FAIL: go build failed (go string with < broken)"; exit 1; }

echo ok
