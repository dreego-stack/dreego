package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestTranspilerAllSections(t *testing.T) {
	t.Parallel()
	dreegotest.MustCompile(t, `<head><title>All</title></head>
<server>x := "ok"</server>
<body><p>{{ x }}</p></body>
<client>console.log("js")</client>
<style>p { color: red; }</style>`)
}
