package main

import (
	"log"

	"demo/www"
	dreego "github.com/dreego-stack/dreego/core"
)

func main() {
	app := dreego.New()
	if err := www.Register(app); err != nil {
		log.Fatal(err)
	}
	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
