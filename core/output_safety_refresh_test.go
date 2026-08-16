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

func TestHeadSafeFuncClassifies(t *testing.T) {
	cases := []struct {
		content string
		want    string
	}{
		{`<title>{{ t }}</title>`, "SafeText"},
		{`<meta name="description" content="{{ t }}">`, "SafeAttr"},
		{`<link href="{{ u }}">`, "SafeURL"},
		{`<meta http-equiv="refresh" content="{{ u }}">`, "SafeRefresh"},
		{`<meta http-equiv=refresh content="{{ u }}">`, "SafeRefresh"},
		{`<script src="{{ u }}"></script>`, "SafeURL"},
	}
	for _, c := range cases {
		i := strings.Index(c.content, "{{")
		if i < 0 {
			t.Fatalf("no placeholder in %q", c.content)
		}
		if got := headSafeFunc(c.content, i); got != c.want {
			t.Errorf("headSafeFunc(%q) = %q, want %q", c.content, got, c.want)
		}
	}
}

func TestHeadSafeFuncClassifiesWhitespaceAroundEquals(t *testing.T) {
	cases := []struct {
		content string
		want    string
	}{
		{`<meta http-equiv = "refresh" content="{{ u }}">`, "SafeRefresh"},
		{`<meta http-equiv = refresh content="{{ u }}">`, "SafeRefresh"},
		{`<meta HTTP-EQUIV="refresh" content="{{ u }}">`, "SafeRefresh"},
		{`<meta http-equiv = "refresh" CONTENT = "{{ u }}">`, "SafeRefresh"},
	}
	for _, c := range cases {
		i := strings.Index(c.content, "{{")
		if i < 0 {
			t.Fatalf("no placeholder in %q", c.content)
		}
		if got := headSafeFunc(c.content, i); got != c.want {
			t.Errorf("headSafeFunc(%q) = %q, want %q", c.content, got, c.want)
		}
	}
}

func TestHeadSafeFuncClassifiesUnquotedValue(t *testing.T) {
	cases := []struct {
		content string
		want    string
	}{
		{`<link href={{ u }}>`, "SafeURL"},
		{`<meta name=description content={{ t }}>`, "SafeAttr"},
		{`<meta http-equiv=refresh content={{ u }}>`, "SafeRefresh"},
	}
	for _, c := range cases {
		i := strings.Index(c.content, "{{")
		if i < 0 {
			t.Fatalf("no placeholder in %q", c.content)
		}
		if got := headSafeFunc(c.content, i); got != c.want {
			t.Errorf("headSafeFunc(%q) = %q, want %q", c.content, got, c.want)
		}
	}
}

func TestAttrValueWhitespaceTolerance(t *testing.T) {
	cases := []struct {
		tag, attr, want string
	}{
		{`<meta http-equiv="refresh">`, "http-equiv", "refresh"},
		{`<meta http-equiv = "refresh">`, "http-equiv", "refresh"},
		{`<meta HTTP-EQUIV = refresh>`, "http-equiv", "refresh"},
		{`<meta http-equiv=refresh>`, "http-equiv", "refresh"},
		{`<meta charset="utf-8" http-equiv = 'x'>`, "http-equiv", "x"},
		{`<meta name="description">`, "http-equiv", ""},
	}
	for _, c := range cases {
		if got := attrValue(c.tag, c.attr); got != c.want {
			t.Errorf("attrValue(%q, %q) = %q, want %q", c.tag, c.attr, got, c.want)
		}
	}
}

func TestGenHeadURLAttrUsesSafeURL(t *testing.T) {
	out, err := genHead(`<link href="{{ u }}">`, "b")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "dreego.SafeURL") {
		t.Errorf("genHead link href must use SafeURL, got:\n%s", out)
	}
}

func TestGenHeadMetaRefreshUsesSafeRefresh(t *testing.T) {
	out, err := genHead(`<meta http-equiv="refresh" content="{{ u }}">`, "b")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "dreego.SafeRefresh") {
		t.Errorf("genHead meta refresh content must use SafeRefresh, got:\n%s", out)
	}
}
