package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugExpressionCloseBraceInStringLiteral(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/routes/get.dreego": "<go>\nname := \"}}\"\n</go>\n<div><p>{{ name == \"}}\" }}</p></div>",
	})
	dreegotest.MustContain(t, gen["www/routes/dree.go"], `name == "}}"`)
}
