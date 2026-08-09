#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: .gitignore does not ignore dreego/routes/, dreego/components/, or dreego/config.json
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

$DREEGO_BIN new testapp 2>&1

[ -d testapp ] || { echo "FAIL: testapp directory not created"; exit 1; }
[ -f testapp/.gitignore ] || { echo "FAIL: missing .gitignore"; exit 1; }

gitignore_lines() {
    sed -E 's/^[[:space:]]+//; s/[[:space:]]+$//' testapp/.gitignore | grep -v '^#' | grep -v '^$'
}

if gitignore_lines | grep -Eqx '/dreego|dreego|dreego/'; then
    echo "FAIL: .gitignore contains a top-level 'dreego' ignore pattern"
    exit 1
fi

if gitignore_lines | grep -Eq '^/dreego([[:space:]]+.*)?$'; then
    echo "FAIL: .gitignore contains '/dreego' which ignores the source directory"
    exit 1
fi
if gitignore_lines | grep -Eq '^dreego/[[:space:]]*$'; then
    echo "FAIL: .gitignore contains 'dreego/' which ignores the source directory"
    exit 1
fi
if gitignore_lines | grep -Eq 'dreego[[:space:]]*$'; then
    echo "FAIL: .gitignore contains trailing 'dreego' which ignores the source directory"
    exit 1
fi

echo "ok"
