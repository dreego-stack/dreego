package i18n

import "testing"

func TestLocalizeDocumentAddsMetadata(t *testing.T) {
	cases := []struct {
		document string
		locale   string
		want     string
	}{
		{"<!doctype html><html><body></body></html>", "de", `<!doctype html><html lang="de"><body></body></html>`},
		{`<HTML class="app" lang="de"><body></body></HTML>`, "en-GB", `<HTML class="app" lang="en-GB"><body></body></HTML>`},
		{"<html><body></body></html>", "ar", `<html lang="ar" dir="rtl"><body></body></html>`},
	}
	for _, tc := range cases {
		if got := LocalizeDocument(tc.document, tc.locale); got != tc.want {
			t.Errorf("LocalizeDocument(%q) = %q, want %q", tc.locale, got, tc.want)
		}
	}
}

func TestIsRTL(t *testing.T) {
	if !IsRTL("ar") || !IsRTL("he-IL") || IsRTL("de") {
		t.Fatal("unexpected text direction")
	}
}
