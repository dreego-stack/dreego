#!/bin/sh
set -e
DIR="$(cd "$(dirname "$0")" && pwd)"
PASS=0
FAIL=0
FILTER="${DREEGO_FILTER:-}"

for test_dir in $(find "$DIR" -type d -not -path "$DIR" | sort); do
    test_script="$test_dir/test.sh"
    [ -f "$test_script" ] || continue
    name="${test_dir#$DIR/}"
    if [ -n "$FILTER" ] && ! echo "$name" | grep -Eq "$FILTER"; then
        continue
    fi
    if sh "$test_script" >/dev/null 2>&1; then
        PASS=$((PASS + 1))
    else
        echo "FAIL $name"
        FAIL=$((FAIL + 1))
    fi
done

echo "=== $PASS passed, $FAIL failed ==="
[ $FAIL -eq 0 ] || exit 1
