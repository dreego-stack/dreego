package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestTranspilerComment(t *testing.T) {
	dreegotest.MustCompile(t, `<head><title>T</title></head>
<go>x := "hi"</go>
<div><p>{x}</p></div>`)
}
