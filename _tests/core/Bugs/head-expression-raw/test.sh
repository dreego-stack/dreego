#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: A template expression inside a route's <head> section must be resolved, not emitted raw
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

mkdir -p dreego/routes

cat > dreego/routes/get.dreego << 'DREEGO'
<go>
type Doc struct{ Title string }
doc := Doc{Title: "PeerNet Docs"}
</go>
<head><title>{doc.Title} — PeerNet Docs</title></head>
<div><h1>{doc.Title}</h1></div>
DREEGO

$DREEGO_BIN generate

generated="dreego/gen/routes.go"
if grep -Fq '{doc.Title}' "$generated"; then
    echo "FAIL: head expression {doc.Title} emitted raw instead of resolved"
    exit 1
fi
grep -Fq 'doc.Title' "$generated" || {
    echo "FAIL: resolved doc.Title missing from generated code"
    exit 1
}
go build -o /dev/null .
echo ok