package tests

import (
	"bytes"
	"net/http/httptest"
	"strings"
	"testing"

	dreego "github.com/dreego-stack/dreego/core"
)

func TestBugSessionWriteOrdering(t *testing.T) {
	t.Parallel()
	store := dreego.NewCookieStore(bytes.Repeat([]byte("s"), 32))
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	if err := store.Set(w, r, "ok", "value", nil); err != nil {
		t.Fatalf("first Set: %v", err)
	}
	if err := store.Set(w, r, "big", strings.Repeat("x", 8192), nil); err != dreego.ErrSessionTooLarge {
		t.Fatalf("second Set: got %v, want ErrSessionTooLarge", err)
	}
	if got, _ := store.Get(r, "big"); got != "" {
		t.Errorf("oversized value must not be visible after failed write, got %q", got)
	}
}
