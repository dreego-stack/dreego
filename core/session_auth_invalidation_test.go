package core

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dreego-stack/dreego/core/internal/middleware"
)

func TestAuthLoginReplacesPreAuthState(t *testing.T) {
	store := NewCookieStore(testSecret)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	store.Set(w, r, "pre_auth_cart", "items", nil)
	store.Set(w, r, "pending_oauth_state", "abc123", nil)

	loginW := httptest.NewRecorder()
	loginReq := httptest.NewRequest("GET", "/", nil)
	for _, c := range w.Result().Cookies() {
		loginReq.AddCookie(c)
	}

	store.Set(loginW, loginReq, "user_id", "u42", nil)
	store.Delete(loginW, loginReq, "pending_oauth_state")

	cookies := loginW.Result().Cookies()
	last := cookies[len(cookies)-1]
	verifyReq := httptest.NewRequest("GET", "/", nil)
	verifyReq.AddCookie(last)

	uid, _ := store.Get(verifyReq, "user_id")
	if uid != "u42" {
		t.Errorf("after login: user_id = %q, want u42", uid)
	}
	oauth, _ := store.Get(verifyReq, "pending_oauth_state")
	if oauth != "" {
		t.Errorf("after login: pending_oauth_state should be invalidated, got %q", oauth)
	}
	cart, _ := store.Get(verifyReq, "pre_auth_cart")
	if cart != "items" {
		t.Errorf("after login: pre_auth_cart should be preserved, got %q", cart)
	}
}

func TestPrivilegeChangeReplacesAuthState(t *testing.T) {
	store := NewCookieStore(testSecret)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	store.Set(w, r, "user_id", "u42", nil)
	store.Set(w, r, "role", "member", nil)

	promoteW := httptest.NewRecorder()
	promoteReq := httptest.NewRequest("GET", "/", nil)
	for _, c := range w.Result().Cookies() {
		promoteReq.AddCookie(c)
	}

	store.Set(promoteW, promoteReq, "role", "admin", nil)
	store.Set(promoteW, promoteReq, "privilege_version", "v2", nil)

	cookies := promoteW.Result().Cookies()
	last := cookies[len(cookies)-1]
	verifyReq := httptest.NewRequest("GET", "/", nil)
	verifyReq.AddCookie(last)

	role, _ := store.Get(verifyReq, "role")
	if role != "admin" {
		t.Errorf("after privilege change: role = %q, want admin", role)
	}
	pv, _ := store.Get(verifyReq, "privilege_version")
	if pv != "v2" {
		t.Errorf("after privilege change: privilege_version = %q, want v2", pv)
	}
}

func TestLogoutInvalidatesCompleteAuthState(t *testing.T) {
	store := NewCookieStore(testSecret)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	store.Set(w, r, "user_id", "u42", nil)
	store.Set(w, r, "role", "admin", nil)
	store.Set(w, r, "csrf_token", "tok", nil)

	logoutW := httptest.NewRecorder()
	logoutReq := httptest.NewRequest("GET", "/", nil)
	logoutReq.TLS = &tls.ConnectionState{}
	for _, c := range w.Result().Cookies() {
		logoutReq.AddCookie(c)
	}

	if err := store.Destroy(logoutW, logoutReq); err != nil {
		t.Fatalf("Destroy failed: %v", err)
	}

	cookies := logoutW.Result().Cookies()
	var sess *http.Cookie
	for _, c := range cookies {
		if c.Name == "dreego_session" {
			sess = c
		}
	}
	if sess == nil {
		t.Fatal("expected dreego_session cookie after destroy")
	}
	if sess.MaxAge != -1 {
		t.Errorf("destroy cookie MaxAge = %d, want -1", sess.MaxAge)
	}
	if !sess.Secure {
		t.Error("destroy cookie should preserve Secure=true")
	}
	if !sess.HttpOnly {
		t.Error("destroy cookie should preserve HttpOnly=true")
	}
	if sess.SameSite == http.SameSiteDefaultMode {
		t.Error("destroy cookie should preserve SameSite")
	}
	if sess.Path != "/" {
		t.Errorf("destroy cookie Path = %q, want /", sess.Path)
	}
}

