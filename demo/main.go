package main

import (
	_ "codeberg.org/dreego/dreego/demo/dreego/routes"
	_ "codeberg.org/dreego/dreego/demo/dreego/routes/users"
	"codeberg.org/dreego/dreego/pkg/runtime"
)

func main() {
	runtime.Listen(":8080")
}
