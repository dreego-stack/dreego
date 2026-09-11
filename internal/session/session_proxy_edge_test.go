package session

import (
	"net/http/httptest"
	"testing"
)

func TestIsTLSRejectsUntrustedForwardedProto(t *testing.T) {
	r := httptest.NewRequest("GET", "http://example.test/", nil)
	r.RemoteAddr = "203.0.113.10:443"
	r.Header.Set("X-Forwarded-Proto", "https")
	if IsTLS(r, map[string]bool{"198.51.100.5": true}) {
		t.Fatal("IsTLS trusted an untrusted proxy")
	}
}

func TestIsTLSAcceptsTrustedProxyWithPortAndHTTPS(t *testing.T) {
	r := httptest.NewRequest("GET", "http://example.test/", nil)
	r.RemoteAddr = "198.51.100.5:8080"
	r.Header.Set("X-Forwarded-Proto", "https")
	if !IsTLS(r, map[string]bool{"198.51.100.5": true}) {
		t.Fatal("IsTLS rejected trusted HTTPS proxy")
	}
}

func TestIsTLSDoesNotAcceptNonHTTPSForwardedProto(t *testing.T) {
	r := httptest.NewRequest("GET", "http://example.test/", nil)
	r.RemoteAddr = "198.51.100.5:8080"
	r.Header.Set("X-Forwarded-Proto", "http")
	if IsTLS(r, map[string]bool{"198.51.100.5": true}) {
		t.Fatal("IsTLS accepted non-HTTPS forwarded proto")
	}
}
