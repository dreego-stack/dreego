#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Component prop value from expression title={user.Name} compiles
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
Component Card (title string)
<div><article><h2>{title}</h2></article></div>
DREEGO
cat > dreego/routes/get.dreego << 'DREEGO'
<go>type User struct { Name string }
user := User{Name: "Ada"}</go>
<div><@Card title={user.Name}/></div>
DREEGO

$DREEGO_BIN generate
grep -q 'user.Name' dreego/gen/routes.go || { echo "FAIL: prop expression user.Name not passed to component"; exit 1; }
go build -o /dev/null .
echo ok
