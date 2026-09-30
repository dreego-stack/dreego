package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugScopedStyleDeclarationsLost(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/app/routes/+page.dreego": `<style>p { background: radial-gradient(circle, red, blue); }</style>
<body><p>hi</p></body>`,
	})
	dreegotest.MustContain(t, gen["www/app/routes/dree.go"], "radial-gradient(circle, red, blue)")
}
