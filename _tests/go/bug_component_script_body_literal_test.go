package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugComponentScriptBodyLiteral(t *testing.T) {
	gen := dreegotest.Build(t, map[string]string{
		"dreego/components/Snippet.dreego": `Component Snippet (x string)
<div><script>const s = "literal {x}";</script></div>`,
		"dreego/routes/get.dreego": `<div><@Snippet x=42/></div>`,
	})
	dreegotest.MustContain(t, gen["dreego/gen/components.go"], "literal {x}")
}
