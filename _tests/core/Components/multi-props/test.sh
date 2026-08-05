#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Component with 3+ props all rendered
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
require codeberg.org/dreego/dreego/core v0.0.0
replace codeberg.org/dreego/dreego/core => $realrepo/core
EOF

cat > main.go << 'GO'
package main
import _ "t/dreego/gen"
func main() {}
GO

mkdir -p dreego/components dreego/routes
cat > dreego/components/Profile.dreego << 'DREEGO'
Component Profile (name string, role string, email string)
<div><h2>{name}</h2><p>{role}</p><a href="mailto:{email}">{email}</a></div>
DREEGO
cat > dreego/routes/get.dreego << 'DREEGO'
<div><@Profile name="Ada" role="Admin" email="ada@example.com"/></div>
DREEGO

$DREEGO_BIN generate
grep -q 'Profile("Ada"' dreego/gen/routes.go || { echo "FAIL: first prop not passed"; exit 1; }
grep -q '"Admin"' dreego/gen/routes.go || { echo "FAIL: second prop not passed"; exit 1; }
grep -q '"ada@example.com"' dreego/gen/routes.go || { echo "FAIL: third prop not passed"; exit 1; }
go build -o /dev/null .
echo ok
