package transpiler

import "testing"

func TestMimeByExt(t *testing.T) {
	cases := map[string]string{
		".css":   "text/css; charset=utf-8",
		".js":    "application/javascript; charset=utf-8",
		".svg":   "image/svg+xml",
		".png":   "image/png",
		".ico":   "image/x-icon",
		".jpg":   "image/jpeg",
		".jpeg":  "image/jpeg",
		".html":  "text/html; charset=utf-8",
		".json":  "application/json; charset=utf-8",
		".woff2": "font/woff2",
		".woff":  "font/woff",
	}
	for ext, want := range cases {
		if got := MimeByExt(ext); got != want {
			t.Errorf("MimeByExt(%q) = %q, want %q", ext, got, want)
		}
	}
}

func TestMimeByExtDefault(t *testing.T) {
	for _, ext := range []string{".txt", ".md", ".bin", "", ".unknown"} {
		if got := MimeByExt(ext); got != "application/octet-stream" {
			t.Errorf("MimeByExt(%q) = %q, want application/octet-stream", ext, got)
		}
	}
}

func TestMimeByExtCaseInsensitive(t *testing.T) {
	for _, ext := range []string{".CSS", ".Js", ".PNG", ".HTML", ".Woff2"} {
		if got := MimeByExt(ext); got == "application/octet-stream" {
			t.Errorf("MimeByExt(%q) returned default, expected case-insensitive match", ext)
		}
	}
}
