package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugScopedCSSMedia(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/routes/get.dreego": `<style>@media (max-width: 600px) { p { color: red; } }</style>
<div><p>hi</p></div>`,
	})
	dreegotest.MustContain(t, gen["www/routes/dree.go"], "@media")
}
