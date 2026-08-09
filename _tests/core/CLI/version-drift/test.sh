#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: CLI builds against the current core version, not a stale one
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

# The CLI version must match the version injected at build time (the latest
# git tag), not a stale core version.
want="${DREEGO_VERSION:-dev}"
out="$($DREEGO_BIN version)"
case "$out" in
    *"$want"*) ;;
    *) echo "FAIL: dreego version is '$out', want it to contain '$want' (DREEGO_VERSION)"; exit 1 ;;
esac

# The repo module must declare the module path matching the repo's latest tag.
# A stale require (e.g. v0.0.23 while the latest tag is v0.0.27) is the drift.
# In-repo builds resolve core via the root module (local core wins); external
# `go install` resolves via the published vX.Y.Z tag.
gomod="$realrepo/go.mod"
if ! grep -q "^module github.com/dreego-stack/dreego$" "$gomod"; then
    echo "FAIL: go.mod does not declare module github.com/dreego-stack/dreego"
    grep '^module ' "$gomod" || true
    exit 1
fi

echo ok
