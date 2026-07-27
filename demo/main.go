package main

import (
	_ "codeberg.org/dreego/dreego/demo/dreego/gen"
	"codeberg.org/dreego/dreego/pkg/runtime"
	"codeberg.org/dreego/dreego/pkg/session"
)

func main() {
	runtime.SetSessionStore(session.NewCookieStore([]byte("demo-secret-change-me")))
	runtime.Listen(":8080")
}
