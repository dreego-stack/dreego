package main

import (
	"log"
	"os"

	dreego "github.com/dreego-stack/dreego/core"
	"github.com/dreego-stack/dreego/adapter/ssr"
	webapp "importcheck/www/app"
)

func main() {
	app := dreego.New(webapp.App)
	if err := app.SetCSP("default-src 'self'; script-src 'self' 'unsafe-inline' https://cdn.tailwindcss.com; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' https:; connect-src 'self'; base-uri 'self'; form-action 'self'"); err != nil {
		log.Fatal(err)
	}
	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}
	if err := ssr.Listen(app, addr); err != nil {
		log.Fatal(err)
	}
}
