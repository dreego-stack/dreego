#!/bin/sh
# Recursive test orchestrator. Every folder under _tests/e2e that contains a
# test.sh or a *_test.go file is one test. All tests run in parallel; each runs
# with its own folder as the working directory.
set -e

DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_DIR="$(cd "$DIR/.." && pwd)"
E2E_DIR="$DIR/e2e"

RUNS="${DREEGO_RUNS:-1}"
case "$RUNS" in ''|*[!0-9]*) RUNS=1 ;; esac
FILTER="${DREEGO_FILTER:-}"
JOBS="${DREEGO_JOBS:-$(nproc 2>/dev/null || echo 4)}"
RACE="${DREEGO_RACE:-1}"

export REPO_DIR DREEGO_LOCAL_REPO="$REPO_DIR" DREEGO_RACE="$RACE"

DREEGO_BIN_DIR="$(mktemp -d)"
DREEGO_BIN="$DREEGO_BIN_DIR/dreego"
RESULTDIR="$(mktemp -d)"
trap 'rm -rf "$DREEGO_BIN_DIR" "$RESULTDIR"' EXIT
export DREEGO_BIN

VERSION="${DREEGO_VERSION:-$(cd "$REPO_DIR" && git describe --tags --match 'v[0-9]*.[0-9]*.[0-9]*' --abbrev=0 2>/dev/null || echo dev)}"
(cd "$REPO_DIR" && go build -ldflags "-X main.version=$VERSION" -o "$DREEGO_BIN" ./cmd/dreego) || {
    echo "FAIL: could not build dreego CLI"
    exit 1
}

"$DREEGO_BIN" tools install typescript >/dev/null 2>&1 || {
    echo "FAIL: could not install pinned TypeScript compiler"
    exit 1
}

if ! command -v curl >/dev/null 2>&1; then
    apk add --no-cache curl >/dev/null 2>&1 || { echo "FAIL: curl unavailable"; exit 1; }
fi

discover_tests() {
    find "$E2E_DIR" -type f \( -name 'test.sh' -o -name '*_test.go' \) 2>/dev/null |
        sed "s#^$E2E_DIR/##" |
        awk -F/ '{ sub(/\/[^\/]+$/, "", $0); print }' |
        sort -u
}

run_suite() {
    run="$1"
    RUN_DIR="$RESULTDIR/run-$run"
    mkdir -p "$RUN_DIR"
    : > "$RUN_DIR/failures.log"

    tests="$(discover_tests)"
    if [ -n "$FILTER" ]; then
        tests="$(printf '%s\n' "$tests" | grep -E "$FILTER" || true)"
    fi
    total="$(printf '%s\n' "$tests" | grep -c . || true)"

    printf '==> Run %s: %s tests, %s jobs (race=%s)\n' "$run" "$total" "$JOBS" "$RACE"

    if [ "$total" -gt 0 ]; then
        printf '%s\n' "$tests" |
            xargs -P "$JOBS" -I{} sh "$DIR/run-one.sh" "$E2E_DIR/{}" "{}" "$RUN_DIR"
    fi

    pass=0
    fail=0
    for f in "$RUN_DIR"/*.status; do
        [ -f "$f" ] || continue
        if [ "$(head -n 1 "$f")" = "PASS" ]; then pass=$((pass + 1)); else fail=$((fail + 1)); fi
    done

    if [ -s "$RUN_DIR/failures.log" ]; then
        cat "$RUN_DIR/failures.log"
    fi
    printf '==> %s <=> %s passed <=> %s failed\n' \
        "$([ "$fail" -eq 0 ] && echo PASS || echo FAIL)" "$pass" "$fail"
    [ "$fail" -eq 0 ]
}

run=1
while [ "$run" -le "$RUNS" ]; do
    if ! run_suite "$run"; then
        echo "==> FAILED on run $run/$RUNS"
        exit 1
    fi
    run=$((run + 1))
done

echo "==> ALL $RUNS RUNS PASSED"