func TestReplayedCookieRejectedAfterRotation(t *testing.T) {
	old := NewCookieStore(testSecret)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	old.Set(w, r, "user_id", "u42", &Options{Encrypt: true})

	newStore := NewCookieStore(testSecret2)

	replayReq := httptest.NewRequest("GET", "/", nil)
	for _, c := range w.Result().Cookies() {
		replayReq.AddCookie(c)
	}

	uid, _ := newStore.Get(replayReq, "user_id")
	if uid != "" {
		t.Errorf("replayed cookie signed with old secret should be rejected, got %q", uid)
	}
}

func TestReplayedSignedOnlyCookieRejectedAfterRotation(t *testing.T) {
	old := NewCookieStore(testSecret)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	old.Set(w, r, "user_id", "u42", nil)

	newStore := NewCookieStore(testSecret2)

	replayReq := httptest.NewRequest("GET", "/", nil)
	for _, c := range w.Result().Cookies() {
		replayReq.AddCookie(c)
	}

	uid, _ := newStore.Get(replayReq, "user_id")
	if uid != "" {
		t.Errorf("replayed signed-only cookie with old secret should be rejected, got %q", uid)
	}
}

func TestInvalidSessionReturnsError(t *testing.T) {
	store := NewCookieStore(testSecret)
	r := httptest.NewRequest("GET", "/", nil)
	r.AddCookie(&http.Cookie{Name: "dreego_session", Value: "garbage-value-not-valid"})

	_, err := store.Get(r, "user_id")
	if err == nil {
		t.Error("invalid session cookie should return observable error, got nil")
	}
}

func TestLoginAfterRotationEstablishesNewSession(t *testing.T) {
	old := NewCookieStore(testSecret)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	old.Set(w, r, "user_id", "old-user", nil)
	old.Set(w, r, "pending_oauth_state", "abc123", nil)

	newStore := NewCookieStore(testSecret2)
	loginReq := httptest.NewRequest("GET", "/", nil)
	for _, c := range w.Result().Cookies() {
		loginReq.AddCookie(c)
	}
	loginW := httptest.NewRecorder()

	if err := newStore.Set(loginW, loginReq, "user_id", "new-user", nil); err != nil {
		t.Fatalf("login after rotation must succeed, got %v", err)
	}

	verifyReq := httptest.NewRequest("GET", "/", nil)
	for _, c := range loginW.Result().Cookies() {
		verifyReq.AddCookie(c)
	}
	uid, err := newStore.Get(verifyReq, "user_id")
	if err != nil {
		t.Fatalf("Get after re-login failed: %v", err)
	}
	if uid != "new-user" {
		t.Errorf("user_id = %q, want new-user", uid)
	}
	oauth, err := newStore.Get(verifyReq, "pending_oauth_state")
	if err != nil {
		t.Fatalf("Get pending_oauth_state failed: %v", err)
	}
	if oauth != "" {
		t.Errorf("old keys must not survive rotation, got pending_oauth_state=%q", oauth)
	}
}

func TestCSRFPostAfterRotationSucceeds(t *testing.T) {
	old := NewCookieStore(testSecret)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	old.Set(w, r, "csrf_token", "stale-token", nil)
	stale := findCookie(t, w.Result().Cookies(), "dreego_session")

	newStore := NewCookieStore(testSecret2)
	mw := middleware.CSRF(newStore)

	getW := httptest.NewRecorder()
	getR := httptest.NewRequest("GET", "/", nil)
	getR.AddCookie(stale)
	mw(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})).ServeHTTP(getW, getR)

	readable := csrfReadablePolicyCookie(t, getW.Result().Cookies())
	if readable == nil {
		t.Fatal("expected readable csrf_token cookie")
	}

	postW := httptest.NewRecorder()
	postR := httptest.NewRequest("POST", "/", nil)
	for _, c := range getW.Result().Cookies() {
		postR.AddCookie(c)
	}
	postR.Header.Set("X-CSRF-Token", readable.Value)
	ok := false
	mw(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { ok = true })).ServeHTTP(postW, postR)
	if !ok {
		t.Error("CSRF POST after rotation should succeed")
	}
	if postW.Code != http.StatusOK {
		t.Errorf("CSRF POST status = %d, want 200", postW.Code)
	}
}
