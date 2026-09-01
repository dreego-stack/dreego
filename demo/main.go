package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"demo/blog"
	"demo/saas"
	"demo/www"
	dreego "github.com/dreego-stack/dreego/core"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	public := dreego.New()
	if err := configure(public); err != nil {
		return err
	}
	if err := www.Register(public); err != nil {
		return err
	}

	product := dreego.New()
	if err := configure(product); err != nil {
		return err
	}
	if err := saas.Register(product); err != nil {
		return err
	}

	blogApp := dreego.New()
	if err := configure(blogApp); err != nil {
		return err
	}
	if err := blog.Register(blogApp); err != nil {
		return err
	}

	handler := hostRouter(public.Handler(), product.Handler(), blogApp.Handler())
	addr := ":8080"
	if port := os.Getenv("DREEGO_PORT"); port != "" {
		addr = ":" + port
	}
	return http.ListenAndServe(addr, handler)
}

func configure(app *dreego.App) error {
	return errors.Join(
		app.SetCSP("default-src 'self'; script-src 'self' 'unsafe-inline' https://cdn.tailwindcss.com; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' https:; connect-src 'self'; base-uri 'self'; form-action 'self'"),
	)
}

func hostRouter(public, product, blog http.Handler) http.Handler {
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
		public.ServeHTTP(w, r)
	})
}
