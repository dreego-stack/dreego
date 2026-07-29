#!/bin/sh
set -e
cd "$(dirname "$0")"
rm -rf dreego/* 2>/dev/null || true
mkdir -p dreego/routes
cat > dreego/routes/get.dreego << 'DREEGO'
<go>
    name := "Lukas"
</go>
<go type="json">
    c.JSON(200, map[string]string{"format": "json", "name": name})
</go>
<go type="xml">
    user := struct{XMLName struct{} `xml:"user"`; Name string `xml:"name"`}{Name: name}
    c.XML(200, user)
</go>
<div><h1>{name}</h1></div>
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
J=$(curl -s -H "Accept: application/json" http://localhost:8080/)
echo "$J" | grep -q '"format":"json"' || { echo "FAIL: JSON broken, got: $J"; exit 1; }
X=$(curl -s -H "Accept: application/xml" http://localhost:8080/)
echo "$X" | grep -q '<user>' || { echo "FAIL: XML broken, got: $X"; exit 1; }
H=$(curl -s http://localhost:8080/)
echo "$H" | grep -q '<h1>Lukas</h1>' || { echo "FAIL: HTML broken, got: $H"; exit 1; }
echo ok
