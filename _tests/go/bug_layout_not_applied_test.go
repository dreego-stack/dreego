package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugLayoutNotApplied(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/layouts/default.dreego": `<body><html><body><nav>Nav</nav>{#slot}<footer>Footer</footer></body></html></body>`,
		"www/routes/get.dreego":      `<body><p>Page</p></body>`,
	})
	dreegotest.MustContain(t, gen["www/layouts/dree.go"], "<html>")
	dreegotest.MustContain(t, gen["www/layouts/dree.go"], "Nav")
	dreegotest.MustContain(t, gen["www/layouts/dree.go"], "Footer")
	dreegotest.MustContain(t, gen["www/routes/dree.go"], "Page")
	dreegotest.MustContain(t, gen["www/routes/dree.go"], "layouts.Default(c, pageContent, head)")
}
