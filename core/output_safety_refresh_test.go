package core

import (
	"strings"
	"testing"
)

func TestSafeRefreshRejectsUnsafeURL(t *testing.T) {
	if got := SafeRefresh(`0;url=javascript:alert(1)`); strings.Contains(got, "javascript:") {
		t.Errorf("SafeRefresh must reject javascript: URL, got %q", got)
	}
	if got := SafeRefresh(`5`); got != "5" {
		t.Errorf("SafeRefresh without url= must pass through, got %q", got)
	}
	if got := SafeRefresh(`0;url=https://example.com`); !strings.Contains(got, "https://example.com") {
		t.Errorf("SafeRefresh must keep https URLs, got %q", got)
	}
}

func TestSafeRefreshRejectsWhitespaceAndCaseBypass(t *testing.T) {
	cases := []struct{ in, want string }{
		{"0; url = javascript:alert(1)", "0; url = #"},
		{"0;url =javascript:alert(1)", "0;url =#"},
		{"0;URL = javascript:alert(1)", "0;URL = #"},
		{"0;Url\t=\tjavascript:alert(1)", "0;Url\t=\t#"},
		{"0; url=javascript:alert(1)", "0; url=#"},
		{"0;url =\n javascript:alert(1)", "0;url =\n #"},
		{"0;url=java\nscript:alert(1)", "0;url=#"},
		{"0;url=java\rscript:alert(1)", "0;url=#"},
		{"0;url=java\tscript:alert(1)", "0;url=#"},
	}
	for _, c := range cases {
		got := SafeRefresh(c.in)
		if got != c.want {
			t.Errorf("SafeRefresh(%q) = %q, want sanitized %q", c.in, got, c.want)
		}
		if strings.Contains(got, "javascript:") {
			t.Errorf("SafeRefresh(%q) = %q, must reject javascript: value", c.in, got)
		}
	}
}

func TestSafeRefreshKeepsWhitespaceVariantWithSafeURL(t *testing.T) {
	for _, in := range []string{
		"0; url = https://example.com",
		"0;URL=https://example.com",
	} {
		got := SafeRefresh(in)
		if !strings.Contains(got, "https://example.com") {
			t.Errorf("SafeRefresh(%q) = %q, want https URL kept", in, got)
		}
		if strings.Contains(got, "javascript:") {
			t.Errorf("SafeRefresh(%q) = %q, must not contain javascript:", in, got)
		}
	}
}

func TestSafeSrcdocKeepsMarkupInertAfterNestedParse(t *testing.T) {
	got := SafeSrcdoc(`<script>parent.postMessage("run", "*")</script>`)
	if strings.Contains(got, "&lt;script") {
		t.Fatalf("srcdoc is escaped only once: %q", got)
	}
	if !strings.Contains(got, "&amp;lt;script&amp;gt;") {
		t.Fatalf("srcdoc does not preserve escaped markup for its nested document: %q", got)
	}
}
