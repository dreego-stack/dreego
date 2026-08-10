package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugComponentAttrSpace(t *testing.T) {
	dreegotest.MustBuild(t, map[string]string{
		"dreego/components/Greet.dreego": `Component Greet (name string)
<div>Hello {name}</div>`,
		"dreego/routes/get.dreego": `<go>name := "Ada"</go>
<div><@Greet name={ name }/></div>`,
	})
}
