#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that dreego version prints a non-empty version and exits 0
set -e

realrepo="$(cd "$(dirname "$0")"/../../../.. && pwd)"
workdir="$(mktemp -d)"
trap "rm -rf $workdir" EXIT

cd "$workdir"

cat > go.mod << EOF
module t
go 1.22
require github.com/dreego-stack/dreego v0.0.0
replace github.com/dreego-stack/dreego => $realrepo
EOF

out="$($DREEGO_BIN version)"
[ -n "$out" ] || { echo "FAIL: version output is empty"; exit 1; }
[ "$out" = "(devel)" ] && { echo "FAIL: version is the build-info placeholder (devel)"; exit 1; }

echo ok
