package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugScopedStyleKeyframes(t *testing.T) {
	gen := dreegotest.Build(t, map[string]string{
		"dreego/routes/get.dreego": `<style>@keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }</style>
<div><p>hi</p></div>`,
	})
	dreegotest.MustContain(t, gen["dreego/gen/routes.go"], "@keyframes spin")
}
