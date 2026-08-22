#!/bin/sh
# import-check: builds the importcheck fixture, starts the server, and verifies
# that the Tailwind CDN script declared in the page <head> is served in the
# rendered HTML. This is the black-box counterpart to the head-merging unit
# tests in _tests/go/bug_layout_head_ok_test.go.
set -e

DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_DIR="$(cd "$DIR/../.." && pwd)"
FIXTURE_DIR="$REPO_DIR/_tests/fixtures/importcheck"

# Temp working copy so generated dree.go files and build artifacts never touch
# the fixture source. Also lets multiple CI runs run in parallel.
WORK="$(mktemp -d)"
trap 'rm -rf "$WORK"' EXIT

cp -R "$FIXTURE_DIR/." "$WORK/"
sed -i.bak 's|^replace github.com/dreego-stack/dreego => .*$|replace github.com/dreego-stack/dreego => '"$REPO_DIR"'|' "$WORK/go.mod"
rm -f "$WORK/go.mod.bak"

DREEGO_BIN="${DREEGO_BIN:-}"
if [ -z "$DREEGO_BIN" ]; then
    DREEGO_BIN="$(mktemp -d)/dreego"
    (cd "$REPO_DIR" && go build -o "$DREEGO_BIN" ./cli/dreego)
fi

(cd "$WORK" && "$DREEGO_BIN" generate)
(cd "$WORK" && go build -o server .)

PORT="${PORT:-$(shuf -i 20000-29999 -n 1 2>/dev/null || echo 20000)}"
(cd "$WORK" && PORT="$PORT" ./server) &
SERVER_PID=$!
trap 'kill "$SERVER_PID" 2>/dev/null; rm -rf "$WORK"' EXIT

for _ in $(seq 1 50); do
    if curl -sf "http://127.0.0.1:$PORT" >/dev/null 2>&1; then
        break
    fi
    sleep 0.1
done

BODY="$(curl -sf "http://127.0.0.1:$PORT")" || {
    echo "FAIL: server did not respond on :$PORT"
    exit 1
}

case "$BODY" in
    *'<script src="https://cdn.tailwindcss.com"></script>'*)
        echo "PASS: Tailwind CDN script present in rendered HTML"
        ;;
    *)
        echo "FAIL: Tailwind CDN script missing from rendered HTML"
        exit 1
        ;;
esac
