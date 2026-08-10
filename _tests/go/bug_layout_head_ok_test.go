package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugLayoutHeadOk(t *testing.T) {
	gen := dreegotest.Build(t, map[string]string{
		"dreego/layouts/default.dreego": `<head><title>Layout Title</title><link rel="stylesheet" href="cdn.tailwindcss.com"></head>
<div>{#slot}</div>`,
		"dreego/routes/get.dreego": `<head><meta name="description" content="route meta"></head>
<div><p>hi</p></div>`,
	})
	dreegotest.MustContain(t, gen["dreego/gen/routes.go"], "cdn.tailwindcss.com")
	dreegotest.MustContain(t, gen["dreego/gen/routes.go"], "Layout Title")
	dreegotest.MustContain(t, gen["dreego/gen/routes.go"], "route meta")
}
