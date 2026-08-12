package tests

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	dreego "github.com/dreego-stack/dreego/core"
)

func TestBugBindFormNonString(t *testing.T) {
	t.Parallel()
	type Profile struct {
		Name   string
		Age    int
		Admin  bool
		Labels map[string]string
	}
	form := url.Values{}
	form.Set("name", "Ada")
	form.Set("age", "42")
	form.Set("admin", "on")
	form.Set("labels", "x")
	req, _ := http.NewRequest("POST", "/", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	var p Profile
	if err := dreego.BindForm(req, &p); err == nil {
		t.Fatal("expected error for unsupported map field, got nil")
	}
}
