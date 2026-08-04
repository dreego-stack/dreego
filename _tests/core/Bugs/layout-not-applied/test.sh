#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Route with layout renders inside layout {#slot}
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

mkdir -p dreego/layouts dreego/routes

cat > dreego/layouts/default.dreego << 'DREEGO'
<html>
<head><title>Site</title></head>
<body>
<nav>Nav</nav>
<main>{#slot}</main>
<footer>Footer</footer>
</body>
</html>
DREEGO

cat > dreego/routes/get.dreego << 'DREEGO'
<div><h1>Page</h1></div>
DREEGO

cat > main.go << 'GO'
package main
import (_ "t/dreego/gen"; core "codeberg.org/dreego/dreego/core")
func main() { core.Listen(":0") }
GO

$DREEGO_BIN generate

grep -q '<html>' dreego/gen/routes.go || { echo "FAIL: layout not applied — no <html> in generated routes"; exit 1; }
grep -q 'Nav' dreego/gen/routes.go || { echo "FAIL: layout nav missing in generated routes"; exit 1; }
grep -q 'Footer' dreego/gen/routes.go || { echo "FAIL: layout footer missing in generated routes"; exit 1; }
grep -q 'Page' dreego/gen/routes.go || { echo "FAIL: route content missing in generated routes"; exit 1; }

go build -o /dev/null .
echo ok
