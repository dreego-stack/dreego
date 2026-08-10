package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugScopedCSSMedia(t *testing.T) {
	gen := dreegotest.Build(t, map[string]string{
		"dreego/routes/get.dreego": `<style>@media (max-width: 600px) { p { color: red; } }</style>
<div><p>hi</p></div>`,
	})
	dreegotest.MustContain(t, gen["dreego/gen/routes.go"], "@media")
}
