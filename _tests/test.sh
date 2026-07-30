#!/bin/sh
set -e
DIR="$(cd "$(dirname "$0")" && pwd)"
PASS=0
FAIL=0
FILTER="${DREEGO_FILTER:-}"
JOBS="${DREEGO_JOBS:-$(nproc 2>/dev/null || echo 4)}"
RESULTDIR="$(mktemp -d)"
trap "rm -rf $RESULTDIR" EXIT
RUNNING=0

for test_dir in $(find "$DIR" -type d -not -path "$DIR" | sort); do
    test_script="$test_dir/test.sh"
    [ -f "$test_script" ] || continue
    name="${test_dir#$DIR/}"
    if [ -n "$FILTER" ] && ! echo "$name" | grep -Eq "$FILTER"; then
        continue
    fi
    result_file="$RESULTDIR/$(echo "$name" | tr '/' '_')"
    (
        if sh "$test_script" >/dev/null 2>&1; then
            echo "PASS" > "$result_file"
        else
            echo "FAIL $name" > "$result_file"
        fi
    ) &
    RUNNING=$((RUNNING + 1))
    if [ "$RUNNING" -ge "$JOBS" ]; then
        wait -n 2>/dev/null || true
        RUNNING=$((RUNNING - 1))
    fi
done

wait

for f in "$RESULTDIR"/*; do
    [ -f "$f" ] || continue
    read -r line < "$f"
    case "$line" in
        PASS) PASS=$((PASS + 1)) ;;
        FAIL*) echo "$line"; FAIL=$((FAIL + 1)) ;;
    esac
done

echo "=== $PASS passed, $FAIL failed ==="
[ $FAIL -eq 0 ] || exit 1
