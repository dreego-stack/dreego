package main

import (
	"log"
	"os"

	webapp "forms/www/app"
	dreego "github.com/dreego-stack/dreego/core"
	"github.com/dreego-stack/dreego/adapter/ssr"
)

func main() {
	store := dreego.NewCookieStore([]byte("reference-apps-secret-key-32-bytes!"))
	app := dreego.New(webapp.App)
	if err := app.SetSessionStore(store); err != nil {
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
