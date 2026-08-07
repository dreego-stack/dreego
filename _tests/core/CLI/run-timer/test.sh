#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that dreego run -t 1 starts the server and auto-stops after the timer
set -e

realrepo="$(cd "$(dirname "$0")"/../../../.. && pwd)"
workdir="$(mktemp -d)"
trap "rm -rf $workdir" EXIT

cd "$workdir"

port="${DREEGO_PORT:-$(( ( $(od -An -N2 -i /dev/urandom | tr -d ' ') % 50000 ) + 10000 ))}"

cat > go.mod << EOF
module t
go 1.22
require codeberg.org/dreego/dreego/core v0.0.0
replace codeberg.org/dreego/dreego/core => $realrepo/core
EOF

mkdir -p dreego/routes

cat > dreego/routes/get.dreego << 'DREEGO'
<div><p>timer</p></div>
DREEGO

cat > main.go << GO
package main

import (
	_ "t/dreego/gen"
	dreego "codeberg.org/dreego/dreego/core"
)

func main() {
	dreego.Listen(":$port")
}
GO

apk add --no-cache curl >/dev/null 2>&1 || true

# Start the server with a 1s auto-stop timer; wait until it is up, hit it,
# then verify the process exits on its own. The first go build can be slow.
$DREEGO_BIN run -t 1 > run.log 2>&1 &
pid=$!

up=0
for i in $(seq 1 60); do
    if curl -fsS "http://localhost:$port/" >/dev/null 2>&1; then
        up=1
        break
    fi
    sleep 0.5
done
[ "$up" -eq 1 ] || { echo "FAIL: server did not come up"; kill $pid 2>/dev/null; exit 1; }

curl -fsS "http://localhost:$port/" | grep -q "timer" || { echo "FAIL: server did not respond with expected content"; kill $pid 2>/dev/null; exit 1; }

# Wait for the 1s timer to stop the server and reap the process.
deadline=$(( $(date +%s) + 10 ))
while kill -0 $pid 2>/dev/null; do
    now=$(date +%s)
    [ "$now" -ge "$deadline" ] && { echo "FAIL: server did not auto-stop after timer"; kill $pid 2>/dev/null; exit 1; }
    sleep 0.3
done

grep -q "timer: server stopped" run.log || { echo "FAIL: expected 'timer: server stopped' log line"; exit 1; }

echo ok
