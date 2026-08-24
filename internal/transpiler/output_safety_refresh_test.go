package transpiler

import (
	"strings"
	"testing"
)

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
		{`<client src="{{ u }}"></client>`, "SafeURL"},
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
