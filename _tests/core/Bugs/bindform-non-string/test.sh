#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: BindForm returns an error (not panic) on genuinely unsupported field
#       types (e.g. map), while string/int/bool/[]string bind successfully.
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

cat > bindform_test.go << 'GO'
package t

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	dreego "codeberg.org/dreego/dreego/core"
)

type Profile struct {
	Name   string
	Age    int
	Admin  bool
	Labels map[string]string
}

func TestBindFormUnsupportedFieldReturnsError(t *testing.T) {
	form := url.Values{}
	form.Set("name", "Ada")
	form.Set("age", "42")
	form.Set("admin", "on")
	r := &http.Request{
		Method: "POST",
		URL:    &url.URL{Path: "/"},
		Body:   http.NoBody,
		Header: http.Header{"Content-Type": []string{"application/x-www-form-urlencoded"}},
	}
	r.PostForm = form
	r.ParseForm()

	var p Profile
	err := dreego.BindForm(r, &p)
	if err == nil {
		t.Fatal("expected error for unsupported field type (map)")
	}
	if !strings.Contains(err.Error(), "unsupported field type") {
		t.Fatalf("expected unsupported field type error, got: %s", err.Error())
	}
	if p.Name != "Ada" {
		t.Errorf("Name not bound")
	}
	if p.Age != 42 {
		t.Errorf("Age not bound as int")
	}
	if !p.Admin {
		t.Errorf("Admin not bound as bool")
	}
}
GO

go test .
echo ok
