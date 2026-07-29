#!/bin/sh
set -e
cd "$(dirname "$0")"
rm -rf dreego/* 2>/dev/null || true
mkdir -p dreego/routes
cat > dreego/routes/post.dreego << 'DREEGO'
<go type="json">
    var input map[string]any
    c.Bind(&input)
    input["echo"] = true
    c.JSON(200, input)
</go>
DREEGO
cat > main.go << 'GO'
package main
import (_ "t/dreego/gen"; core "codeberg.org/dreego/dreego/core")
func main() { core.Listen(":8080") }
GO
go mod init t >/dev/null 2>&1
go mod edit -replace codeberg.org/dreego/dreego=../../..
go mod edit -require codeberg.org/dreego/dreego@v0.0.0
sed -i 's|_ "gen"|_ "t/dreego/gen"|' main.go
dreego generate 2>&1
go build -o /tmp/srv .
/tmp/srv &
PID=$!
trap "kill $PID 2>/dev/null" EXIT
sleep 1
RESP=$(curl -s -H "Accept: application/json" -H "Content-Type: application/json" -d '{"name":"Lukas"}' http://localhost:8080/)
echo "$RESP" | grep -q '"name":"Lukas"' || { echo "FAIL: Bind not working, got: $RESP"; exit 1; }
echo "$RESP" | grep -q '"echo":true' || { echo "FAIL: no echo field, got: $RESP"; exit 1; }
echo ok
