package dreefile

import (
	"mime"
	"strings"
)

// mimeByExt pins the extensions whose value Dreego needs to control or that the
// host MIME database does not map consistently across platforms (fonts and
// icons). Every other known extension falls through to the Go standard library.
var mimeByExt = map[string]string{
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

func MimeByExt(ext string) string {
	ext = strings.ToLower(ext)
	if value, ok := mimeByExt[ext]; ok {
		return value
	}
	if value := mime.TypeByExtension(ext); value != "" {
		return value
	}
	return "application/octet-stream"
}
