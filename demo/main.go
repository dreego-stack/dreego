package main

import (
	"demo/dreego/gen"
	dreego "github.com/dreego-stack/dreego/core"
)

func main() {
	app := dreego.New()
	gen.Register(app)
	app.Listen(":8080")
}
