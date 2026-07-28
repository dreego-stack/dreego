package core

import (
	"net/http"
	"path/filepath"
	"strings"

)

func RegisterStatic(path, mime string, content []byte) {
	data := make([]byte, len(content))
	copy(data, content)

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", mime)
		w.Write(data)
	}
	routes = append(routes, route{method: "GET", pattern: path, handler: handler})
}

func MimeByExt(ext string) string {
	switch strings.ToLower(ext) {
	case ".css":
		return "text/css; charset=utf-8"
	case ".js":
		return "application/javascript; charset=utf-8"
	case ".svg":
		return "image/svg+xml"
	case ".png":
		return "image/png"
	case ".ico":
		return "image/x-icon"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".html":
		return "text/html; charset=utf-8"
	case ".json":
		return "application/json; charset=utf-8"
	case ".woff2":
		return "font/woff2"
	case ".woff":
		return "font/woff"
	default:
		return "application/octet-stream"
	}
}

func staticPattern(filePath string) string {
	return "/" + filepath.ToSlash(filePath)
}
