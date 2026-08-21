package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugComponentAttrSpace(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/components/Greet.dreego": `Component Greet (name string)
<div>Hello {{ name }}</div>`,
		"www/routes/get.dreego": `<go>name := "Ada"</go>
<div><@Greet name={ name }/></div>`,
	})
}
