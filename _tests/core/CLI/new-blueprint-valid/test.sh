#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that dreego new scaffolds a valid project: no unreplaced placeholders, generate succeeds, config.json is valid JSON
set -e

realrepo="$(cd "$(dirname "$0")"/../../../.. && pwd)"
workdir="$(mktemp -d)"
trap "rm -rf $workdir" EXIT

cd "$workdir"

cat > go.mod << EOF
module t
go 1.22
require github.com/dreego-stack/dreego/core v0.0.0
replace github.com/dreego-stack/dreego/core => $realrepo/core
EOF

$DREEGO_BIN new testapp >/dev/null 2>&1

# No unreplaced placeholder may remain anywhere in the scaffold.
if grep -rq '§\$name\$§' testapp; then
    echo "FAIL: unreplaced §\$name\$§ placeholder found in scaffold"
    exit 1
fi

# config.json must be valid JSON (use python3 when present, else grep fallback).
if command -v python3 >/dev/null 2>&1; then
    python3 -c "import json,sys; json.load(open(sys.argv[1]))" testapp/dreego/config.json \
        || { echo "FAIL: config.json is not valid JSON"; exit 1; }
else
    cat testapp/dreego/config.json | grep -q '"logging"' || { echo "FAIL: config.json invalid/missing logging"; exit 1; }
fi

# dreego generate must succeed inside the scaffold.
(
    cd testapp
    "$DREEGO_BIN" generate >/dev/null 2>&1
) || { echo "FAIL: dreego generate failed in scaffold"; exit 1; }

[ -f testapp/dreego/gen/routes.go ] || { echo "FAIL: gen/routes.go not produced"; exit 1; }

echo ok
