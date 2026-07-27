package core

import (
	"strings"
)

type RouteInfo struct {
	HandlerName string
	RoutePath   string
	Method      string
}

func GenerateRouter(routes []RouteInfo) string {
	var buf strings.Builder

	buf.WriteString("package main\n\n")
	buf.WriteString("import \"net/http\"\n\n")
	buf.WriteString("func RegisterRoutes(mux *http.ServeMux) {\n")
	for _, r := range routes {
		buf.WriteString("\tmux.HandleFunc(\"" + r.Method + " " + r.RoutePath + "\", " + r.HandlerName + ")\n")
	}
	buf.WriteString("}\n")

	return buf.String()
}
