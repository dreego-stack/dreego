package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCSRFConstantTimeCompare(t *testing.T) {
	store := newCSRFMockStore()
	token := seedCSRFToken(t, store)
	mw := CSRF(store)

	run := func(clientToken string) int {
		hit := false
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/", nil)
		r.Header.Set("X-CSRF-Token", clientToken)
		mw(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			hit = true
		})).ServeHTTP(w, r)
		if hit {
			return w.Code
		}
		return w.Code
	}

	if code := run(token); code != http.StatusOK {
		t.Errorf("correct token: expected 200, got %d", code)
	}

	wrongSameLen := token[:len(token)-1] + "X"
	if wrongSameLen == token {
		wrongSameLen = token[:len(token)-1] + "Y"
	}
	if code := run(wrongSameLen); code != http.StatusForbidden {
		t.Errorf("wrong same-length token: expected 403, got %d", code)
	}

	if code := run("wrong-token"); code != http.StatusForbidden {
		t.Errorf("wrong different-length token: expected 403, got %d", code)
	}

	if code := run(""); code != http.StatusForbidden {
		t.Errorf("empty token: expected 403, got %d", code)
	}
}

func TestCSRFConstantTimeCompareLastByteDiff(t *testing.T) {
	store := newCSRFMockStore()
	token := seedCSRFToken(t, store)
	mw := CSRF(store)

	if len(token) < 2 {
		t.Fatalf("seeded token too short: %q", token)
	}

	lastByte := token[len(token)-1]
	diff := "A"
	if string(lastByte) == "A" {
		diff = "B"
	}
	clientToken := token[:len(token)-1] + diff
	if clientToken == token {
		t.Fatalf("setup invariant violated: tokens must differ")
	}

	hit := false
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/", nil)
	r.Header.Set("X-CSRF-Token", clientToken)
	mw(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		hit = true
	})).ServeHTTP(w, r)

	if w.Code != http.StatusForbidden {
		t.Errorf("token differing only in last byte: expected 403, got %d", w.Code)
	}
	if hit {
		t.Error("handler must not be reached when only the last byte differs")
	}
}
