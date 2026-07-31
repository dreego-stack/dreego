#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that validate min/max rules reject non-digit values (B17)
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
	core "codeberg.org/dreego/dreego/core"
)

type Form struct {
	Name string `validate:"min=abc"`
}

func main() {
	errs := core.ValidateForm(Form{Name: "x"})
	if errs == nil || errs["name"] == "" {
		fmt.Println("FAIL: non-digit min rule silently accepted")
		return
	}
	fmt.Println("ok")
}
GO

go run .
