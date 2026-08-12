package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugHeadDroppedByLayout(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"dreego/layouts/default.dreego": `<head><title>Layout Title</title></head>
<div>{#slot}</div>`,
		"dreego/routes/get.dreego": `<head><script src="route-script.js"></script></head>
<div><p>hi</p></div>`,
	})
	dreegotest.MustContain(t, gen["dreego/gen/routes.go"], "Layout Title")
	dreegotest.MustContain(t, gen["dreego/gen/routes.go"], "route-script.js")
}
