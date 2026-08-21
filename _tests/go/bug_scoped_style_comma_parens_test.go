package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugScopedStyleCommaParens(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/routes/get.dreego": `<style>a, b { color: rgb(1, 2, 3); }</style>
<div><p>hi</p></div>`,
	})
	dreegotest.MustContain(t, gen["www/routes/dree.go"], "rgb(1, 2, 3)")
}
