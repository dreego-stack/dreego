package dreegotest_test

import (
	"net/http"
	"net/url"
	"testing"

	dreego "github.com/dreego-stack/dreego/core"
	"github.com/dreego-stack/dreego/dreegotest"
)

func TestAssertHelpers(t *testing.T) {
	ok := &testing.T{}
	dreegotest.MustStatus(ok, 200, 200)
	dreegotest.MustContainBody(ok, "hello world", "hello")
	dreegotest.MustNotContainBody(ok, "hello world", "goodbye")
	dreegotest.MustEqual(ok, "a", "a")
	dreegotest.MustEqual(ok, 1, 1)
	dreegotest.MustNotEqual(ok, "a", "b")
	dreegotest.MustHeader(ok, http.Header{"X-Test": {"v"}}, "X-Test", "v")
	h := http.Header{"Location": {"/next"}}
	dreegotest.MustRedirect(ok, 302, h, "/next")
}

func TestMustStatusFails(t *testing.T) {
	f := failT{}
	dreegotest.MustStatus(&f, 200, 404)
	if !f.failed {
		t.Fatal("MustStatus must fail on mismatch")
	}
}

func TestMustContainBodyFails(t *testing.T) {
	f := failT{}
	dreegotest.MustContainBody(&f, "abc", "xyz")
	if !f.failed {
		t.Fatal("MustContainBody must fail on missing substring")
	}
}

func TestMustNotContainBodyFails(t *testing.T) {
	f := failT{}
	dreegotest.MustNotContainBody(&f, "abc", "ab")
	if !f.failed {
		t.Fatal("MustNotContainBody must fail on present substring")
	}
}

func TestMustEqualFails(t *testing.T) {
	f := failT{}
	dreegotest.MustEqual(&f, 1, 2)
	if !f.failed {
		t.Fatal("MustEqual must fail on mismatch")
	}
}

func TestMustNotEqualFails(t *testing.T) {
	f := failT{}
	dreegotest.MustNotEqual(&f, 1, 1)
	if !f.failed {
		t.Fatal("MustNotEqual must fail on equality")
	}
}

func TestMustHeaderFails(t *testing.T) {
	f := failT{}
	dreegotest.MustHeader(&f, http.Header{}, "X-Missing", "v")
	if !f.failed {
		t.Fatal("MustHeader must fail on missing header")
	}
}

func TestMustRedirectFails(t *testing.T) {
	f := failT{}
	dreegotest.MustRedirect(&f, 200, http.Header{}, "/next")
	if !f.failed {
		t.Fatal("MustRedirect must fail on non-redirect status")
	}
}

func TestResponseIncludesHeaderAndCookies(t *testing.T) {
	app := dreego.New()
	app.Register("GET", "/dgt-resp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		http.SetCookie(w, &http.Cookie{Name: "sid", Value: "abc"})
		w.WriteHeader(http.StatusFound)
		w.Header().Set("Location", "/landing")
	})
	resp := dreegotest.NewApp(app).Get(t, "/dgt-resp")
	if resp.StatusCode != http.StatusFound {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusFound)
	}
	dreegotest.MustHeader(t, resp.Header, "Content-Type", "text/plain")
	dreegotest.MustHeader(t, resp.Header, "Location", "/landing")
	if resp.Location() != "/landing" {
		t.Errorf("Location() = %q, want /landing", resp.Location())
	}
	if len(resp.Cookies) == 0 || resp.Cookies[0].Name != "sid" {
		t.Errorf("Cookies = %v, want sid cookie", resp.Cookies)
	}
}

func TestAppClientPostForm(t *testing.T) {
	app := dreego.New()
	app.Register("POST", "/dgt-submit", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("got:" + r.FormValue("name")))
	})
	resp := dreegotest.NewApp(app).PostForm(t, "/dgt-submit", url.Values{"name": {"world"}})
	dreegotest.MustStatus(t, resp.StatusCode, 200)
	dreegotest.MustContainBody(t, resp.Body, "got:world")
}

type failT struct {
	testing.TB
	failed bool
}

func (f *failT) Helper() {}
func (f *failT) Fatalf(format string, args ...any) {
	f.failed = true
}
