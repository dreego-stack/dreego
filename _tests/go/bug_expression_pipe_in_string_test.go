package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugExpressionPipeInStringLiteral(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"dreego/routes/get.dreego": "<go>\nname := \"x\"\n</go>\n<div><p>{{ name == \"a|b\" | upper }}</p></div>",
	})
	dreegotest.MustContain(t, gen["dreego/gen/routes.go"], `strings.ToUpper(fmt.Sprintf("%v", name == "a|b"))`)
}
