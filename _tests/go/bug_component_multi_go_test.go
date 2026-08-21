package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugComponentMultiGo(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/components/Greet.dreego": `Component Greet (name string)
<go>greeting := "hello"</go>
<go>msg := greeting + " world"</go>
<div>{{ msg }} {{ name }}</div>`,
		"www/routes/get.dreego": `<div><@Greet name="Ada"/></div>`,
	})
	dreegotest.MustContain(t, gen["www/components/dree.go"], `greeting := "hello"`)
	dreegotest.MustContain(t, gen["www/components/dree.go"], `msg := greeting + " world"`)
}
