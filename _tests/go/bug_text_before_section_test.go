package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugTextBeforeSection(t *testing.T) {
	gen := dreegotest.Build(t, map[string]string{
		"dreego/routes/get.dreego": `<!doctype html>
<html lang="en">
<go>msg := "hi"</go>
<div><p>{msg}</p></div>`,
	})
	dreegotest.MustNotContain(t, gen["dreego/gen/routes.go"], `html.EscapeString("msg := \"hi\""`)
	dreegotest.MustContain(t, gen["dreego/gen/routes.go"], `msg := "hi"`)
}
