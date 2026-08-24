package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugGoStringLt(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/routes/get.dreego": "<server>\nmsg := \"TO: <HASH>\"\nsvg := `<svg viewBox=\"0 0 24 24\"><path d=\"M12 2\"/></svg>`\n</server>\n<body>\n<p>{{ msg }}</p>\n<p>{{ svg }}</p>\n</body>",
	})
	dreegotest.MustContain(t, gen["www/routes/dree.go"], "TO: <HASH>")
	dreegotest.MustContain(t, gen["www/routes/dree.go"], `<svg viewBox="0 0 24 24">`)
	dreegotest.MustContain(t, gen["www/routes/dree.go"], `<path d="M12 2"/>`)
}
