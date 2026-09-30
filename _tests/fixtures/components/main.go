package main

import (
	"log"
	"os"

	webapp "components/www/app"
	dreego "github.com/dreego-stack/dreego/core"
	"github.com/dreego-stack/dreego/adapter/ssr"
)

func main() {
	app := dreego.New(webapp.App)
	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}
	if err := ssr.Listen(app, addr); err != nil {
		log.Fatal(err)
	}
}
