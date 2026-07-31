#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that dreego new testapp creates a project from the landing blueprint
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

[ -d testapp ] || { echo "FAIL: testapp directory not created"; exit 1; }
[ -f testapp/main.go ] || { echo "FAIL: missing main.go"; exit 1; }
[ -f testapp/go.mod ] || { echo "FAIL: missing go.mod"; exit 1; }
[ -f testapp/dreego/routes/get.dreego ] || { echo "FAIL: missing get.dreego"; exit 1; }
[ -f testapp/dreego/layouts/default.dreego ] || { echo "FAIL: missing default.dreego"; exit 1; }
[ -f testapp/dreego/components/Hero.dreego ] || { echo "FAIL: missing Hero.dreego"; exit 1; }
[ -f testapp/dreego/components/FeatureCard.dreego ] || { echo "FAIL: missing FeatureCard.dreego"; exit 1; }

echo ok
