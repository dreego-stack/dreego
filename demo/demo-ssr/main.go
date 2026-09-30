package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	dreego "github.com/dreego-stack/dreego/core"
	"demo/dreego/blog"
	luaDemo "demo/dreego/lua"
	"demo/dreego/saas"
	"demo/dreego/www"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	public := dreego.New(www.App)
	product := dreego.New(saas.App)
	blogApp := dreego.New(blog.App)
	luaApp := dreego.New(luaDemo.App)

	if err := errors.Join(
		configure(public),
		configure(product),
		configure(blogApp),
		configure(luaApp),
		registerLocaleSelection(public),
	); err != nil {
		return err
	}

	handler := hostRouter(public.Handler(), product.Handler(), blogApp.Handler(), luaApp.Handler())
	addr := ":8080"
	if port := os.Getenv("DREEGO_PORT"); port != "" {
		addr = ":" + port
	}
	return http.ListenAndServe(addr, handler)
}

func registerLocaleSelection(app *dreego.App) error {
	return app.Register(http.MethodPost, "/locale", dreego.LocaleSelectionHandler(dreego.LocaleSelectionOptions{
		FallbackPath: "/",
	}))
}

func configure(app *dreego.App) error {
	return app.SetCSP("default-src 'self'; script-src 'self' 'unsafe-inline' https://cdn.tailwindcss.com; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' https:; connect-src 'self'; base-uri 'self'; form-action 'self'")
}

func hostRouter(public, product, blog, lua http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host := strings.Split(r.Host, ":")[0]
		if host == "saas.localhost" || strings.HasPrefix(host, "saas.") {
			product.ServeHTTP(w, r)
			return
		}
		if host == "blog.localhost" || strings.HasPrefix(host, "blog.") {
			blog.ServeHTTP(w, r)
			return
		}
		if host == "lua.localhost" || strings.HasPrefix(host, "lua.") {
			lua.ServeHTTP(w, r)
			return
		}
		path := r.URL.Path
		if path == "/_dreego/lua.js" {
			lua.ServeHTTP(w, r)
			return
		}
		if path == "/blog" {
			http.Redirect(w, r, "/blog/", http.StatusMovedPermanently)
			return
		}
		if path == "/saas" {
			http.Redirect(w, r, "/saas/", http.StatusMovedPermanently)
			return
		}
		if path == "/lua" {
			http.Redirect(w, r, "/lua/", http.StatusMovedPermanently)
			return
		}
		if strings.HasPrefix(path, "/blog/") {
			http.StripPrefix("/blog", blog).ServeHTTP(w, r)
			return
		}
		if strings.HasPrefix(path, "/saas/") {
			http.StripPrefix("/saas", product).ServeHTTP(w, r)
			return
		}
		if strings.HasPrefix(path, "/lua/") {
			http.StripPrefix("/lua", lua).ServeHTTP(w, r)
			return
		}
		public.ServeHTTP(w, r)
	})
}
