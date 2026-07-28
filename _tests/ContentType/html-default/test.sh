#!/bin/sh
set -e
cd "$(dirname "$0")"
rm -rf dreego/* 2>/dev/null || true
mkdir -p dreego/routes
cat > dreego/routes/get.dreego << 'DREEGO'
<go>
    msg := "hello"
</go>
<div><h1>{msg}</h1></div>
DREEGO
mkdir -p dreego/gen
go mod init t >/dev/null 2>&1
go mod edit -replace codeberg.org/dreego/dreego=../../..
go mod edit -require codeberg.org/dreego/dreego@v0.0.0
cat > main.go << 'GO'
package main
import (_ "t/dreego/gen"; core "codeberg.org/dreego/dreego/dreego-core")
func main() { core.Listen(":0") }
GO
dreego generate 2>&1
grep -q "text/html" dreego/gen/routes.go || { echo "FAIL: no text/html content-type"; exit 1; }
grep -q "b.WriteString" dreego/gen/routes.go || { echo "FAIL: template rendering missing"; exit 1; }
go build -o /dev/null .
echo ok
