package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestComponentForwardReferenceBindsPropsByName(t *testing.T) {
	t.Parallel()
	generated := dreegotest.Build(t, map[string]string{
		"www/components/Alpha.dreego": "Component Alpha ()\n<div><@Zulu second=\"B\" first=\"A\"/></div>",
		"www/components/Zulu.dreego":  "Component Zulu (first string, second string)\n<div><p>{{ first }}{{ second }}</p></div>",
		"www/routes/get.dreego":       `<div><@Alpha/></div>`,
	})
	components := generated["www/components/dree.go"]
	if !strings.Contains(components, `Zulu("A", "B")`) {
		t.Fatalf("forward component props were not bound by name, got:\n%s", components)
	}
}
