package main

import (
	"log"
	"os"

	dreego "github.com/dreego-stack/dreego/core"
	"hello/dreego/gen"
)

func main() {
	app := dreego.New()
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
