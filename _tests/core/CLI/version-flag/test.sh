#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: dreego --version and -v print the version, like the version subcommand
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

# check_version runs the CLI with the given flag and asserts exit 0 plus a
# non-empty, non-placeholder version string. Prints the output.
check_version() {
    flag="$1"
    if ! out="$($DREEGO_BIN $flag 2>&1)"; then
        echo "FAIL: dreego $flag exited non-zero"
        exit 1
    fi
    [ -n "$out" ] || { echo "FAIL: dreego $flag output is empty"; exit 1; }
    [ "$out" = "(devel)" ] && { echo "FAIL: dreego $flag output is the build-info placeholder (devel)"; exit 1; }
    echo "$out"
}

flag_out="$(check_version --version)"
short_out="$(check_version -v)"
sub_out="$(check_version version)"

# --version and -v must behave exactly like the version subcommand.
[ "$flag_out" = "$sub_out" ] || { echo "FAIL: --version and version outputs differ: '$flag_out' vs '$sub_out'"; exit 1; }
[ "$short_out" = "$sub_out" ] || { echo "FAIL: -v and version outputs differ: '$short_out' vs '$sub_out'"; exit 1; }

echo ok
