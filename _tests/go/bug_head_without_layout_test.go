package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugHeadWithoutLayout(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"dreego/routes/get.dreego": `<head><script src="script.js"></script></head>
<div><p>hi</p></div>`,
	})
	dreegotest.MustContain(t, gen["dreego/gen/routes.go"], "script.js")
}
