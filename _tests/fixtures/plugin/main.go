package main

import (
	"log"
	"os"

	dreego "github.com/dreego-stack/dreego/core"
	"github.com/dreego-stack/dreego/adapter/ssr"
	"pluginapp/plugin"
	webapp "pluginapp/www/app"
)

func main() {
	app := dreego.New(webapp.App)
	if err := plugin.Register(app, plugin.Options{Prefix: "/plugin"}); err != nil {
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
