#!/bin/sh
set -e
DIR="$(cd "$(dirname "$0")" && pwd)"
PASS=0
FAIL=0

echo "=== dreego _tests ==="
echo

for test_dir in $(find "$DIR" -type d -not -path "$DIR" | sort); do
    test_script="$test_dir/test.sh"
    [ -f "$test_script" ] || continue

    name="${test_dir#$DIR/}"
    printf "  %-45s " "$name"

    out=$(cd "$test_dir" && sh test.sh 2>&1)
    rc=$?

    if [ $rc -eq 0 ]; then
        echo "PASS"
        PASS=$((PASS + 1))
    else
        echo "FAIL"
        echo "$out" | sed 's/^/      /'
        FAIL=$((FAIL + 1))
    fi
done

echo
echo "=== $PASS passed, $FAIL failed ==="
[ $FAIL -eq 0 ] || exit 1
