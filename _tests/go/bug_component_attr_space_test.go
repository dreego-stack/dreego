package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugComponentAttrSpace(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/components/Greet.dreego": `Component Greet (name string)
<body>Hello {{ name }}</body>`,
		"www/routes/get.dreego": `<server>name := "Ada"</server>
<body><@Greet name={ name }/></body>`,
	})
}
