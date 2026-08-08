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
    (cd "$realrepo" && go build -o "$DREEGO_BIN" ./cmd/dreego) || { echo "FAIL: could not build dreego CLI"; exit 1; }
    export DREEGO_BIN
fi

# The CLI version must match the repo's VERSION file (v0.0.25), not a stale
# core version (v0.0.23) that cmd/dreego/go.mod currently requires.
want="$(cat "$realrepo/VERSION")"
out="$($DREEGO_BIN version)"
case "$out" in
    *"$want"*) ;;
    *) echo "FAIL: dreego version is '$out', want it to contain '$want' (VERSION file)"; exit 1 ;;
esac

# The CLI module must require a core version matching the repo's VERSION
# file. A stale require (e.g. v0.0.23 while VERSION is v0.0.25) is the drift.
# In-repo builds resolve core via go.work (local core wins); external
# `go install` resolves via the published core/vX.Y.Z tag.
gomod="$realrepo/cmd/dreego/go.mod"
if ! grep -q "^require github.com/dreego-stack/dreego/core $want\$" "$gomod"; then
    echo "FAIL: cmd/dreego/go.mod does not require github.com/dreego-stack/dreego/core $want"
    grep 'github.com/dreego-stack/dreego/core' "$gomod" || true
    exit 1
fi

echo ok
