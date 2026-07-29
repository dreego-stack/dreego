package main

import (
	_ "gen"

	core "codeberg.org/dreego/dreego/core"
)

func main() {
	core.Listen(":8080")
}
