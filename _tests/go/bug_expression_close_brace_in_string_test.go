package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugExpressionCloseBraceInStringLiteral(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/app/routes/+page.dreego": "<server>\nname := \"}}\"\n</server>\n<body><p>{{ name == \"}}\" }}</p></body>",
	})
	dreegotest.MustContain(t, gen["www/app/routes/dree.go"], `name == "}}"`)
}
