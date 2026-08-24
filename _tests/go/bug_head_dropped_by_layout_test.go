package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugHeadDroppedByLayout(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/layouts/default.dreego": `<head><title>Layout Title</title></head>
<body>{#slot}</body>`,
		"www/routes/get.dreego": `<head><script src="route-script.js"></script></head>
<body><p>hi</p></body>`,
	})
	dreegotest.MustContain(t, gen["www/routes/dree.go"], "Layout Title")
	dreegotest.MustContain(t, gen["www/routes/dree.go"], "route-script.js")
}
