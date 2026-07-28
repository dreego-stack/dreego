package main

import (
	_ "§$name$§/dreego/gen"
	core "codeberg.org/dreego/dreego/dreego-core"
)

func main() {
	core.Listen(":8080")
}
