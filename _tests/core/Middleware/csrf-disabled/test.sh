#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that CSRF middleware can be disabled at runtime
set -e
realrepo="$(cd "$(dirname "$0")"/../../../.. && pwd)"
workdir="$(mktemp -d)"
trap "rm -rf $workdir" EXIT

cd "$workdir"

cat > go.mod << EOF
module t
go 1.22
require codeberg.org/dreego/dreego v0.0.0
replace codeberg.org/dreego/dreego => $realrepo
EOF

mkdir -p dreego/routes
cat > dreego/routes/get.dreego << 'DREEGO'
<div><p>csrf off</p></div>
DREEGO
cat > main.go << 'GO'
package main
import (_ "t/dreego/gen"; core "codeberg.org/dreego/dreego/core")
func main() { core.SetSessionStore(core.NewCookieStore([]byte("test"))); core.SetCSRF(false); core.Listen(":0") }
GO
go run $realrepo/cmd/dreego generate
go build -o /dev/null .
echo ok
