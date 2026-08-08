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
require github.com/dreego-stack/dreego/core v0.0.0
replace github.com/dreego-stack/dreego/core => $realrepo/core
EOF

cat > validate_test.go << 'GO'
package t

import (
	"testing"

	dreego "github.com/dreego-stack/dreego/core"
)

type Form struct {
	Name string `validate:"min=abc"`
}

func TestValidateNonDigitRejected(t *testing.T) {
	errs := dreego.ValidateForm(Form{Name: "x"})
	if errs == nil || errs["name"] == "" {
		t.Fatal("non-digit min rule silently accepted")
	}
}
GO

go test .
echo ok
