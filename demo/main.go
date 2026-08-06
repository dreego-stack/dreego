package main

import (
	_ "demo/dreego/gen"
	dreego "codeberg.org/dreego/dreego/core"
)

func main() {
	dreego.Listen(":8080")
}
