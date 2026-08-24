package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugComponentScriptBodyLiteral(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/components/Snippet.dreego": `Component Snippet (x string)
<body><script>const s = "literal {x}";</script></body>`,
		"www/routes/get.dreego": `<body><@Snippet x=42/></body>`,
	})
	dreegotest.MustContain(t, gen["www/components/dree.go"], "literal {x}")
}
