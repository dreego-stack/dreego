package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestTranspilerBasicPage(t *testing.T) {
	dreegotest.MustCompile(t, `<head><title>Test</title></head>
<go>msg := "hello"</go>
<div><h1>{msg}</h1></div>`)
}
