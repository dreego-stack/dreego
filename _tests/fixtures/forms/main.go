package main

import (
	"log"
	"os"

	"forms/dreego/gen"
	dreego "github.com/dreego-stack/dreego/core"
)

func main() {
	app := dreego.New()
	store := dreego.NewCookieStore([]byte("reference-apps-secret-key-32-bytes!"))
	if err := app.SetSessionStore(store); err != nil {
		log.Fatal(err)
	}
	if err := gen.Register(app); err != nil {
		log.Fatal(err)
	}
	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}
	if err := app.Listen(addr); err != nil {
		log.Fatal(err)
	}
}
