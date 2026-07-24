package main

import (
	_ "codeberg.org/dreego/dreego/demo/dreego/gen"
	"codeberg.org/dreego/dreego/pkg/runtime"
)

func main() {
	runtime.Listen(":8080")
}
