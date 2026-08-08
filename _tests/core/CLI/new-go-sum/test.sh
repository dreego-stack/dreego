#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: fresh dreego new scaffold builds without go mod tidy
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

$DREEGO_BIN new testapp 2>&1

[ -d testapp ] || { echo "FAIL: testapp directory not created"; exit 1; }

cd testapp

# `dreego new` adds a `replace` directive pointing the required dreego version
# at the local repo root, so `go mod tidy` and the build run fully offline.
# Verify the replace directive was written and that the scaffold builds.
grep -q '^replace github.com/dreego-stack/dreego => ' go.mod \
    || { echo "FAIL: go.mod has no replace directive for the local dreego"; exit 1; }

GOWORK=off go mod tidy

# With a local-dir replace there is no remote module to verify, so no go.sum is
# written. The offline proof is that the scaffold generates and builds.
$DREEGO_BIN generate
GOWORK=off go build .

echo "ok"
