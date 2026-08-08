#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that dreego run -t uses SIGTERM (graceful shutdown) instead of SIGKILL (B20)
set -e

realrepo="$(cd "$(dirname "$0")"/../../../.. && pwd)"
workdir="$(mktemp -d)"
trap "rm -rf $workdir" EXIT

cd "$workdir"

if [ -z "$DREEGO_BIN" ]; then
    DREEGO_BIN="$workdir/.dreego-bin"
    (cd "$realrepo" && go build -o "$DREEGO_BIN" ./cmd/dreego) || { echo "FAIL: could not build dreego CLI"; exit 1; }
    export DREEGO_BIN
fi

cat > go.mod << EOF
module t
go 1.22
require github.com/dreego-stack/dreego v0.0.0
replace github.com/dreego-stack/dreego => $realrepo
EOF

port="${DREEGO_PORT:-$(( ( $(od -An -N2 -i /dev/urandom | tr -d ' ') % 50000 ) + 10000 ))}"

cat > main.go << GO
package main
import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	_ "t/dreego/gen"
	dreego "github.com/dreego-stack/dreego/core"
)

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-done
		fmt.Println("SIGTERM received")
		os.Exit(0)
	}()
	dreego.SetLogging(false)
	dreego.Listen(":$port")
}
GO

mkdir -p dreego/routes
cat > dreego/routes/get.dreego << 'DREEGO'
<div>hello</div>
DREEGO

$DREEGO_BIN generate
outfile="$workdir/run.out"
$DREEGO_BIN run -t 15 > "$outfile" 2>&1
grep -q "SIGTERM received" "$outfile" || { echo "FAIL: server did not receive SIGTERM (B20)"; cat "$outfile"; exit 1; }
echo ok
