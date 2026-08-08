#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Verify the raw filter outputs unescaped HTML while default escaping is applied
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

mkdir -p dreego/routes

cat > dreego/routes/get.dreego << 'DREEGO'
<go>rawHtml := "<b>bold</b>"</go>
<div><p>{rawHtml|raw}</p><p>{rawHtml}</p></div>
DREEGO

$DREEGO_BIN generate
go build -o /dev/null .
grep -q 'b.WriteString(fmt.Sprintf("%v", rawHtml))' dreego/gen/routes.go
grep -q 'b.WriteString(html.EscapeString(fmt.Sprintf("%v", rawHtml)))' dreego/gen/routes.go
echo ok
