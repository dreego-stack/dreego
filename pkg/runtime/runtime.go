package runtime

import "net/http"

var routes []route

type route struct {
	method  string
	pattern string
	handler http.HandlerFunc
}

func Register(method, pattern string, handler http.HandlerFunc) {
	routes = append(routes, route{method, pattern, handler})
}

func ServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	for _, r := range routes {
		mux.HandleFunc(r.method+" "+r.pattern, r.handler)
	}
	return mux
}

func Listen(addr string) error {
	return http.ListenAndServe(addr, ServeMux())
}
