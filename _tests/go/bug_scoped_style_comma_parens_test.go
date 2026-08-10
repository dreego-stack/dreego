package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugScopedStyleCommaParens(t *testing.T) {
	gen := dreegotest.Build(t, map[string]string{
		"dreego/routes/get.dreego": `<style>a, b { color: rgb(1, 2, 3); }</style>
<div><p>hi</p></div>`,
	})
	dreegotest.MustContain(t, gen["dreego/gen/routes.go"], "rgb(1, 2, 3)")
}
