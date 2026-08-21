package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugComponentQuotedProp(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/components/Card.dreego": `Component Card (title string, active bool)
<div><h1>{{ title }}</h1><span>{{ active }}</span></div>`,
		"www/routes/get.dreego": `<go>myTitle := "Hello"</go>
<div><@Card title={myTitle} active={true}/></div>`,
	})
	dreegotest.MustNotContain(t, gen["www/routes/dree.go"], "{myTitle}")
	dreegotest.MustContain(t, gen["www/routes/dree.go"], "Card(myTitle, true)")
}
