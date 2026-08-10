package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugComponentQuotedProp(t *testing.T) {
	gen := dreegotest.Build(t, map[string]string{
		"dreego/components/Card.dreego": `Component Card (title string, active bool)
<div><h1>{title}</h1><span>{active}</span></div>`,
		"dreego/routes/get.dreego": `<go>myTitle := "Hello"</go>
<div><@Card title="{myTitle}" active={true}/></div>`,
	})
	dreegotest.MustNotContain(t, gen["dreego/gen/routes.go"], "{myTitle}")
	dreegotest.MustContain(t, gen["dreego/gen/routes.go"], "Card(myTitle, true)")
}
