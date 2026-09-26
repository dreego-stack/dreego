package context

import (
	gcontext "context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dreego-stack/dreego/internal/session"
)

func TestFlashRoundTripThroughStore(t *testing.T) {
	t.Parallel()
	store := session.NewCookieStore([]byte("01234567890123456789012345678901"))
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r = session.WithStore(r, store)
	ctx := NewSSR(w, r)
	ctx.Flash("notice", "saved")

	next := httptest.NewRequest(http.MethodGet, "/", nil)
	for _, ck := range w.Result().Cookies() {
		next.AddCookie(ck)
	}
	next = session.WithStore(next, store)
	ctx2 := NewSSR(httptest.NewRecorder(), next)
	if got := ctx2.FlashGet("notice"); got != "saved" {
		t.Fatalf("FlashGet = %q, want saved", got)
	}
	if got := ctx2.FlashPeek("notice"); got != "" {
		t.Fatalf("FlashPeek after consume = %q, want empty", got)
	}
}

func TestCSRFInputEscapesToken(t *testing.T) {
	t.Parallel()
	store := session.NewCookieStore([]byte("01234567890123456789012345678901"))
	w := httptest.NewRecorder()
	r := session.WithStore(httptest.NewRequest(http.MethodGet, "/", nil), store)
	ctx := NewSSR(w, r)
	ctx.SetSessionVal("csrf_token", `a"b<c>`)

	got := ctx.CSRFInput()
	want := `<input type="hidden" name="csrf_token" value="a&#34;b&lt;c&gt;">`
	if got != want {
		t.Fatalf("CSRFInput = %q, want %q", got, want)
	}
}

func TestCSRFInputWithoutStoreIsEmptyValue(t *testing.T) {
	t.Parallel()
	ctx := &SSRContext{Context: gcontext.Background()}
	if got := ctx.CSRFInput(); got != `<input type="hidden" name="csrf_token" value="">` {
		t.Fatalf("CSRFInput = %q", got)
	}
}
