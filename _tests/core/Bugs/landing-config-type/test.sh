#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that landing blueprint config.json has correct logging type (B8)
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

go run $realrepo/cmd/dreego new testapp 2>&1

cd testapp
if ! grep -q '"logging": {' dreego/config.json; then
    echo "FAIL: landing config.json logging field has wrong type (B8)"
    exit 1
fi
if ! grep -q '"enabled":' dreego/config.json; then
    echo "FAIL: landing config.json logging.enabled missing (B8)"
    exit 1
fi
echo ok
