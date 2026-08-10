package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugCleanSegmentOptional(t *testing.T) {
	gen := dreegotest.Build(t, map[string]string{
		"dreego/routes/[[opt]]/get.dreego": `<div>optional</div>`,
	})
	dreegotest.MustNotContain(t, gen["dreego/gen/routes.go"], "[opt]")
	dreegotest.MustContain(t, gen["dreego/gen/routes.go"], `dreego.Register("GET", "/{opt}",`)
}
