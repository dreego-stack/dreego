package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugHeadWithoutLayout(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/routes/get.dreego": `<head><script src="script.js"></script></head>
<body><p>hi</p></body>`,
	})
	dreegotest.MustContain(t, gen["www/routes/dree.go"], "script.js")
}
