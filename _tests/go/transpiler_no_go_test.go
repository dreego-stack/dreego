package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestTranspilerNoGo(t *testing.T) {
	dreegotest.MustCompile(t, `<head><title>No Go</title></head>
<div><h1>Static</h1></div>`)
}
