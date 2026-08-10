package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugComponentAttrPropSubstitution(t *testing.T) {
	gen := dreegotest.Build(t, map[string]string{
		"dreego/components/Link.dreego": `Component Link (url string, label string)
<a href="{url}">{label}</a>`,
		"dreego/routes/get.dreego": `<div><@Link url="https://example.com" label="Go"/></div>`,
	})
	dreegotest.MustNotContain(t, gen["dreego/gen/components.go"], "{url}")
	dreegotest.MustNotContain(t, gen["dreego/gen/components.go"], "{label}")
	dreegotest.MustContain(t, gen["dreego/gen/components.go"], "EscapeString")
}
