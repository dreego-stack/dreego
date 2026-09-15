#!/bin/sh
set -eu

threshold="${DREEGO_COVERAGE_MIN:-35}"
profile="$(mktemp)"
trap 'rm -f "$profile"' EXIT
output="$(go test -coverprofile="$profile" ./internal/... ./core/... ./adapter/ssr/... ./adapter/wails/... ./dreegotest/... ./cmd/dreego/... 2>&1)"
printf '%s\n' "$output"

if ! printf '%s\n' "$threshold" | awk '/^[0-9]+(\.[0-9]+)?$/ { ok=1 } END { exit ok ? 0 : 1 }'; then
    echo "error: DREEGO_COVERAGE_MIN must be a non-negative number" >&2
    exit 1
fi

total="$(awk '
    NR > 1 { statements += $2; if (($3 + 0) > 0) covered += $2 }
    END { if (statements > 0) printf "%.1f\n", (covered * 100) / statements }
' "$profile")"
if [ -z "$total" ]; then
    total="$(printf '%s\n' "$output" | awk '/coverage: [0-9.]+% of statements/ { value=$0; sub(/^.*coverage: /, "", value); sub(/%.*$/, "", value); print value; exit }')"
fi

if printf '%s\n' "$total" | awk -v minimum="$threshold" '
    /^[0-9]+(\.[0-9]+)?$/ { valid=1; if (($1 + 0) < (minimum + 0)) exit 1 }
    END { if (!valid) exit 2 }
'; then
    echo "coverage gate passed (minimum ${threshold}%)"
else
    echo "error: coverage ${total:-unknown}% is below ${threshold}%" >&2
    exit 1
fi
