package session

import (
	"net/http"
	"testing"
)

func findCookie(t *testing.T, cookies []*http.Cookie, name string) *http.Cookie {
	t.Helper()
	for _, c := range cookies {
		if c.Name == name {
			return c
		}
	}
	t.Fatalf("expected cookie %s, got %d cookies", name, len(cookies))
	return nil
}
