package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugComponentMultiPlaceholderAttr(t *testing.T) {
	gen := dreegotest.Build(t, map[string]string{
		"dreego/components/Card.dreego": `Component Card (url string)
<div><a href="{url}">go</a></div>`,
		"dreego/routes/get.dreego": `<go>left := "x"; right := "y"</go>
<div><@Card href="{left}-{right}"/></div>`,
	})
	dreegotest.MustNotContain(t, gen["dreego/gen/routes.go"], "left}-{right")
}
