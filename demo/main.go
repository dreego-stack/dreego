package main

import (
	_ "codeberg.org/dreego/dreego/demo/dreego/routes"
	_ "codeberg.org/dreego/dreego/demo/dreego/routes/about"
	_ "codeberg.org/dreego/dreego/demo/dreego/routes/users/_id_"
	"codeberg.org/dreego/dreego/pkg/runtime"
)

func main() {
	runtime.Listen(":8080")
}
