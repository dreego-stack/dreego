#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
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

go run codeberg.org/dreego/dreego/cmd/dreego init .
[ -f main.go ] || { echo "missing main.go"; exit 1; }
[ -f dreego/routes/get.dreego ] || { echo "missing get.dreego"; exit 1; }
echo ok
