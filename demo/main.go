package main

import (
	_ "demo/dreego/gen"
	core "codeberg.org/dreego/dreego/core"
)

func main() {
	core.Listen(":8080")
}
