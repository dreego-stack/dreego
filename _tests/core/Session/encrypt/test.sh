#!/bin/sh
# What: encrypted session round-trip via the core module
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

cat > session_encrypt_test.go << 'GO'
package t

import (
	"encoding/base64"
	"net/http/httptest"
	"strings"
	"testing"

	core "codeberg.org/dreego/dreego/core"
)

func TestSessionEncryptRoundTrip(t *testing.T) {
	store := core.NewCookieStore([]byte("secret-key"))
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	store.Set(w, r, "user", "alice", &core.Options{Encrypt: true})

	req := httptest.NewRequest("GET", "/", nil)
	for _, c := range w.Result().Cookies() {
		req.AddCookie(c)
	}
	v, _ := store.Get(req, "user")
	if v != "alice" {
		t.Fatalf("round-trip failed, got %q", v)
	}
}

func TestSessionEncryptValueNotPlaintext(t *testing.T) {
	store := core.NewCookieStore([]byte("secret-key"))
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	store.Set(w, r, "user", "alice", &core.Options{Encrypt: true})

	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(cookies))
	}
	decoded, err := base64.RawURLEncoding.DecodeString(cookies[0].Value)
	if err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if strings.Contains(string(decoded), "alice") {
		t.Error("encrypted payload must not contain plaintext value")
	}
}

func TestSessionEncryptTamperRejected(t *testing.T) {
	store := core.NewCookieStore([]byte("secret-key"))
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	store.Set(w, r, "role", "admin", &core.Options{Encrypt: true})

	req := httptest.NewRequest("GET", "/", nil)
	for _, c := range w.Result().Cookies() {
		c.Value = base64.RawURLEncoding.EncodeToString([]byte("tampered"))
		req.AddCookie(c)
	}
	v, _ := store.Get(req, "role")
	if v != "" {
		t.Errorf("tampered cookie rejected, got %q", v)
	}
}

func TestSessionEncryptKeyRotationRejected(t *testing.T) {
	store := core.NewCookieStore([]byte("secret-key"))
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	store.Set(w, r, "role", "admin", &core.Options{Encrypt: true})

	other := core.NewCookieStore([]byte("other-secret"))
	req := httptest.NewRequest("GET", "/", nil)
	for _, c := range w.Result().Cookies() {
		req.AddCookie(c)
	}
	v, _ := other.Get(req, "role")
	if v != "" {
		t.Errorf("wrong key should reject cookie, got %q", v)
	}
}
GO

go test .
echo ok
