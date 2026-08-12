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

REPO_DIR="$(cd "$DIR/.." && pwd)"

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

go_failed=0
go_count=0
go_run=0
for pkg in ./core/... ./cli/dreego/...; do
    go_run=$((go_run + 1))
    go_out="$RESULTDIR/gotest-$go_run.out"
    if ! (cd "$REPO_DIR" && go test "$pkg" > "$go_out" 2>&1); then
        go_failed=$((go_failed + 1))
        echo "-> FAIL -> go test $pkg"
        cat "$go_out"
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
goit_out="$RESULTDIR/gotest-goit.out"
if ! (cd "$REPO_DIR" && go test ./_tests/go/... > "$goit_out" 2>&1); then
    echo "-> FAIL -> go test ./_tests/go/..."
    cat "$goit_out"
    FAIL=$((FAIL + 1))
else
    echo "==> PASS <=> _tests/go (Go integration tests, $goit_count) <========="
fi

DREEGO_BIN="$workdir/.dreego-bin"
DREEGO_BIN="$(mktemp -d)/dreego"
(cd "$REPO_DIR" && go build -ldflags "-X main.version=${DREEGO_VERSION:-dev}" -o "$DREEGO_BIN" ./cli/dreego) || {
    echo "FAIL: could not build dreego CLI"
    exit 1
}
export DREEGO_BIN
export REPO_DIR

# Install curl once, sequentially, before any parallel test starts. Every server
# test needs curl; doing apk add concurrently per-test races on apk's database
# lock and flakes. A single pre-install here is deterministic and race-free.
apk add --no-cache curl >/dev/null 2>&1 || { echo "FAIL: curl unavailable"; exit 1; }

DREEGO_PORT_BASE="${DREEGO_PORT_BASE:-20000}"
export DREEGO_PORT_BASE
port_counter=$DREEGO_PORT_BASE

for test_dir in $(find "$DIR/core" -type d 2>/dev/null | sort); do
    test_script="$test_dir/test.sh"
    [ -f "$test_script" ] || continue
    name="${test_dir#$DIR/}"
    if [ -n "$FILTER" ] && ! echo "$name" | grep -Eq "$FILTER"; then
        continue
    fi
    result_file="$RESULTDIR/$(echo "$name" | tr '/' '_')"
    export DREEGO_PORT=$port_counter
    port_counter=$((port_counter + 1))
    (
        if sh "$test_script" >/dev/null 2>&1; then
            echo "PASS" > "$result_file"
        else
            echo "-> FAIL -> $name" > "$result_file"
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
        "-> FAIL"*) echo "$line"; FAIL=$((FAIL + 1)) ;;
    esac
done

echo "==> $([ "$FAIL" -eq 0 ] && echo PASS || echo FAIL) <=> $PASS Passed <=> $FAIL Failed ==="
[ $FAIL -eq 0 ] || exit 1