package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugLayoutHeadLost(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/layouts/default.dreego": `<head><title>Layout Title</title><link rel="stylesheet" href="cdn.tailwindcss.com"></head>
<div>{#slot}</div>`,
		"www/routes/get.dreego": `<head><meta name="description" content="route meta"></head>
<div><p>hi</p></div>`,
	})
	dreegotest.MustContain(t, gen["www/routes/dree.go"], "cdn.tailwindcss.com")
	dreegotest.MustContain(t, gen["www/routes/dree.go"], "Layout Title")
	dreegotest.MustContain(t, gen["www/routes/dree.go"], "route meta")
}
