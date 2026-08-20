#!/bin/sh
set -e
DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_DIR="$(cd "$DIR/.." && pwd)"
RUNS="${DREEGO_RUNS:-1}"
case "$RUNS" in
    ''|*[!0-9]*) RUNS=1 ;;
esac
FILTER="${DREEGO_FILTER:-}"
JOBS="${DREEGO_JOBS:-$(nproc 2>/dev/null || echo 4)}"
DREEGO_PORT_BASE="${DREEGO_PORT_BASE:-20000}"
RESULTDIR="$(mktemp -d)"
trap "rm -rf $RESULTDIR" EXIT

# Build the CLI once; every run reuses it.
DREEGO_BIN="$(mktemp -d)/dreego"
(cd "$REPO_DIR" && go build -ldflags "-X main.version=${DREEGO_VERSION:-dev}" -o "$DREEGO_BIN" ./cli/dreego) || {
    echo "FAIL: could not build dreego CLI"
    exit 1
}
export DREEGO_BIN
export REPO_DIR
export DREEGO_LOCAL_REPO="$REPO_DIR"

# Install curl once, before any parallel test starts. Every server test needs
# curl; doing apk add concurrently per-test races on apk's database lock and
# flakes. A single pre-install here is deterministic and race-free.
apk add --no-cache curl >/dev/null 2>&1 || { echo "FAIL: curl unavailable"; exit 1; }

run_suite() {
    run=$1
    PASS=0
    FAIL=0
    RUNNING=0
    port_counter=$((DREEGO_PORT_BASE + (run - 1) * 1000))
    run_dir="$RESULTDIR/run-$run"
    mkdir -p "$run_dir"

    if ! (cd "$REPO_DIR" && sh _scripts/check-core-deps.sh > /dev/null 2>&1); then
        echo "FAIL core-deps"
        FAIL=$((FAIL + 1))
    else
        echo "==> PASS <=> Core deps <========="
    fi

    if ! out=$(sh "$DIR/find-binary.sh" 2>&1); then
        echo "$out" | grep '^->' || true
        echo "==> FAIL   <=>  find-binary <==========="
        FAIL=$((FAIL + 1))
    else
        echo "==> PASS <=> No binary files <========="
    fi

    if ! (cd "$REPO_DIR" && python3 _scripts/release-prep-test.py > "$run_dir/release-prep-test.out" 2>&1); then
        echo "-> FAIL -> release-prep contract tests"
        cat "$run_dir/release-prep-test.out"
        FAIL=$((FAIL + 1))
    else
        echo "==> PASS <=> release-prep contract tests <========="
    fi

    go_failed=0
    go_count=0
    go_run=0
    for pkg in ./core/... ./cli/dreego/...; do
        go_run=$((go_run + 1))
        go_out="$run_dir/gotest-$go_run.out"
        if ! (cd "$REPO_DIR" && go list "$pkg" > /dev/null 2>&1); then
            go_failed=$((go_failed + 1))
            echo "-> FAIL -> missing package $pkg"
            continue
        fi
        if ! (cd "$REPO_DIR" && go test -v "$pkg" > "$go_out" 2>&1); then
            go_failed=$((go_failed + 1))
            echo "-> FAIL -> go test $pkg"
            cat "$go_out"
        fi
        if grep -q -- '--- SKIP:' "$go_out" 2>/dev/null; then
            go_failed=$((go_failed + 1))
            echo "-> FAIL -> skipped tests in $pkg"
            grep -- '--- SKIP:' "$go_out"
        fi
        go_count=$((go_count + $(cd "$REPO_DIR" && go test -list '^Test' "$pkg" 2>/dev/null | grep -c '^Test')))
    done

    if [ "$go_failed" -gt 0 ]; then
        echo "==> FAIL   <=>  GO Tests ($go_count) <==========="
        FAIL=$((FAIL + go_failed))
    else
        echo "==> PASS <=> GO Tests ($go_count) <========="
    fi

    goit_count=$(cd "$REPO_DIR" && go test -list '^Test' ./_tests/go/... 2>/dev/null | grep -c '^Test')
    goit_out="$run_dir/gotest-goit.out"
    if ! (cd "$REPO_DIR" && go list ./_tests/go/... > /dev/null 2>&1); then
        echo "-> FAIL -> missing package ./_tests/go/..."
        FAIL=$((FAIL + 1))
    elif ! (cd "$REPO_DIR" && go test -v ./_tests/go/... > "$goit_out" 2>&1); then
        echo "-> FAIL -> go test ./_tests/go/..."
        cat "$goit_out"
        FAIL=$((FAIL + 1))
    else
        echo "==> PASS <=> _tests/go (Go integration tests, $goit_count) <========="
    fi
    if grep -q -- '--- SKIP:' "$goit_out" 2>/dev/null; then
        if grep -v -- '--- SKIP: TestCLIVersionDrift' "$goit_out" | grep -q -- '--- SKIP:'; then
            echo "-> FAIL -> skipped tests in ./_tests/go/..."
            grep -- '--- SKIP:' "$goit_out"
            FAIL=$((FAIL + 1))
        fi
    fi

    wait

    for f in "$run_dir"/*; do
        [ -f "$f" ] || continue
        read -r line < "$f"
        case "$line" in
            PASS) PASS=$((PASS + 1)) ;;
            FAIL*) echo "$line"; FAIL=$((FAIL + 1)) ;;
            "-> FAIL"*) echo "$line"; FAIL=$((FAIL + 1)) ;;
        esac
    done

    echo "==> $([ "$FAIL" -eq 0 ] && echo PASS || echo FAIL) <=> $PASS Passed <=> $FAIL Failed ==="
    [ "$FAIL" -eq 0 ]
}

if [ "$RUNS" -gt 1 ]; then
    echo "==> Running suite $RUNS times (DREEGO_RUNS) <========="
fi

run=1
while [ "$run" -le "$RUNS" ]; do
    echo "==> Run $run/$RUNS <========="
    if ! run_suite "$run"; then
        echo "==> FAILED on run $run/$RUNS <========="
        exit 1
    fi
    run=$((run + 1))
done

if [ "$RUNS" -gt 1 ]; then
    echo "==> ALL $RUNS RUNS PASSED <========="
fi
