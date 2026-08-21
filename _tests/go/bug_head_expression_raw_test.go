package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugHeadExpressionRaw(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/routes/get.dreego": `<go>
type Doc struct{ Title string }
doc := Doc{Title: "PeerNet Docs"}
</go>
<head><title>{{ doc.Title }} — PeerNet Docs</title></head>
<div><h1>{{ doc.Title }}</h1></div>`,
	})
	dreegotest.MustNotContain(t, gen["www/routes/dree.go"], "{{ doc.Title }}")
	dreegotest.MustContain(t, gen["www/routes/dree.go"], "doc.Title")
}
