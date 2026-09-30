package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugExpressionPipeInStringLiteral(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/app/routes/+page.dreego": "<server>\nname := \"x\"\n</server>\n<body><p>{{ name == \"a|b\" | upper }}</p></body>",
	})
	dreegotest.MustContain(t, gen["www/app/routes/dree.go"], `strings.ToUpper(fmt.Sprintf("%v", name == "a|b"))`)
}
