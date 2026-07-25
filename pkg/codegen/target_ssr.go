package codegen

import (
	"codeberg.org/dreego/dreego/pkg/context"
	"net/http"
)

func NewHandler(render func(*context.SSRContext) (string, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.NewSSR(w, r)
		html, err := render(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(html))
	}
}
