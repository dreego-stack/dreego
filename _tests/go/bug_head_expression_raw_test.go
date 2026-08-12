package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugHeadExpressionRaw(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"dreego/routes/get.dreego": `<go>
type Doc struct{ Title string }
doc := Doc{Title: "PeerNet Docs"}
</go>
<head><title>{doc.Title} — PeerNet Docs</title></head>
<div><h1>{doc.Title}</h1></div>`,
	})
	dreegotest.MustNotContain(t, gen["dreego/gen/routes.go"], "{doc.Title}")
	dreegotest.MustContain(t, gen["dreego/gen/routes.go"], "doc.Title")
}
