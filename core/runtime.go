package core

import (
	"net/http"
	"strings"
)

type route struct {
	method  string
	pattern string
	handler http.HandlerFunc
}

type redirectRule struct {
	from   string
	to     string
	status int
}

type rewriteRule struct {
	from string
	to   string
}

func matchRewrite(rw rewriteRule, path string) bool {
	return strings.HasPrefix(path, strings.TrimSuffix(rw.from, "/*"))
}

func matchRedirect(rd redirectRule, path string) (string, bool) {
	if strings.HasSuffix(rd.from, "/*") {
		prefix := strings.TrimSuffix(rd.from, "/*")
		if strings.HasPrefix(path, prefix) {
			return strings.Replace(path, prefix, strings.TrimSuffix(rd.to, "/*"), 1), true
		}
		return "", false
	}
	if path == rd.from {
		return rd.to, true
	}
	return "", false
}
