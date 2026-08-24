package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestTranspilerBasicPage(t *testing.T) {
	t.Parallel()
	dreegotest.MustCompile(t, `<head><title>Test</title></head>
<server>msg := "hello"</server>
<body><h1>{{ msg }}</h1></body>`)
}
