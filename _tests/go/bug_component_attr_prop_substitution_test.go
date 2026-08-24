package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugComponentAttrPropSubstitution(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/components/Link.dreego": `Component Link (url string, label string)
<body><a href="{{ url }}">{{ label }}</a></body>`,
		"www/routes/get.dreego": `<body><@Link url="https://example.com" label="Go"/></body>`,
	})
	dreegotest.MustNotContain(t, gen["www/components/dree.go"], "{{ url }}")
	dreegotest.MustNotContain(t, gen["www/components/dree.go"], "{{ label }}")
	dreegotest.MustContain(t, gen["www/components/dree.go"], "dreego.SafeURL")
	dreegotest.MustContain(t, gen["www/components/dree.go"], "dreego.SafeText")
}
