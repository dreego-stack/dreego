package main

import (
	_ "demo/dreego/gen"
	dreego "github.com/dreego-stack/dreego/core"
)

func main() {
	dreego.Listen(":8080")
}
