package main

import (
	"errors"
	"log"

	"demo/www"
	dreego "github.com/dreego-stack/dreego/core"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	app := dreego.New()
	return errors.Join(
		app.SetCSP("default-src 'self'; script-src 'self' 'unsafe-inline' https://cdn.tailwindcss.com; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' https:; connect-src 'self'; base-uri 'self'; form-action 'self'"),
		www.Register(app),
		app.Listen(":8080"),
	)
}
