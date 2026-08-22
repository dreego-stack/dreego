#!/bin/sh
set -eu

threshold="${DREEGO_COVERAGE_MIN:-35}"
output="$(go test -cover ./core/... 2>&1)"
printf '%s\n' "$output"

if ! printf '%s\n' "$threshold" | awk '/^[0-9]+(\.[0-9]+)?$/ { ok=1 } END { exit ok ? 0 : 1 }'; then
    echo "error: DREEGO_COVERAGE_MIN must be a non-negative number" >&2
    exit 1
fi

if printf '%s\n' "$output" | awk -v minimum="$threshold" '
    /coverage: [0-9.]+% of statements/ {
        value = $0
        sub(/^.*coverage: /, "", value)
        sub(/%.*$/, "", value)
        if (value < minimum) {
            failed = 1
            print "error: coverage " value "% is below " minimum "%" > "/dev/stderr"
        }
    }
    END { exit failed ? 1 : 0 }
'; then
    echo "coverage gate passed (minimum ${threshold}%)"
else
    exit 1
fi
