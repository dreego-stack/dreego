package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugLayoutHeadOk(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/layouts/default.dreego": `<head><title>Layout Title</title><link rel="stylesheet" href="cdn.tailwindcss.com"></head>
<body>{#slot}</body>`,
		"www/routes/get.dreego": `<head><meta name="description" content="route meta"></head>
<body><p>hi</p></body>`,
	})
	dreegotest.MustContain(t, gen["www/routes/dree.go"], "cdn.tailwindcss.com")
	dreegotest.MustContain(t, gen["www/routes/dree.go"], "Layout Title")
	dreegotest.MustContain(t, gen["www/routes/dree.go"], "route meta")
}
