#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: <script> body in a component is emitted literal, {x} not resolved as Go expr
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

mkdir -p dreego/components dreego/routes

# The <script> body must stay literal even though it contains {x}. The lexer
# treats <script> content as a raw text block where {…} is not a Go expression.
cat > dreego/components/Snippet.dreego << 'DREEGO'
Component Snippet (x string)
<div><script>const s = "literal {x}";</script></div>
DREEGO

cat > dreego/routes/get.dreego << 'DREEGO'
<div><@Snippet x=42/></div>
DREEGO

$DREEGO_BIN generate

# Desired: the generated component keeps the script body verbatim with {x} intact
# (goLiteral emits `` `literal {x}` `` since there is no backtick in the body).
if grep -q 'literal {x}' dreego/gen/components.go; then
    echo "ok: script body kept literal {x}"
else
    echo "FAIL: script body {x} was resolved as a Go expression"
    grep -n 'literal' dreego/gen/components.go || true
    exit 1
fi

# The whole project must still build (x stays literal, no undefined-var break).
go build -o /dev/null .
echo ok
