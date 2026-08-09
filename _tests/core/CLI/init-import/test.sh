#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: dreego init scaffolds a buildable project: main.go imports the module-qualified gen package
set -e

realrepo="$(cd "$(dirname "$0")"/../../../.. && pwd)"
workdir="$(mktemp -d)"
trap "rm -rf $workdir" EXIT

cd "$workdir"

if [ -z "$DREEGO_BIN" ]; then
    DREEGO_BIN="$workdir/.dreego-bin"
    (cd "$realrepo" && go build -o "$DREEGO_BIN" ./cli/dreego) || { echo "FAIL: could not build dreego CLI"; exit 1; }
    export DREEGO_BIN
fi

cat > go.mod << EOF
module t
go 1.22
require github.com/dreego-stack/dreego v0.0.0
replace github.com/dreego-stack/dreego => $realrepo
EOF

$DREEGO_BIN init .

# main.go must import the generated package via the module-qualified path
# (module t + dreego/gen), never the bare "gen" import which fails to build
# with "package gen is not in std".
grep -q '_ "t/dreego/gen"' main.go || { echo "FAIL: main.go does not import _ \"t/dreego/gen\""; cat main.go; exit 1; }
if grep -qE '_\s*"gen"' main.go; then
    echo "FAIL: main.go contains bare _ \"gen\" import"
    cat main.go
    exit 1
fi

# A fresh init project must generate and build (GOWORK=off: no go.work in the
# temp dir, the replace directive resolves core locally).
$DREEGO_BIN generate
GOWORK=off go build .

echo ok
