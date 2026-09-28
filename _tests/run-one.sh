#!/bin/sh
# Run exactly one test folder. Invoked by test.sh.
# Args: <folder> <name> <result-dir>
# Env: DREEGO_RACE (1 enables -race for Go tests), DREEGO_BIN.
set -u

dir="$1"
name="$2"
outdir="$3"
safe="$(printf '%s' "$name" | tr '/' '_')"
status_file="$outdir/$safe.status"

output=""
status=0
if [ -f "$dir/test.sh" ]; then
    output="$(cd "$dir" && sh test.sh 2>&1)" || status=1
fi
if ls "$dir"/*_test.go >/dev/null 2>&1; then
    if [ "${DREEGO_RACE:-1}" = "1" ]; then
        go_output="$(cd "$dir" && go test -race -count=1 . 2>&1)" || status=1
    else
        go_output="$(cd "$dir" && go test -count=1 . 2>&1)" || status=1
    fi
    case "$go_output" in
        *"build constraints exclude all Go files"*)
            status=0
            go_output="skipped: build constraints exclude all Go files"
            ;;
    esac
    output="$output
$go_output"
fi

if [ "$status" -eq 0 ]; then
    printf 'PASS\n' > "$status_file"
else
    printf 'FAIL\n' > "$status_file"
    {
        printf '=== FAIL %s ===\n' "$name"
        printf '%s\n' "$output"
    } >> "$outdir/failures.log"
fi
