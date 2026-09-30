package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugHeadExpressionRaw(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/app/routes/+page.dreego": `<server>
type Doc struct{ Title string }
doc := Doc{Title: "PeerNet Docs"}
</server>
<head><title>{{ doc.Title }} — PeerNet Docs</title></head>
<body><h1>{{ doc.Title }}</h1></body>`,
	})
	dreegotest.MustNotContain(t, gen["www/app/routes/dree.go"], "{{ doc.Title }}")
	dreegotest.MustContain(t, gen["www/app/routes/dree.go"], "doc.Title")
}
