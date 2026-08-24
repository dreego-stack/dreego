package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugComponentMultiGo(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/components/Greet.dreego": `Component Greet (name string)
<server>greeting := "hello"</server>
<server>msg := greeting + " world"</server>
<body>{{ msg }} {{ name }}</body>`,
		"www/routes/get.dreego": `<body><@Greet name="Ada"/></body>`,
	})
	dreegotest.MustContain(t, gen["www/components/dree.go"], `greeting := "hello"`)
	dreegotest.MustContain(t, gen["www/components/dree.go"], `msg := greeting + " world"`)
}
