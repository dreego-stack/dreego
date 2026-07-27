package main

import (
	_ "codeberg.org/dreego/dreego/demo/dreego/gen"

	core "codeberg.org/dreego/dreego/dreego-core"
)

func main() {
	core.SetSessionStore(core.NewCookieStore([]byte("demo-secret-change-me")))
	core.Listen(":8080")
}
