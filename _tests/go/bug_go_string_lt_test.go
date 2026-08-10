package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugGoStringLt(t *testing.T) {
	gen := dreegotest.Build(t, map[string]string{
		"dreego/routes/get.dreego": "<go>\nmsg := \"TO: <HASH>\"\nsvg := `<svg viewBox=\"0 0 24 24\"><path d=\"M12 2\"/></svg>`\n</go>\n<div>\n<p>{msg}</p>\n<p>{svg}</p>\n</div>",
	})
	dreegotest.MustContain(t, gen["dreego/gen/routes.go"], "TO: <HASH>")
	dreegotest.MustContain(t, gen["dreego/gen/routes.go"], `<svg viewBox="0 0 24 24">`)
	dreegotest.MustContain(t, gen["dreego/gen/routes.go"], `<path d="M12 2"/>`)
}
