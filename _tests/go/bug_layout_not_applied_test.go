package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugLayoutNotApplied(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"dreego/layouts/default.dreego": `<div><html><body><nav>Nav</nav>{#slot}<footer>Footer</footer></body></html></div>`,
		"dreego/routes/get.dreego":      `<div><p>Page</p></div>`,
	})
	dreegotest.MustContain(t, gen["dreego/gen/routes.go"], "<html>")
	dreegotest.MustContain(t, gen["dreego/gen/routes.go"], "Nav")
	dreegotest.MustContain(t, gen["dreego/gen/routes.go"], "Footer")
	dreegotest.MustContain(t, gen["dreego/gen/routes.go"], "Page")
}
