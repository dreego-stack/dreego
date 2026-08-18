package core

import (
	"net/http"
	"testing"
)

func TestRedirectRejectsRootWildcardTarget(t *testing.T) {
	app := New()
	if err := app.RegisterRedirect("/old/*", "/*", http.StatusFound); err == nil {
		t.Fatal("root wildcard target can emit a protocol-relative redirect")
	}
}

func TestRewriteRejectsRootWildcardTarget(t *testing.T) {
	app := New()
	if err := app.RegisterRewrite("/old/*", "/*"); err == nil {
		t.Fatal("root wildcard target can emit a double-slash path")
	}
}
