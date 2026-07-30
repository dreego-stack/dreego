#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that static files in subdirectories compile
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

mkdir -p dreego/routes dreego/static/images

cat > dreego/routes/get.dreego << 'DREEGO'
<div><p>hello</p></div>
DREEGO

printf '<svg></svg>' > dreego/static/images/logo.svg

go run $realrepo/cmd/dreego generate

grep -q 'RegisterStatic("/images/logo.svg"' dreego/gen/dree.go
echo "ok"
