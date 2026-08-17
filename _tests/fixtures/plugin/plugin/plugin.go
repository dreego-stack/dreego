package plugin

import (
	"fmt"
	"net/http"

	dreego "github.com/dreego-stack/dreego/core"
)

type Options struct {
	Prefix string
}

func Register(app *dreego.App, options Options) error {
	if options.Prefix == "" {
		options.Prefix = "/plugin"
	}
	if err := app.Register(http.MethodGet, options.Prefix+"/hello", helloHandler); err != nil {
		return err
	}
	if err := app.Register(http.MethodGet, options.Prefix+"/hello/{id}", helloIDHandler); err != nil {
		return err
	}
	return app.Register(http.MethodGet, options.Prefix+"/health", healthHandler)
}

func helloHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("Hello from the plugin"))
}

func helloIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "Hello %s", r.PathValue("id"))
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("plugin ok"))
}
