#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: dreego docs reads local module docs from go.mod (vendor/cache), -p
#       for plugins, and --list for the sitemap. No HTTP, no embedded copy.
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

# 1) Core docs: run inside the dreego repo itself (module == self). `dreego docs`
#    reads the realrepo _docs/ directly; --json emits structured JSON.
out="$(cd "$realrepo" && "$DREEGO_BIN" docs 2>&1)"
echo "$out" | grep -q "Dreego Documentation" || { echo "FAIL: docs did not show core index"; exit 1; }
out="$(cd "$realrepo" && "$DREEGO_BIN" docs --json 2>&1)"
echo "$out" | grep -q '"headings"' || { echo "FAIL: docs --json missing headings"; exit 1; }
out="$(cd "$realrepo" && "$DREEGO_BIN" docs --list 2>&1)"
echo "$out" | grep -q "github.com/dreego-stack/dreego" || { echo "FAIL: --list missing core"; exit 1; }

# 2) Plugin docs: a fake project requiring plugin-sse with the plugin under
#    vendor/. `dreego docs -p plugin-sse` and `--list` resolve it locally.
mkdir -p vendor/github.com/dreego-stack/plugin-sse/_docs
cat > go.mod <<EOF
module example.com/myapp

go 1.22

require (
    github.com/dreego-stack/dreego v0.0.27
    github.com/dreego-stack/plugin-sse v0.1.0
)
EOF
cat > vendor/github.com/dreego-stack/plugin-sse/_docs/index.md <<EOF
# Plugin SSE

Vendor-local plugin docs.
EOF

out="$("$DREEGO_BIN" docs -p plugin-sse 2>&1)"
echo "$out" | grep -q "Vendor-local plugin docs" || { echo "FAIL: docs -p did not show vendor plugin docs"; exit 1; }

out="$("$DREEGO_BIN" docs --list 2>&1)"
echo "$out" | grep -q "github.com/dreego-stack/plugin-sse" || { echo "FAIL: --list missing plugin"; exit 1; }

echo ok
