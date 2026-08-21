package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugComponentMultiPlaceholderAttr(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/components/Card.dreego": `Component Card (url string)
<div><a href="{{ url }}">go</a></div>`,
		"www/routes/get.dreego": `<go>left := "x"; right := "y"</go>
<div><@Card url={left + "-" + right}/></div>`,
	})
	dreegotest.MustNotContain(t, gen["www/routes/dree.go"], "left}-{right")
}
