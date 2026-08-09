#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: {prop} in an HTML attribute is substituted, not emitted literally
set -e

realrepo="$(cd "$(dirname "$0")"/../../../.. && pwd)"
workdir="$(mktemp -d)"
trap "rm -rf $workdir" EXIT

cd "$workdir"

if [ -z "$DREEGO_BIN" ]; then
    DREEGO_BIN="$workdir/.dreego-bin"
    (cd "$realrepo" && go build -o "$DREEGO_BIN" ./cli/dreego) || { echo "FAIL: could not build dreego CLI"; exit 1; }
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
cat > dreego/components/Link.dreego << 'DREEGO'
Component Link (url string, label string)
<a href="{url}">{label}</a>
DREEGO
cat > dreego/routes/get.dreego << 'DREEGO'
<div><@Link url="https://example.com" label="Go"/></div>
DREEGO

$DREEGO_BIN generate

generated="dreego/gen/components.go"
if grep -q '{url}' "$generated"; then
    echo "FAIL: {url} emitted literally in attribute, expected resolved expression"
    exit 1
fi
if grep -q '{label}' "$generated"; then
    echo "FAIL: {label} emitted literally in text content"
    exit 1
fi
if ! grep -q 'EscapeString' "$generated"; then
    echo "FAIL: attribute expression must be html-escaped (XSS-safe)"
    exit 1
fi
go build -o /dev/null .
echo ok
