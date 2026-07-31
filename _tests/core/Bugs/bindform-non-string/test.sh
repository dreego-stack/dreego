#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that BindForm returns an error (not panic) on non-string fields (B2)
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
	"net/http"
	"net/url"
	"strings"
	core "codeberg.org/dreego/dreego/core"
)

type Profile struct {
	Name  string
	Age   int
	Admin bool
}

func main() {
	form := url.Values{}
	form.Set("name", "Ada")
	form.Set("age", "42")
	form.Set("admin", "true")
	r := &http.Request{
		Method: "POST",
		URL:    &url.URL{Path: "/"},
		Body:   http.NoBody,
		Header: http.Header{"Content-Type": []string{"application/x-www-form-urlencoded"}},
	}
	r.PostForm = form
	r.ParseForm()

	var p Profile
	err := core.BindForm(r, &p)
	if err == nil {
		panic("expected error for non-string fields")
	}
	if !strings.Contains(err.Error(), "unsupported field type") {
		panic("expected unsupported field type error, got: " + err.Error())
	}
	if p.Name != "Ada" {
		panic("Name not bound")
	}
}
GO

go run .
echo ok
