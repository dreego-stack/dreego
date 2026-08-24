package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestTranspilerComment(t *testing.T) {
	t.Parallel()
	dreegotest.MustCompile(t, `<head><title>T</title></head>
<server>x := "hi"</server>
<body><p>{{ x }}</p></body>`)
}
