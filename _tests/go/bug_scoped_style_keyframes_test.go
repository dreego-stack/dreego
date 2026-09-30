package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugScopedStyleKeyframes(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/app/routes/+page.dreego": `<style>@keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }</style>
<body><p>hi</p></body>`,
	})
	dreegotest.MustContain(t, gen["www/app/routes/dree.go"], "@keyframes spin")
}
