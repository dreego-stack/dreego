#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: dreego.ServeMux() returns the same handler on repeated calls (cached stack)
set -e

realrepo="$(cd "$(dirname "$0")"/../../../.. && pwd)"
workdir="$(mktemp -d)"
trap "rm -rf $workdir" EXIT

cd "$workdir"

cat > go.mod << EOF
module t
go 1.22
require codeberg.org/dreego/dreego/core v0.0.0
replace codeberg.org/dreego/dreego/core => $realrepo/core
EOF

cat > main.go << 'GO'
package main

import (
	"fmt"

	dreego "codeberg.org/dreego/dreego/core"
	_ "t/dreego/gen"
)

func main() {
	first := fmt.Sprintf("%p", dreego.ServeMux())
	second := fmt.Sprintf("%p", dreego.ServeMux())
	if first != second {
		fmt.Println("FAIL: ServeMux returned different handlers")
		panic("servemux cache miss")
	}
	fmt.Println("ok")
}
GO

mkdir -p dreego/routes

cat > dreego/routes/get.dreego << 'DREEGO'
<div><p>hello</p></div>
DREEGO

$DREEGO_BIN generate
go build -o app .
./app
