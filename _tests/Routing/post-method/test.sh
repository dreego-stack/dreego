#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that a POST method route compiles
set -e

realrepo="$(cd "$(dirname "$0")"/../../.. && pwd)"
workdir="$(mktemp -d)"
trap "rm -rf $workdir" EXIT

cd "$workdir"

cat > go.mod << EOF
module t
go 1.22
require codeberg.org/dreego/dreego v0.0.0
replace codeberg.org/dreego/dreego => $realrepo
EOF

go run $realrepo/cmd/dreego init .
sed -i 's|_ "gen"|_ "t/dreego/gen"|' main.go
go run $realrepo/cmd/dreego generate
go build -o /dev/null .
echo ok
