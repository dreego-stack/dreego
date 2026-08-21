package main

import (
	"log"
	"os"

	dreego "github.com/dreego-stack/dreego/core"
	"hello/www"
)

func main() {
	app := dreego.New()
	if err := www.Register(app); err != nil {
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
