#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that dreego run -t uses SIGTERM (graceful shutdown) instead of SIGKILL (B20)
set -e

realrepo="$(cd "$(dirname "$0")"/../../.. && pwd)"
workdir="$(mktemp -d)"
trap "rm -rf $workdir" EXIT

cd "$workdir"

cat > go.mod << EOF
module t
go 1.22
require codeberg.org/dreego/dreego v0.0.0
replace codeberg.org/dreego/dreego => $realrepo
EOF

cat > main.go << 'GO'
package main
import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	_ "t/dreego/gen"
	core "codeberg.org/dreego/dreego/core"
)

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-done
		fmt.Println("SIGTERM received")
		os.Exit(0)
	}()
	core.SetLogging(false)
	core.Listen(":8080")
}
GO

mkdir -p dreego/routes
cat > dreego/routes/get.dreego << 'DREEGO'
<div>hello</div>
DREEGO

go run $realrepo/cmd/dreego generate
go run $realrepo/cmd/dreego run -t 1 2>&1 | grep -q "SIGTERM received" || { echo "FAIL: server did not receive SIGTERM (B20)"; exit 1; }
echo ok
