package sample

import (
	"net/http"

	dreego "codeberg.org/dreego/dreego/core"
)

type SamplePlugin struct{}

func (p *SamplePlugin) Middleware(next http.Handler) http.Handler {
	_ = dreego.NewSSR
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}