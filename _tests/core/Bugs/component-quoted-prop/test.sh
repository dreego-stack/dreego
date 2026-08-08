#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: <@Card title="{myTitle}" active={true}/> resolves quoted {var} as expression
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

mkdir -p dreego/components dreego/routes

cat > dreego/components/Card.dreego << 'DREEGO'
Component Card (title string, active bool)
<div><h1>{title}</h1><span>{active}</span></div>
DREEGO

cat > dreego/routes/get.dreego << 'DREEGO'
<go>myTitle := "Hello"</go>
<div><@Card title="{myTitle}" active={true}/></div>
DREEGO

$DREEGO_BIN generate

# title="{myTitle}" must resolve to the expression myTitle, not a literal "{myTitle}".
if grep -q '{myTitle}' dreego/gen/routes.go; then
    echo "FAIL: title=\"{myTitle}\" emitted literally, expected resolved expression"
    grep -n 'Card(' dreego/gen/routes.go || true
    exit 1
fi
if ! grep -q 'Card(myTitle, true)' dreego/gen/routes.go; then
    echo "FAIL: expected Card(myTitle, true) call, got:"
    grep -n 'Card(' dreego/gen/routes.go || true
    exit 1
fi

# The generated project must compile (bool prop with unquoted {true}).
go build -o /dev/null . || { echo "FAIL: go build failed"; exit 1; }

echo ok
