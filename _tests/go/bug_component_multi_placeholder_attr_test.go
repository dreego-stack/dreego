package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugComponentMultiPlaceholderAttr(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/components/Card.dreego": `Component Card (url string)
<body><a href="{{ url }}">go</a></body>`,
		"www/routes/get.dreego": `<server>left := "x"; right := "y"</server>
<body><@Card url={left + "-" + right}/></body>`,
	})
	dreegotest.MustNotContain(t, gen["www/routes/dree.go"], "left}-{right")
}
