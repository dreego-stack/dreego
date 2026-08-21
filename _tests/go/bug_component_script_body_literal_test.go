package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugComponentScriptBodyLiteral(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/components/Snippet.dreego": `Component Snippet (x string)
<div><script>const s = "literal {x}";</script></div>`,
		"www/routes/get.dreego": `<div><@Snippet x=42/></div>`,
	})
	dreegotest.MustContain(t, gen["www/components/dree.go"], "literal {x}")
}
