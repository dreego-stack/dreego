#!/bin/sh
set -e
DIR="$(cd "$(dirname "$0")" && pwd)"
PASS=0
FAIL=0
FILTER="${DREEGO_FILTER:-}"
FAILS=""

echo "=== dreego _tests ==="
echo

for test_dir in $(find "$DIR" -type d -not -path "$DIR" | sort); do
    test_script="$test_dir/test.sh"
    [ -f "$test_script" ] || continue

    name="${test_dir#$DIR/}"

    if [ -n "$FILTER" ] && ! echo "$name" | grep -Eq "$FILTER"; then
        continue
    fi

    out=$(cd "$test_dir" && sh test.sh 2>&1)
    rc=$?

    if [ $rc -eq 0 ]; then
        echo "  PASS  $name"
        PASS=$((PASS + 1))
    else
        echo "  FAIL  $name"
        FAILS="$FAILS\n  FAIL  $name\n$(echo "$out" | sed 's/^/         /')\n"
        FAIL=$((FAIL + 1))
    fi
done

if [ $FAIL -gt 0 ]; then
    echo
    echo "FAILURES:"
    echo -e "$FAILS"
fi

echo
echo "=== $PASS passed, $FAIL failed ==="
[ $FAIL -eq 0 ] || exit 1
